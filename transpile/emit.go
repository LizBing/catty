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
func Emit(method *rtda.Method, ir *lowering.IR, loader rtda.Loader) (string, error) {
	if !method.IsStatic() {
		return "", fmt.Errorf("transpile: A2.1 supports static methods only (got %s)", method.Name())
	}
	if err := stackMergeError(method, ir); err != nil {
		return "", err
	}
	cp := method.Owner().ConstantPool()
	targets := collectTargets(ir)

	e := &emitter{slotTemp: map[int]string{}, slotType: map[int]string{}, loader: loader}
	var body strings.Builder
	for pc := 0; pc < len(ir.Insts); pc++ {
		inst := &ir.Insts[pc]
		if !inst.Present {
			continue
		}
		if targets[pc] {
			fmt.Fprintf(&body, "pc%d:\n", pc)
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
	slotTemp map[int]string // slot index → temp name of its most recent def
	slotType map[int]string // slot index → that temp's Go type (for dup)
	loader   rtda.Loader
	temps    []tempDecl
	counter  int
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
		case classfile.ConstantString:
			t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
			w("%s = runtime.NewString(%q)", t, cp.String(inst.Index))
		default:
			return fmt.Errorf("ldc: A2.2 supports int/String constants only (tag %d)", cp.Tag(inst.Index))
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
		w("%s = %s", t, localName(inst))
	case opcode.Aload, opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		t := e.defTemp(int(inst.Defs[0]), "*rtda.Object")
		w("%s = %s", t, localName(inst))

	// --- stores ---
	case opcode.Istore, opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		w("%s = %s", localName(inst), e.use(int(inst.Uses[0])))
	case opcode.Astore, opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		w("%s = %s", localName(inst), e.use(int(inst.Uses[0])))

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
		w("%s += %d", localName(inst), int32(inst.Const8))

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
		args := make([]string, md.ArgSlots())
		for i := range args {
			args[i] = e.use(int(inst.Uses[i]))
		}
		call := fmt.Sprintf("%s(%s)", mangle(className, name), strings.Join(args, ", "))
		if md.ReturnType == "V" {
			w("%s", call)
		} else {
			goType, err := descToGo(md.ReturnType)
			if err != nil {
				return err
			}
			t := e.defTemp(int(inst.Defs[0]), goType)
			w("%s = %s", t, call)
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

	default:
		return fmt.Errorf("opcode %s not supported (A2.2b: int/ref/arrays, loops, OOP)", opcode.Name(inst.Op))
	}
	return nil
}

// localName returns the Go name for the local an iload/istore/iinc references.
func localName(inst *lowering.IRInst) string {
	if inst.Op == opcode.Iinc {
		return fmt.Sprintf("l%d", inst.IncIndex) // iinc's local index lives in IncIndex, not Index
	}
	idx := int(inst.Index)
	switch inst.Op {
	case opcode.Iload0, opcode.Istore0:
		idx = 0
	case opcode.Iload1, opcode.Istore1:
		idx = 1
	case opcode.Iload2, opcode.Istore2:
		idx = 2
	case opcode.Iload3, opcode.Istore3:
		idx = 3
	case opcode.Aload0, opcode.Astore0:
		idx = 0
	case opcode.Aload1, opcode.Astore1:
		idx = 1
	case opcode.Aload2, opcode.Astore2:
		idx = 2
	case opcode.Aload3, opcode.Astore3:
		idx = 3
	}
	return fmt.Sprintf("l%d", idx)
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
	for i := len(params); i < int(method.MaxLocals()); i++ {
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
	case opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		return int(inst.Op - opcode.Istore0), "int32", true
	case opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		return int(inst.Op - opcode.Astore0), "*rtda.Object", true
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

// --- merge gate ---
//
// A merge (a pc with >1 predecessor) is fine as long as no operand-stack value
// crosses it: loop state lives in mutable locals (no phi), and the operand stack
// is empty at loop heads. Only a merge with a value on the stack (a diamond like
// `cond ? x : y`) needs phi insertion — those are refused here (a later step).

func stackMergeError(method *rtda.Method, ir *lowering.IR) error {
	preds := map[int]int{}
	for pc := 0; pc < len(ir.Insts); pc++ {
		inst := &ir.Insts[pc]
		if !inst.Present {
			continue
		}
		for _, s := range successors(inst, pc) {
			preds[s]++
		}
	}
	for pc, n := range preds {
		if n < 2 || pc >= len(ir.Insts) {
			continue
		}
		inst := &ir.Insts[pc]
		if inst.Present && len(inst.InTypes) > 0 {
			return fmt.Errorf("transpile: %s: merge at pc %d has a value on the operand stack (phi needed — diamonds are a later step)",
				method.Name(), pc)
		}
	}
	return nil
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
		return "", fmt.Errorf("long not supported in A2.1")
	case "F":
		return "", fmt.Errorf("float not supported in A2.1")
	case "D":
		return "", fmt.Errorf("double not supported in A2.1")
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
	case opcode.Iadd:
		return "+"
	case opcode.Isub:
		return "-"
	case opcode.Imul:
		return "*"
	case opcode.Idiv:
		return "/"
	case opcode.Irem:
		return "%"
	case opcode.Iand:
		return "&"
	case opcode.Ior:
		return "|"
	case opcode.Ixor:
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
