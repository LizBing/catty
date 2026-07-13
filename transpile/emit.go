// Package transpile lowers a method's IR to Go source — the AOT transpiler's
// emitter (ROADMAP A1/A2). Each operand-stack value becomes a Go local; bytecode
// control flow becomes goto/labels; the Go toolchain compiles it, with the Go
// runtime as GC/scheduler.
//
// A2.1 uses fresh-per-def temps (each def = a new typed Go local; uses reference
// the defining temp), which resolves the slot-type-reuse that broke ref methods
// under A1's position-stable slots. Scope: static, merge-free methods (straight-
// line + one-armed if); int + ref(+array) types. Merges (loops/diamonds → phis),
// fields, and the invoke bridge are later milestones.
package transpile

import (
	"fmt"
	"strings"

	"catty/classfile"
	"catty/lowering"
	"catty/opcode"
	"catty/rtda"
)

// Emit turns one method's IR into a Go function definition. See the package doc
// for scope. Returns an error (not wrong code) for methods outside scope. The
// loader resolves field offsets (getfield/putfield) at emit time.
func Emit(method *rtda.Method, ir *lowering.IR, loader rtda.Loader, emittable map[string]bool) (string, error) {
	if !method.IsStatic() {
		return "", fmt.Errorf("transpile: A2.1 supports static methods only (got %s)", method.Name())
	}
	cp := method.Owner().ConstantPool()
	targets := collectTargets(ir)

	e := &emitter{
		slotTemp: map[int]string{}, slotType: map[int]string{}, loader: loader,
		merges: map[int]bool{}, fallThrough: map[int]bool{},
		mergeTemps: map[int][]string{}, mergeTempTypes: map[int][]string{},
		localMap: buildLocalMap(method),
		emittable: emittable,
	}
	e.merges, e.fallThrough = cfgAnalysis(ir)
	if err := e.allocMergeTemps(ir, method); err != nil { // validates + refuses deferred cases
		return "", err
	}

	var body strings.Builder
	for pc := 0; pc < len(ir.Insts); pc++ {
		inst := &ir.Insts[pc]
		if !inst.Present {
			continue
		}
		// Fall-through into a non-empty-stack merge: copy the predecessor's temps
		// into the merge temps, then read from them after the join (the phi).
		if e.merges[pc] && e.fallThrough[pc] && len(e.mergeTemps[pc]) > 0 {
			e.emitMergeCopies(&body, pc)
			e.resetToMergeTemps(pc)
		}
		if targets[pc] {
			fmt.Fprintf(&body, "pc%d:\n", pc)
		}
		// A branch into a non-empty-stack merge: copy this path's temps into the
		// merge temps before jumping (the phi, predecessor-edge side).
		if isBranch(inst.Op) && e.merges[int(inst.Branch)] && len(e.mergeTemps[int(inst.Branch)]) > 0 {
			e.emitMergeCopies(&body, int(inst.Branch))
		}
		if err := e.emitOne(&body, inst, cp); err != nil {
			return "", fmt.Errorf("at pc %d: %w", pc, err)
		}
	}

	var b strings.Builder
	emitSignature(&b, method, e.temps, localTypes(ir))
	b.WriteString(body.String())
	emitSink(&b, e.temps)
	emitTerminator(&b, method)
	b.WriteString("}\n")
	return b.String(), nil
}

// emitter carries the fresh-per-def state: the current temp name per operand-stack
// slot, and the list of all temps (for top-of-function declarations + the sink).
type emitter struct {
	slotTemp       map[int]string   // slot index → temp name of its most recent def
	slotType       map[int]string   // slot index → that temp's Go type (for dup)
	loader         rtda.Loader
	temps          []tempDecl
	counter        int
	merges         map[int]bool      // pcs with >1 predecessor (control-flow merges)
	fallThrough    map[int]bool      // pcs reached by fall-through from a predecessor
	mergeTemps     map[int][]string  // merge pc → temp name per stack slot (the phi)
	mergeTempTypes map[int][]string  // merge pc → Go type per stack slot
	localMap       map[int]string    // JVM local slot → Go param name (cat-2 aware)
	emittable      map[string]bool   // AOT'd methods in this build (invokestatic dispatch)
}

type tempDecl struct{ name, gotype string }

// defTemp allocates a fresh temp of goType for a def of slot, records it, and
// returns the temp name.
func (e *emitter) defTemp(slot int, goType string) string {
	name := fmt.Sprintf("t%d", e.counter)
	e.counter++
	e.temps = append(e.temps, tempDecl{name, goType})
	e.slotTemp[slot] = name
	e.slotType[slot] = goType
	return name
}

// defTempCat2 allocates one temp for a category-2 def (long/double spanning 2
// slots) and binds BOTH slots to it — the Go value is one int64/float64.
func (e *emitter) defTempCat2(slot int, goType string) string {
	name := e.defTemp(slot, goType)
	e.slotTemp[slot+1] = name
	e.slotType[slot+1] = goType
	return name
}

func (e *emitter) use(slot int) string { return e.slotTemp[slot] }

