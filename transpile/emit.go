// Package transpile lowers a method's IR to Go source — the AOT transpiler's
// emitter. This is the first executable proof of the "emit Go, let go build
// optimize" thesis (ROADMAP A1): a method compiled to native code via the Go
// toolchain instead of interpreted.
//
// A1 scope: int-only, static methods (see Emit for the supported opcode subset).
// Slots and locals become Go locals; bytecode control flow becomes goto/labels.
package transpile

import (
	"fmt"
	"strings"

	"catty/classfile"
	"catty/lowering"
	"catty/opcode"
	"catty/rtda"
)

// Emit turns one method's IR into a Go function definition. The operand stack is
// eliminated: each stack slot is a Go local `sK`, each JVM local is `lK` (the
// first ArgSlotCount of which are the function's parameters). Control flow uses
// Go goto/labels. A1 is int-only — any non-int opcode returns an error.
//
// The Go-source rules that shape the emitter:
//   - all slot/extra-local declarations come before any label, so goto never
//     crosses a var declaration;
//   - a `pcNN:` label is emitted only at branch/switch targets (no unused labels);
//   - a trailing `_ = sK` sink marks every slot used (no unused-local errors).
func Emit(method *rtda.Method, ir *lowering.IR) (string, error) {
	if !method.IsStatic() {
		return "", fmt.Errorf("transpile: A1 supports static methods only (got instance method %s)", method.Name())
	}
	cp := method.Owner().ConstantPool()
	targets := collectTargets(ir)

	var b strings.Builder
	emitSignature(&b, method)
	var body strings.Builder
	if err := emitInst(&body, ir, cp, targets); err != nil {
		return "", err
	}
	b.WriteString(body.String())
	emitSink(&b, method)
	// A non-void function must end with a terminating statement. The body already
	// returns on every path; this trailing return is unreachable but satisfies
	// Go's "missing return" check (the sink statements above are not terminating).
	if method.ReturnType() != "V" {
		b.WriteString("\treturn 0\n")
	}
	b.WriteString("}\n")
	return b.String(), nil
}

// collectTargets returns the set of pcs that are explicit branch/switch targets
// (the only places that need a Go label).
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
			opcode.IfIcmpgt, opcode.IfIcmple:
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

// emitSignature writes `func <mangled>(l0 int32, ...) int32 {` plus the slot and
// extra-local declarations.
func emitSignature(b *strings.Builder, method *rtda.Method) {
	argSlots := int(method.ArgSlotCount())
	maxStack := int(method.MaxStack())
	maxLocals := int(method.MaxLocals())

	fmt.Fprintf(b, "func %s(", mangle(method.Owner().Name(), method.Name()))
	for i := 0; i < argSlots; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "l%d int32", i)
	}
	b.WriteString(") int32 {\n")

	if maxStack > 0 {
		b.WriteString("\tvar ")
		writeIndexed(b, "s", maxStack)
		b.WriteString(" int32\n")
	}
	if extra := maxLocals - argSlots; extra > 0 {
		b.WriteString("\tvar ")
		writeIndexedFrom(b, "l", argSlots, maxLocals)
		b.WriteString(" int32\n")
	}
}

// emitInst writes the body: one or more Go statements per IR instruction, with a
// `pcNN:` label at branch targets. Returns an error (with pc context) on any
// unsupported opcode, so Emit never emits wrong code.
func emitInst(b *strings.Builder, ir *lowering.IR, cp *classfile.ConstantPool, targets map[int]bool) error {
	for pc := 0; pc < len(ir.Insts); pc++ {
		inst := &ir.Insts[pc]
		if !inst.Present {
			continue
		}
		if targets[pc] {
			fmt.Fprintf(b, "pc%d:\n", pc)
		}
		if err := emitOne(b, inst, cp); err != nil {
			return fmt.Errorf("at pc %d: %w", pc, err)
		}
	}
	return nil
}