// emitOne emits the Go statement(s) for one IR instruction.
func (e *emitter) emitOne(b *strings.Builder, inst *lowering.IRInst, cp *classfile.ConstantPool) error {
	w := func(format string, args ...any) { fmt.Fprintf(b, "\t"+format+"\n", args...) }
	switch inst.Op {

	// --- int constants ---
	case opcode.IconstM1, opcode.Iconst0, opcode.Iconst1, opcode.Iconst2,
		opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %d", t, int32(inst.Op-opcode.Iconst0))
	case opcode.Bipush:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %d", t, int32(inst.Const8))
	case opcode.Sipush:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %d", t, int32(inst.Const16))
	case opcode.Ldc, opcode.LdcW:
		switch cp.Tag(inst.Index) {
		case classfile.ConstantInteger:
			t := e.defTemp(int(inst.Defs[0]), "int32")
			w("%s = %d", t, cp.Integer(inst.Index))
		case classfile.ConstantFloat:
			t := e.defTemp(int(inst.Defs[0]), "float32")
			w("%s = %g", t, cp.Float(inst.Index))
		case classfile.ConstantString:
			t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
			w("%s = runtime.NewString(%q)", t, cp.String(inst.Index))
		default:
			return fmt.Errorf("ldc: unsupported constant tag %d", cp.Tag(inst.Index))
		}
	case opcode.Ldc2W:
		switch cp.Tag(inst.Index) {
		case classfile.ConstantLong:
			t := e.defTempCat2(int(inst.Defs[0]), "int64")
			w("%s = %d", t, cp.Long(inst.Index))
		case classfile.ConstantDouble:
			t := e.defTempCat2(int(inst.Defs[0]), "float64")
			w("%s = %g", t, cp.Double(inst.Index))
		default:
			return fmt.Errorf("ldc2_w: unsupported constant tag %d", cp.Tag(inst.Index))
		}

	// --- getstatic (via the runtime bridge) ---
	case opcode.Getstatic:
		className, name, desc := cp.MemberRef(inst.Index)
		goType, err := descToGo(desc)
		if err != nil {
			return err
		}
		t := e.defTemp(int(inst.Defs[0]), goType)
		call := fmt.Sprintf("runtime.GetStatic(%q, %q, %q)", className, name, desc)
		w("%s = %s", t, slotExtract(call, desc))

	// --- invokevirtual (native target via the runtime bridge) ---
	case opcode.Invokevirtual:
		className, name, desc := cp.MemberRef(inst.Index)
		md := rtda.ParseMethodDescriptor(desc)
		argSlots := md.ArgSlots()
		// args[0] = receiver; Uses[0] is the receiver, Uses[1..] the params.
		args := []string{"rtda.RefSlot(" + e.use(int(inst.Uses[0])) + ")"}
		for i := 0; i < argSlots; i++ {
			temp := e.use(int(inst.Uses[1+i]))
			args = append(args, slotConstructor(md.ParameterTypes[i], temp))
		}
		call := fmt.Sprintf("runtime.InvokeVirtual(%q, %q, %q, []rtda.Slot{%s})",
			className, name, desc, strings.Join(args, ", "))
		if md.ReturnType == "V" {
			w("%s", call)
		} else {
			goType, err := descToGo(md.ReturnType)
			if err != nil {
				return err
			}
			t := e.defTemp(int(inst.Defs[0]), goType)
			w("%s = %s", t, slotExtract(call, md.ReturnType))
		}

	// --- loads ---
	case opcode.Iload, opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %s", t, e.localName(inst))
	case opcode.Aload, opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
		w("%s = %s", t, e.localName(inst))

	// --- stores ---
	case opcode.Istore, opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		w("%s = %s", e.localName(inst), e.use(int(inst.Uses[0])))
	case opcode.Astore, opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		w("%s = %s", e.localName(inst), e.use(int(inst.Uses[0])))

	// --- int arithmetic ---
	case opcode.Iadd, opcode.Isub, opcode.Imul, opcode.Idiv, opcode.Irem,
		opcode.Iand, opcode.Ior, opcode.Ixor:
		// Read uses before allocating the def temp: the def often reuses an
		// operand's slot (e.g. iadd writes to slot d-2 = Uses[0]).
		a, b := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %s %s %s", t, a, binop(inst.Op), b)
	case opcode.Ineg:
		a := e.use(int(inst.Uses[0]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = -%s", t, a)
	case opcode.Ishl, opcode.Ishr, opcode.Iushr:
		a, b := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %s", t, shiftExpr(inst.Op, a, b))
	case opcode.Iinc:
		w("%s += %d", e.localName(inst), int32(inst.Const8))

	// --- arrays ---
	case opcode.Iaload, opcode.Baload, opcode.Caload, opcode.Saload:
		arr, idx := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = %s.ArrayElementSlot(int(%s)).Num()", t, arr, idx)
	case opcode.Aaload:
		arr, idx := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
		w("%s = %s.ArrayElementSlot(int(%s)).Ref()", t, arr, idx)
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore:
		arr, idx, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1])), e.use(int(inst.Uses[2]))
		w("%s.ArrayElementSlot(int(%s)).SetNum(%s)", arr, idx, val)
	case opcode.Aastore:
		arr, idx, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1])), e.use(int(inst.Uses[2]))
		w("%s.ArrayElementSlot(int(%s)).SetRef(%s)", arr, idx, val)
	case opcode.Arraylength:
		arr := e.use(int(inst.Uses[0]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(%s.ArrayLength())", t, arr)

	// --- refs ---
	case opcode.AconstNull:
		t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
		w("%s = (*rtda.Object)(nil)", t)

	// --- branches ---
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle:
		w("if %s %s 0 { goto pc%d }", e.use(int(inst.Uses[0])), cmp0(inst.Op), inst.Branch)
	case opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple:
		w("if %s %s %s { goto pc%d }", e.use(int(inst.Uses[0])), icmp(inst.Op), e.use(int(inst.Uses[1])), inst.Branch)
	case opcode.IfAcmpeq, opcode.IfAcmpne:
		w("if %s %s %s { goto pc%d }", e.use(int(inst.Uses[0])), icmp(inst.Op), e.use(int(inst.Uses[1])), inst.Branch)
	case opcode.Ifnull:
		w("if %s == nil { goto pc%d }", e.use(int(inst.Uses[0])), inst.Branch)
	case opcode.Ifnonnull:
		w("if %s != nil { goto pc%d }", e.use(int(inst.Uses[0])), inst.Branch)
	case opcode.Goto, opcode.GotoW:
		w("goto pc%d", inst.Branch)

	// --- returns ---
	case opcode.Ireturn, opcode.Areturn:
		w("return %s", e.use(int(inst.Uses[0])))
	case opcode.Return:
		w("return")

	// --- invokestatic: direct Go call to the mangled (emitted) target ---
	case opcode.Invokestatic:
		className, name, desc := cp.MemberRef(inst.Index)
		md := rtda.ParseMethodDescriptor(desc)
		// Go args: raw for direct AOT calls, Slot-wrapped for the runtime bridge.
		rawArgs := make([]string, 0, len(md.ParameterTypes))
		bridgeArgs := make([]string, 0, len(md.ParameterTypes))
		slot := 0
		for i, p := range md.ParameterTypes {
			temp := e.use(int(inst.Uses[slot]))
			rawArgs = append(rawArgs, temp)
			bridgeArgs = append(bridgeArgs, slotConstructor(md.ParameterTypes[i], temp))
			if p == "J" || p == "D" {
				slot += 2
			} else {
				slot++
			}
		}
		// If the target is AOT'd in this build, emit a direct Go call; otherwise
		// route through the runtime.InvokeStatic bridge (interpreted/native).
		key := className + "\x00" + name + "\x00" + desc
		if e.emittable == nil || e.emittable[key] {
			call := fmt.Sprintf("%s(%s)", mangle(className, name), strings.Join(rawArgs, ", "))
			if md.ReturnType == "V" {
				w("%s", call)
			} else {
				goType, err := descToGo(md.ReturnType)
				if err != nil {
					return err
				}
				if md.ReturnType == "J" || md.ReturnType == "D" {
					t := e.defTempCat2(int(inst.Defs[0]), goType)
					w("%s = %s", t, call)
				} else {
					t := e.defTemp(int(inst.Defs[0]), goType)
					w("%s = %s", t, call)
				}
			}
		} else {
			// Bridge: runtime.InvokeStatic("class", "name", "desc", []rtda.Slot{...})
			call := fmt.Sprintf("runtime.InvokeStatic(%q, %q, %q, []rtda.Slot{%s})",
				className, name, desc, strings.Join(bridgeArgs, ", "))
			if md.ReturnType == "V" {
				w("%s", call)
			} else {
				goType, err := descToGo(md.ReturnType)
				if err != nil {
					return err
				}
				if md.ReturnType == "J" || md.ReturnType == "D" {
					t := e.defTempCat2(int(inst.Defs[0]), goType)
					w("%s = %s", t, slotExtract(call, md.ReturnType))
				} else {
					t := e.defTemp(int(inst.Defs[0]), goType)
					w("%s = %s", t, slotExtract(call, md.ReturnType))
				}
			}
		}

	// --- OOP: new / dup / invokespecial / getfield / putfield (A2.2b) ---
	case opcode.New:
		t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
		w("%s = runtime.NewObject(%q)", t, cp.ClassName(inst.Index))
	case opcode.Dup:
		// stack [v] -> [v, v]: the copy lands in the new top slot (Defs[1]); the
		// original (Defs[0]) keeps its temp.
		src := int(inst.Uses[0])
		t := e.defTemp(int(inst.Defs[1]), e.slotType[src])
		w("%s = %s", t, e.use(src))
	case opcode.Invokespecial:
		className, name, desc := cp.MemberRef(inst.Index)
		md := rtda.ParseMethodDescriptor(desc)
		args := []string{"rtda.RefSlot(" + e.use(int(inst.Uses[0])) + ")"} // this
		for i := 0; i < md.ArgSlots(); i++ {
			args = append(args, slotConstructor(md.ParameterTypes[i], e.use(int(inst.Uses[1+i]))))
		}
		call := fmt.Sprintf("runtime.InvokeSpecial(%q, %q, %q, []rtda.Slot{%s})",
			className, name, desc, strings.Join(args, ", "))
		if md.ReturnType == "V" {
			w("%s", call)
		} else {
			goType, err := descToGo(md.ReturnType)
			if err != nil {
				return err
			}
			t := e.defTemp(int(inst.Defs[0]), goType)
			w("%s = %s", t, slotExtract(call, md.ReturnType))
		}
	case opcode.Getfield:
		className, name, desc := cp.MemberRef(inst.Index)
		goType, err := descToGo(desc)
		if err != nil {
			return err
		}
		field := e.loader.LoadClass(className).LookupField(name, desc)
		obj := e.use(int(inst.Uses[0])) // read the objref before allocating the def (slot reuse)
		t := e.defTemp(int(inst.Defs[0]), goType)
		w("%s = %s", t, slotExtract(fmt.Sprintf("%s.Fields()[%d]", obj, field.SlotID()), desc))
	case opcode.Putfield:
		className, name, desc := cp.MemberRef(inst.Index)
		field := e.loader.LoadClass(className).LookupField(name, desc)
		obj, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		w("%s.Fields()[%d].%s(%s)", obj, field.SlotID(), setAccessor(desc), val)

	// --- float (category-1, float32) ---
	case opcode.Fload, opcode.Fload0, opcode.Fload1, opcode.Fload2, opcode.Fload3:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = %s", t, e.localName(inst))
	case opcode.Fstore, opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3:
		w("%s = %s", e.localName(inst), e.use(int(inst.Uses[0])))
	case opcode.Fconst0, opcode.Fconst1, opcode.Fconst2:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = %d", t, int(inst.Op-opcode.Fconst0))
	case opcode.Fadd, opcode.Fsub, opcode.Fmul, opcode.Fdiv:
		a, b := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = %s %s %s", t, a, binop(inst.Op), b)
	case opcode.Fneg:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = -%s", t, e.use(int(inst.Uses[0])))
	case opcode.Fcmpl, opcode.Fcmpg:
		a, bVal := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		emitCompare(b, t, a, bVal)
	case opcode.Freturn:
		w("return %s", e.use(int(inst.Uses[0])))
	case opcode.Faload:
		arr, idx := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = %s.GetFloatElement(int(%s))", t, arr, idx)
	case opcode.Fastore:
		arr, idx, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1])), e.use(int(inst.Uses[2]))
		w("%s.SetFloatElement(int(%s), %s)", arr, idx, val)

	// --- long (category-2, int64) ---
	case opcode.Lload, opcode.Lload0, opcode.Lload1, opcode.Lload2, opcode.Lload3:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = %s", t, e.localName(inst))
	case opcode.Lstore, opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3:
		w("%s = %s", e.localName(inst), e.use(int(inst.Uses[0])))
	case opcode.Lconst0, opcode.Lconst1:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = %d", t, int(inst.Op-opcode.Lconst0))
	case opcode.Ladd, opcode.Lsub, opcode.Lmul, opcode.Ldiv, opcode.Lrem,
		opcode.Land, opcode.Lor, opcode.Lxor:
		a, b := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = %s %s %s", t, a, binop(inst.Op), b)
	case opcode.Lneg:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = -%s", t, e.use(int(inst.Uses[0])))
	case opcode.Lshl, opcode.Lshr, opcode.Lushr:
		a, s := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = %s", t, longShiftExpr(inst.Op, a, s))
	case opcode.Lcmp:
		a, bVal := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		emitCompare(b, t, a, bVal)
	case opcode.Lreturn:
		w("return %s", e.use(int(inst.Uses[0])))
	case opcode.Laload:
		arr, idx := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = %s.GetLongElement(int(%s))", t, arr, idx)
	case opcode.Lastore:
		arr, idx, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1])), e.use(int(inst.Uses[2]))
		w("%s.SetLongElement(int(%s), %s)", arr, idx, val)

	// --- double (category-2, float64) ---
	case opcode.Dload, opcode.Dload0, opcode.Dload1, opcode.Dload2, opcode.Dload3:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = %s", t, e.localName(inst))
	case opcode.Dstore, opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		w("%s = %s", e.localName(inst), e.use(int(inst.Uses[0])))
	case opcode.Dconst0, opcode.Dconst1:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = %d", t, int(inst.Op-opcode.Dconst0))
	case opcode.Dadd, opcode.Dsub, opcode.Dmul, opcode.Ddiv:
		a, b := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = %s %s %s", t, a, binop(inst.Op), b)
	case opcode.Dneg:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = -%s", t, e.use(int(inst.Uses[0])))
	case opcode.Dcmpl, opcode.Dcmpg:
		a, bVal := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTemp(int(inst.Defs[0]), "int32")
		emitCompare(b, t, a, bVal)
	case opcode.Dreturn:
		w("return %s", e.use(int(inst.Uses[0])))
	case opcode.Daload:
		arr, idx := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = %s.GetDoubleElement(int(%s))", t, arr, idx)
	case opcode.Dastore:
		arr, idx, val := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1])), e.use(int(inst.Uses[2]))
		w("%s.SetDoubleElement(int(%s), %s)", arr, idx, val)

	// --- conversions ---
	case opcode.I2b:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(int8(%s))", t, e.use(int(inst.Uses[0])))
	case opcode.I2c:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(uint16(%s))", t, e.use(int(inst.Uses[0])))
	case opcode.I2s:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(int16(%s))", t, e.use(int(inst.Uses[0])))
	case opcode.I2l:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = int64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.I2f:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = float32(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.I2d:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = float64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.L2i:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.L2f:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = float32(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.L2d:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = float64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.F2i:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.F2l:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = int64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.F2d:
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = float64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.D2i:
		t := e.defTemp(int(inst.Defs[0]), "int32")
		w("%s = int32(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.D2l:
		t := e.defTempCat2(int(inst.Defs[0]), "int64")
		w("%s = int64(%s)", t, e.use(int(inst.Uses[0])))
	case opcode.D2f:
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = float32(%s)", t, e.use(int(inst.Uses[0])))

	// --- frem/drem (Go has no float %; use runtime helpers) ---
	case opcode.Frem:
		a, bVal := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[1]))
		t := e.defTemp(int(inst.Defs[0]), "float32")
		w("%s = runtime.FloatMod(%s, %s)", t, a, bVal)
	case opcode.Drem:
		a, bVal := e.use(int(inst.Uses[0])), e.use(int(inst.Uses[2]))
		t := e.defTempCat2(int(inst.Defs[0]), "float64")
		w("%s = runtime.DoubleMod(%s, %s)", t, a, bVal)

	// --- switch (tableswitch / lookupswitch → Go switch + goto) ---
	case opcode.Tableswitch:
		key := e.use(int(inst.Uses[0]))
		st := inst.Switch
		w("switch %s {", key)
		for i, target := range st.Targets {
			if e.merges[target] && len(e.mergeTemps[target]) > 0 {
				w("case %d:", int(st.Low)+i)
				e.emitMergeCopies(b, target)
				w("goto pc%d", target)
			} else {
				w("case %d: goto pc%d", int(st.Low)+i, target)
			}
		}
		w("default: goto pc%d", st.Default)
		w("}")
	case opcode.Lookupswitch:
		key := e.use(int(inst.Uses[0]))
		st := inst.Switch
		w("switch %s {", key)
		for i, target := range st.Targets {
			if e.merges[target] && len(e.mergeTemps[target]) > 0 {
				w("case %d:", int(st.Keys[i]))
				e.emitMergeCopies(b, target)
				w("goto pc%d", target)
			} else {
				w("case %d: goto pc%d", int(st.Keys[i]), target)
			}
		}
		w("default: goto pc%d", st.Default)
		w("}")

	default:
		return fmt.Errorf("opcode %s not supported", opcode.Name(inst.Op))
	}
	return nil
}

// localName returns the Go name for the local a load/store/iinc references.
// Uses the JVM-local→Go-param map (cat-2 aware: a double at slots 2-3 is one
// Go param "l1"); falls back to l{idx} for extra locals not in the map.
func (e *emitter) localName(inst *lowering.IRInst) string {
	if inst.Op == opcode.Iinc {
		return fmt.Sprintf("l%d", inst.IncIndex)
	}
	idx := int(inst.Index)
	switch inst.Op {
	case opcode.Iload0, opcode.Istore0, opcode.Aload0, opcode.Astore0,
		opcode.Lload0, opcode.Lstore0, opcode.Fload0, opcode.Fstore0,
		opcode.Dload0, opcode.Dstore0:
		idx = 0
	case opcode.Iload1, opcode.Istore1, opcode.Aload1, opcode.Astore1,
		opcode.Lload1, opcode.Lstore1, opcode.Fload1, opcode.Fstore1,
		opcode.Dload1, opcode.Dstore1:
		idx = 1
	case opcode.Iload2, opcode.Istore2, opcode.Aload2, opcode.Astore2,
		opcode.Lload2, opcode.Lstore2, opcode.Fload2, opcode.Fstore2,
		opcode.Dload2, opcode.Dstore2:
		idx = 2
	case opcode.Iload3, opcode.Istore3, opcode.Aload3, opcode.Astore3,
		opcode.Lload3, opcode.Lstore3, opcode.Fload3, opcode.Fstore3,
		opcode.Dload3, opcode.Dstore3:
		idx = 3
	}
	if name, ok := e.localMap[idx]; ok {
		return name
	}
	return fmt.Sprintf("l%d", idx)
}