func emitOne(b *strings.Builder, inst *lowering.IRInst, cp *classfile.ConstantPool) error {
	s := func(i uint8) string { return fmt.Sprintf("s%d", i) }
	u0, u1 := func() string { return s(inst.Uses[0]) }, func() string { return s(inst.Uses[1]) }
	d0 := func() string { return s(inst.Defs[0]) }
	w := func(format string, args ...any) { fmt.Fprintf(b, "\t"+format+"\n", args...) }

	switch inst.Op {
	// --- int loads ---
	case opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3:
		w("%s = l%d", d0(), int(inst.Op-opcode.Iload0))
	case opcode.Iload:
		w("%s = l%d", d0(), inst.Index)

	// --- int constants (iconst_m1..5 collapse via op - Iconst0) ---
	case opcode.IconstM1, opcode.Iconst0, opcode.Iconst1, opcode.Iconst2,
		opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		w("%s = %d", d0(), int32(inst.Op-opcode.Iconst0))
	case opcode.Bipush:
		w("%s = %d", d0(), int32(inst.Const8))
	case opcode.Sipush:
		w("%s = %d", d0(), int32(inst.Const16))
	case opcode.Ldc, opcode.LdcW:
		if cp.Tag(inst.Index) != classfile.ConstantInteger {
			return fmt.Errorf("ldc: A1 supports int constants only (tag %d)", cp.Tag(inst.Index))
		}
		w("%s = %d", d0(), cp.Integer(inst.Index))

	// --- int stores ---
	case opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		w("l%d = %s", int(inst.Op-opcode.Istore0), u0())
	case opcode.Istore:
		w("l%d = %s", inst.Index, u0())

	// --- int arithmetic ---
	case opcode.Iadd, opcode.Isub, opcode.Imul, opcode.Idiv, opcode.Irem,
		opcode.Iand, opcode.Ior, opcode.Ixor:
		w("%s = %s %s %s", d0(), u0(), binop(inst.Op), u1())
	case opcode.Ineg:
		w("%s = -%s", d0(), u0())
	case opcode.Ishl, opcode.Ishr, opcode.Iushr:
		w("%s = %s", d0(), shiftExpr(inst.Op, u0(), u1()))
	case opcode.Iinc:
		w("l%d += %d", inst.IncIndex, int32(inst.Const8))

	// --- conditional branches ---
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle:
		w("if %s %s 0 { goto pc%d }", u0(), cmp0(inst.Op), inst.Branch)
	case opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple:
		w("if %s %s %s { goto pc%d }", u0(), icmp(inst.Op), u1(), inst.Branch)

	case opcode.Goto, opcode.GotoW:
		w("goto pc%d", inst.Branch)

	case opcode.Ireturn:
		w("return %s", u0())

	// --- invokestatic (recursive/cross-method calls resolve by mangled name) ---
	case opcode.Invokestatic:
		className, name, desc := cp.MemberRef(inst.Index)
		md := rtda.ParseMethodDescriptor(desc)
		argSlots := md.ArgSlots()
		args := make([]string, argSlots)
		for i := 0; i < argSlots; i++ {
			args[i] = s(inst.Uses[i])
		}
		call := fmt.Sprintf("%s(%s)", mangle(className, name), strings.Join(args, ", "))
		if md.ReturnType == "V" {
			w("%s", call)
		} else {
			w("%s = %s", d0(), call)
		}

	default:
		return fmt.Errorf("opcode %s not supported in A1 (int-only)", opcode.Name(inst.Op))
	}
	return nil
}

// emitSink writes unreachable `_ = sK` / `_ = lK` reads so every declared slot
// and extra local counts as used (Go otherwise errors on unused locals).
func emitSink(b *strings.Builder, method *rtda.Method) {
	argSlots := int(method.ArgSlotCount())
	for i := 0; i < int(method.MaxStack()); i++ {
		fmt.Fprintf(b, "\t_ = s%d\n", i)
	}
	for i := argSlots; i < int(method.MaxLocals()); i++ {
		fmt.Fprintf(b, "\t_ = l%d\n", i)
	}
}

// binop maps an int binary opcode to its Go operator.
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
	return "?" // unreachable; guarded by the caller's case list
}

// shiftExpr renders a shift with the JVM's 5-bit mask.
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

// cmp0 / icmp map branch opcodes to Go comparison operators.
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
	case opcode.IfIcmpeq:
		return "=="
	case opcode.IfIcmpne:
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

// writeIndexed writes "s0, s1, ..., s{n-1}".
func writeIndexed(b *strings.Builder, prefix string, n int) {
	writeIndexedFrom(b, prefix, 0, n)
}

func writeIndexedFrom(b *strings.Builder, prefix string, from, to int) {
	for i := from; i < to; i++ {
		if i > from {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "%s%d", prefix, i)
	}
}