// buildLocalMap maps JVM local slot indices to Go param names, accounting for
// category-2 params (long/double span 2 slots → one Go var).
func buildLocalMap(method *rtda.Method) map[int]string {
	md := rtda.ParseMethodDescriptor(method.Descriptor())
	m := map[int]string{}
	slot := 0
	paramIdx := 0
	if !method.IsStatic() {
		m[0] = "l0"
		slot = 1
		paramIdx = 1
	}
	for _, p := range md.ParameterTypes {
		sz := 1
		if p == "J" || p == "D" {
			sz = 2
		}
		name := fmt.Sprintf("l%d", paramIdx)
		for j := 0; j < sz; j++ {
			m[slot+j] = name
		}
		slot += sz
		paramIdx++
	}
	return m
}

// totalParamSlots returns the total JVM local slots consumed by the method's
// parameters (long/double = 2 each; +1 for `this` on instance methods).
func totalParamSlots(method *rtda.Method) int {
	md := rtda.ParseMethodDescriptor(method.Descriptor())
	n := 0
	for _, p := range md.ParameterTypes {
		if p == "J" || p == "D" {
			n += 2
		} else {
			n++
		}
	}
	if !method.IsStatic() {
		n++
	}
	return n
}

// --- signature / declarations / sink / terminator ---

func emitSignature(b *strings.Builder, method *rtda.Method, temps []tempDecl, localTypes map[int]string) {
	params, _ := paramGoTypes(method)
	ret, _ := returnGoType(method)
	fmt.Fprintf(b, "func %s(", mangle(method.Owner().Name(), method.Name()))
	for i, pt := range params {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "l%d %s", i, pt)
	}
	if ret == "" {
		b.WriteString(") {\n")
	} else {
		fmt.Fprintf(b, ") %s {\n", ret)
	}
	// Extra locals beyond params: type inferred from the store opcodes that write
	// them (astore → *rtda.Object, istore → int32); default int32.
	for i := totalParamSlots(method); i < int(method.MaxLocals()); i++ {
		gt := localTypes[i]
		if gt == "" {
			gt = "int32"
		}
		fmt.Fprintf(b, "\tvar l%d %s\n", i, gt)
	}
	// Temp declarations, grouped by Go type (all before any label → no goto-over-decl).
	for goType, names := range groupByType(temps) {
		fmt.Fprintf(b, "\tvar %s %s\n", strings.Join(names, ", "), goType)
	}
}

// localTypes infers each local's Go type from the store opcodes that write it
// (istore→int32, astore→*rtda.Object). Used to declare extra locals correctly.
func localTypes(ir *lowering.IR) map[int]string {
	types := map[int]string{}
	for i := range ir.Insts {
		inst := &ir.Insts[i]
		if !inst.Present {
			continue
		}
		if idx, gt, ok := storeLocalType(inst); ok {
			types[idx] = gt
		}
	}
	return types
}

func storeLocalType(inst *lowering.IRInst) (int, string, bool) {
	switch inst.Op {
	case opcode.Istore:
		return int(inst.Index), "int32", true
	case opcode.Astore:
		return int(inst.Index), "*rtda.Object", true
	case opcode.Lstore:
		return int(inst.Index), "int64", true
	case opcode.Fstore:
		return int(inst.Index), "float32", true
	case opcode.Dstore:
		return int(inst.Index), "float64", true
	case opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		return int(inst.Op - opcode.Istore0), "int32", true
	case opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		return int(inst.Op - opcode.Astore0), "*rtda.Object", true
	case opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3:
		return int(inst.Op - opcode.Lstore0), "int64", true
	case opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3:
		return int(inst.Op - opcode.Fstore0), "float32", true
	case opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		return int(inst.Op - opcode.Dstore0), "float64", true
	}
	return 0, "", false
}

func emitSink(b *strings.Builder, temps []tempDecl) {
	for _, td := range temps {
		fmt.Fprintf(b, "\t_ = %s\n", td.name)
	}
}

func emitTerminator(b *strings.Builder, method *rtda.Method) {
	ret, _ := returnGoType(method)
	if ret == "" {
		return // void: may fall off the end
	}
	fmt.Fprintf(b, "\treturn %s\n", zeroValue(ret))
}

func groupByType(temps []tempDecl) map[string][]string {
	g := map[string][]string{}
	for _, td := range temps {
		g[td.gotype] = append(g[td.gotype], td.name)
	}
	return g
}

// --- CFG analysis + merge temps (phi via copy-insertion) ---
//
// A merge (pc with >1 predecessor) with an empty operand stack (a loop head)
// needs no phi — loop state is in mutable locals. A merge with a value on the
// stack (a diamond) needs a phi: a per-slot merge temp, assigned at each
// predecessor edge, read after the join. cfgAnalysis finds merges + fall-through
// edges; allocMergeTemps allocates the merge temps (and refuses deferred types).

func cfgAnalysis(ir *lowering.IR) (merges, fallThrough map[int]bool) {
	merges = map[int]bool{}
	fallThrough = map[int]bool{}
	preds := map[int]int{}
	for pc := 0; pc < len(ir.Insts); pc++ {
		inst := &ir.Insts[pc]
		if !inst.Present {
			continue
		}
		for _, s := range successors(inst, pc) {
			preds[s]++
			if preds[s] > 1 {
				merges[s] = true
			}
		}
		if fallsThrough(inst.Op) {
			fallThrough[pc+inst.Length] = true
		}
	}
	return merges, fallThrough
}

// fallsThrough reports whether control reaches pc+length (i.e. the instruction is
// not an unconditional terminator).
func fallsThrough(op opcode.Opcode) bool {
	switch op {
	case opcode.Goto, opcode.GotoW, opcode.Return, opcode.Ireturn, opcode.Lreturn,
		opcode.Freturn, opcode.Dreturn, opcode.Areturn, opcode.Athrow,
		opcode.Tableswitch, opcode.Lookupswitch:
		return false
	}
	return true
}

func isBranch(op opcode.Opcode) bool {
	switch op {
	case opcode.Goto, opcode.GotoW,
		opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne,
		opcode.Ifnull, opcode.Ifnonnull:
		return true
	}
	return false
}

// allocMergeTemps allots one temp per stack slot at each non-empty-stack merge.
// Category-2 pairs (long/double = [Long,Top] / [Double,Top]) get ONE temp for
// both slots; standalone Top (unused) is skipped.
func (e *emitter) allocMergeTemps(ir *lowering.IR, method *rtda.Method) error {
	for pc := range e.merges {
		if pc >= len(ir.Insts) {
			continue
		}
		inst := &ir.Insts[pc]
		if !inst.Present || len(inst.InTypes) == 0 {
			continue // empty-stack merge (loop) — no phi needed
		}
		n := len(inst.InTypes)
		temps := make([]string, n)
		types := make([]string, n)
		k := 0
		for k < n {
			st := inst.InTypes[k]
			switch st {
			case lowering.TypeLong, lowering.TypeDouble:
				gt, _ := goTypeOf(st) // "int64" / "float64" (no error now)
				name := e.defMergeTemp(gt)
				temps[k] = name
				types[k] = gt
				if k+1 < n { // Top continuation → same temp
					temps[k+1] = name
					types[k+1] = gt
				}
				k += 2
			case lowering.TypeTop:
				k++ // standalone unused slot — skip
			default:
				gt, err := goTypeOf(st)
				if err != nil {
					return fmt.Errorf("transpile: %s: merge at pc %d slot %d: %w", method.Name(), pc, k, err)
				}
				name := e.defMergeTemp(gt)
				temps[k] = name
				types[k] = gt
				k++
			}
		}
		e.mergeTemps[pc] = temps
		e.mergeTempTypes[pc] = types
	}
	return nil
}

// defMergeTemp allocates a temp name recorded for declaration, without binding it
// to a stack slot (merge temps are written at predecessor edges, not by one def).
func (e *emitter) defMergeTemp(goType string) string {
	name := fmt.Sprintf("t%d", e.counter)
	e.counter++
	e.temps = append(e.temps, tempDecl{name, goType})
	return name
}

// emitMergeCopies writes `mergeTemp[k] = slotTemp[k]` for each of a merge's stack
// slots — the predecessor-edge side of the phi.
func (e *emitter) emitMergeCopies(b *strings.Builder, mergePc int) {
	for k, mt := range e.mergeTemps[mergePc] {
		fmt.Fprintf(b, "\t%s = %s\n", mt, e.slotTemp[k])
	}
}

// resetToMergeTemps makes post-join uses read the merge temps.
func (e *emitter) resetToMergeTemps(mergePc int) {
	for k, mt := range e.mergeTemps[mergePc] {
		e.slotTemp[k] = mt
		e.slotType[k] = e.mergeTempTypes[mergePc][k]
	}
}

// goTypeOf maps a lowering slot type to its Go type; long/float/double error.
func goTypeOf(st lowering.SlotType) (string, error) {
	switch st {
	case lowering.TypeInt:
		return "int32", nil
	case lowering.TypeRef:
		return "*rtda.Object", nil
	case lowering.TypeFloat:
		return "float32", nil
	case lowering.TypeLong:
		return "int64", nil
	case lowering.TypeDouble:
		return "float64", nil
	case lowering.TypeTop:
		return "", fmt.Errorf("unused (Top) merge slot")
	}
	return "", fmt.Errorf("unknown slot type %d", st)
}

func successors(inst *lowering.IRInst, pc int) []int {
	switch inst.Op {
	case opcode.Goto, opcode.GotoW:
		return []int{inst.Branch}
	case opcode.Return, opcode.Ireturn, opcode.Lreturn, opcode.Freturn,
		opcode.Dreturn, opcode.Areturn, opcode.Athrow:
		return nil
	case opcode.Tableswitch, opcode.Lookupswitch:
		out := []int{inst.Switch.Default}
		out = append(out, inst.Switch.Targets...)
		return out
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne,
		opcode.Ifnull, opcode.Ifnonnull:
		return []int{pc + inst.Length, inst.Branch}
	default:
		return []int{pc + inst.Length}
	}
}

// --- descriptor → Go type ---

func paramGoTypes(method *rtda.Method) ([]string, error) {
	md := rtda.ParseMethodDescriptor(method.Descriptor())
	out := make([]string, 0, len(md.ParameterTypes))
	for _, p := range md.ParameterTypes {
		gt, err := descToGo(p)
		if err != nil {
			return nil, err
		}
		out = append(out, gt)
	}
	return out, nil
}

func returnGoType(method *rtda.Method) (string, error) {
	ret := rtda.ParseMethodDescriptor(method.Descriptor()).ReturnType
	if ret == "V" || ret == "" {
		return "", nil
	}
	return descToGo(ret)
}

// descToGo maps a field/return/param descriptor to its Go type. A2.1 supports
// int (→ int32) and refs/arrays (→ *rtda.Object); long/float/double are errors.
func descToGo(desc string) (string, error) {
	switch desc {
	case "I", "B", "C", "S", "Z":
		return "int32", nil
	case "J":
		return "int64", nil
	case "F":
		return "float32", nil
	case "D":
		return "float64", nil
	case "V", "":
		return "", nil
	default: // L...; or [
		return "*rtda.Object", nil
	}
}

func zeroValue(goType string) string {
	if goType == "*rtda.Object" {
		return "nil"
	}
	return "0"
}

// --- opcode tables (carried over from A1) ---

func collectTargets(ir *lowering.IR) map[int]bool {
	targets := map[int]bool{}
	for i := range ir.Insts {
		inst := &ir.Insts[i]
		if !inst.Present {
			continue
		}
		switch inst.Op {
		case opcode.Goto, opcode.GotoW,
			opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
			opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
			opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne,
			opcode.Ifnull, opcode.Ifnonnull:
			targets[inst.Branch] = true
		case opcode.Tableswitch, opcode.Lookupswitch:
			targets[inst.Switch.Default] = true
			for _, t := range inst.Switch.Targets {
				targets[t] = true
			}
		}
	}
	return targets
}

func binop(op opcode.Opcode) string {
	switch op {
	case opcode.Iadd, opcode.Ladd, opcode.Fadd, opcode.Dadd:
		return "+"
	case opcode.Isub, opcode.Lsub, opcode.Fsub, opcode.Dsub:
		return "-"
	case opcode.Imul, opcode.Lmul, opcode.Fmul, opcode.Dmul:
		return "*"
	case opcode.Idiv, opcode.Ldiv, opcode.Fdiv, opcode.Ddiv:
		return "/"
	case opcode.Irem, opcode.Lrem:
		return "%"
	case opcode.Iand, opcode.Land:
		return "&"
	case opcode.Ior, opcode.Lor:
		return "|"
	case opcode.Ixor, opcode.Lxor:
		return "^"
	}
	return "?"
}

func shiftExpr(op opcode.Opcode, v, amount string) string {
	switch op {
	case opcode.Ishl:
		return fmt.Sprintf("%s << (%s & 31)", v, amount)
	case opcode.Ishr:
		return fmt.Sprintf("%s >> (%s & 31)", v, amount)
	case opcode.Iushr:
		return fmt.Sprintf("int32(uint32(%s) >> (%s & 31))", v, amount)
	}
	return v
}

// longShiftExpr renders a long shift with the JVM's 6-bit mask.
func longShiftExpr(op opcode.Opcode, v, amount string) string {
	switch op {
	case opcode.Lshl:
		return fmt.Sprintf("%s << (%s & 63)", v, amount)
	case opcode.Lshr:
		return fmt.Sprintf("%s >> (%s & 63)", v, amount)
	case opcode.Lushr:
		return fmt.Sprintf("int64(uint64(%s) >> (%s & 63))", v, amount)
	}
	return v
}

// emitCompare writes the 3-line compare (result = 0; if a>b result=1; if a<b
// result=-1) — Go has no ternary. Used by lcmp/fcmp/dcmp.
func emitCompare(b *strings.Builder, result, a, bVal string) {
	fmt.Fprintf(b, "\t%s = 0\n", result)
	fmt.Fprintf(b, "\tif %s > %s { %s = 1 }\n", a, bVal, result)
	fmt.Fprintf(b, "\tif %s < %s { %s = -1 }\n", a, bVal, result)
}

func cmp0(op opcode.Opcode) string {
	switch op {
	case opcode.Ifeq:
		return "=="
	case opcode.Ifne:
		return "!="
	case opcode.Iflt:
		return "<"
	case opcode.Ifge:
		return ">="
	case opcode.Ifgt:
		return ">"
	case opcode.Ifle:
		return "<="
	}
	return "?"
}

func icmp(op opcode.Opcode) string {
	switch op {
	case opcode.IfIcmpeq, opcode.IfAcmpeq:
		return "=="
	case opcode.IfIcmpne, opcode.IfAcmpne:
		return "!="
	case opcode.IfIcmplt:
		return "<"
	case opcode.IfIcmpge:
		return ">="
	case opcode.IfIcmpgt:
		return ">"
	case opcode.IfIcmple:
		return "<="
	}
	return "?"
}

// slotConstructor wraps a typed temp in the rtda slot constructor for its
// descriptor (ref → RefSlot, int → IntSlot), for boxing invoke args.
func slotConstructor(desc, temp string) string {
	if isRefDesc(desc) {
		return "rtda.RefSlot(" + temp + ")"
	}
	return "rtda.IntSlot(" + temp + ")"
}

// slotExtract extracts a typed value from a Slot-bearing expression.
func slotExtract(call, desc string) string {
	if isRefDesc(desc) {
		return call + ".Ref()"
	}
	return call + ".Num()"
}

// setAccessor maps a field descriptor to the Slot setter (SetNum/SetRef).
func setAccessor(desc string) string {
	if isRefDesc(desc) {
		return "SetRef"
	}
	return "SetNum"
}

// isRefDesc reports whether a descriptor is an object/array reference.
func isRefDesc(desc string) bool {
	return strings.HasPrefix(desc, "L") || strings.HasPrefix(desc, "[")
}
