// Package aot lowers verified JVM bytecode methods to real Go function bodies.
// The generated body owns its control flow: branches become Go branches over
// basic blocks and each supported instruction becomes specialized Go code
// calling shared runtime helpers.
package aot

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// Generate emits one top-level Go function with the common entry signature.
// heights are the verified operand-stack depths at each instruction offset.
func Generate(symbol string, m *classfile.Method, heights map[int32]int) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "func %s(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {\n", symbol)
	fmt.Fprintf(&b, "\tvar locals [%d]rt.Value\n", m.MaxLocals)
	fmt.Fprintf(&b, "\tvar stack [%d]rt.Value\n", m.MaxStack)
	b.WriteString("\t_ = locals\n\t_ = stack\n")
	if m.Flags&rt.AccStatic == 0 {
		b.WriteString("\tlocals[0] = rt.RefValue(receiver)\n")
	}
	mt, err := rt.ParseMethodType(m.Descriptor)
	if err != nil {
		return "", err
	}
	base := 0
	if m.Flags&rt.AccStatic == 0 {
		base = 1
	}
	for i := range mt.Params {
		fmt.Fprintf(&b, "\tlocals[%d] = args[%d]\n", base+i, i)
	}
	b.WriteString("\tblock := 0\n\tfor {\n\t\tswitch block {\n")

	blocks, blockOf, err := buildBlocks(m, heights)
	if err != nil {
		return "", err
	}
	for _, blk := range blocks {
		fmt.Fprintf(&b, "\t\tcase %d:\n", blk.id)
		h := heights[blk.start]
		terminated := false
		for i := blk.first; i <= blk.last; i++ {
			insn := m.Insns[i]
			code, nextH, term, err := emitInsn(insn, m, h, heights, blockOf)
			if err != nil {
				return "", fmt.Errorf("%s.%s%s pc %d: %w", m.Owner, m.Name, m.Descriptor, insn.PC, err)
			}
			for _, line := range strings.Split(strings.TrimSuffix(code, "\n"), "\n") {
				if line == "" {
					continue
				}
				fmt.Fprintf(&b, "\t\t\t%s\n", line)
			}
			h = nextH
			if term {
				terminated = true
			}
		}
		if !terminated {
			next := blk.endPC
			id, ok := blockOf[next]
			if !ok {
				return "", fmt.Errorf("method %s.%s%s: block at pc %d falls through to unknown block %d",
					m.Owner, m.Name, m.Descriptor, blk.start, next)
			}
			fmt.Fprintf(&b, "\t\t\tblock = %d\n\t\t\tcontinue\n", id)
		}
	}
	b.WriteString("\t\t}\n\t}\n}\n")
	return b.String(), nil
}

type block struct {
	id    int
	start int32
	endPC int32
	first int
	last  int
	code  string
}

func buildBlocks(m *classfile.Method, heights map[int32]int) ([]block, map[int32]int, error) {
	leaders := map[int32]bool{0: true}
	index := map[int32]int{}
	for i, insn := range m.Insns {
		index[insn.PC] = i
	}
	for _, insn := range m.Insns {
		if target, ok := classfile.BranchTarget(insn); ok {
			leaders[target] = true
			switch insn.Op {
			case rt.OpGoto:
			default:
				next := insn.PC + insn.Len
				if int(next) < m.CodeLength {
					leaders[next] = true
				}
			}
		}
		switch insn.Op {
		case rt.OpReturn, rt.OpIreturn, rt.OpAreturn, rt.OpAthrow, rt.OpGoto:
			next := insn.PC + insn.Len
			if int(next) < m.CodeLength {
				leaders[next] = true
			}
		}
	}
	for _, h := range m.Handlers {
		leaders[h.HandlerPC] = true
	}
	var starts []int32
	for pc := range leaders {
		starts = append(starts, pc)
	}
	sort.Slice(starts, func(i, j int) bool { return starts[i] < starts[j] })
	var blocks []block
	blockOf := map[int32]int{}
	for i, start := range starts {
		end := int32(m.CodeLength)
		if i+1 < len(starts) {
			end = starts[i+1]
		}
		if _, ok := heights[start]; !ok {
			continue
		}
		first, ok := index[start]
		if !ok {
			return nil, nil, fmt.Errorf("block start %d is not an instruction", start)
		}
		last := first
		for last+1 < len(m.Insns) && m.Insns[last+1].PC < end {
			last++
		}
		id := len(blocks)
		blocks = append(blocks, block{id: id, start: start, endPC: end, first: first, last: last})
		blockOf[start] = id
	}
	if _, ok := blockOf[0]; !ok {
		return nil, nil, fmt.Errorf("entry block is unreachable")
	}
	return blocks, blockOf, nil
}

func emitInsn(insn rt.Insn, m *classfile.Method, h int, heights map[int32]int, blockOf map[int32]int) (string, int, bool, error) {
	pc := insn.PC
	switch insn.Op {
	case rt.OpAconstNull:
		return fmt.Sprintf("stack[%d] = rt.RefValue(nil)", h), h + 1, false, nil
	case rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2, rt.OpIconst3, rt.OpIconst4, rt.OpIconst5:
		return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, int32(insn.Op)-3), h + 1, false, nil
	case rt.OpBipush, rt.OpSipush:
		return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, insn.A), h + 1, false, nil
	case rt.OpLdc, rt.OpLdcW:
		ref := m.Refs[insn.Ref]
		switch ref.Kind {
		case rt.RefInt:
			return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, ref.Int), h + 1, false, nil
		case rt.RefClass:
			return fmt.Sprintf("stack[%d] = t.ClassConstant(%s)", h, strconv.Quote(ref.Class)), h + 1, false, nil
		case rt.RefString:
			var units []string
			for _, u := range ref.Units {
				units = append(units, strconv.Itoa(int(u)))
			}
			return fmt.Sprintf("stack[%d] = t.StringConstant([]uint16{%s})", h, strings.Join(units, ", ")), h + 1, false, nil
		default:
			return "", h, false, fmt.Errorf("unsupported ldc reference")
		}
	case rt.OpIload, rt.OpAload:
		return fmt.Sprintf("stack[%d] = locals[%d]", h, insn.A), h + 1, false, nil
	case 0x1a, 0x1b, 0x1c, 0x1d:
		return fmt.Sprintf("stack[%d] = locals[%d]", h, insn.Op-0x1a), h + 1, false, nil
	case 0x2a, 0x2b, 0x2c, 0x2d:
		return fmt.Sprintf("stack[%d] = locals[%d]", h, insn.Op-0x2a), h + 1, false, nil
	case rt.OpIstore, rt.OpAstore:
		return fmt.Sprintf("locals[%d] = stack[%d]", insn.A, h-1), h - 1, false, nil
	case 0x3b, 0x3c, 0x3d, 0x3e:
		return fmt.Sprintf("locals[%d] = stack[%d]", insn.Op-0x3b, h-1), h - 1, false, nil
	case 0x4b, 0x4c, 0x4d, 0x4e:
		return fmt.Sprintf("locals[%d] = stack[%d]", insn.Op-0x4b, h-1), h - 1, false, nil
	case rt.OpPop:
		return "", h - 1, false, nil
	case rt.OpDup:
		return fmt.Sprintf("stack[%d] = stack[%d]", h, h-1), h + 1, false, nil
	case rt.OpNew:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nobj, thrown := t.NewObject(%s)\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(obj)\n}",
			strconv.Quote(ref.Class), transferCode(pc, m, blockOf), h)
		return code, h + 1, false, nil
	case rt.OpGetfield:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nv, thrown := t.GetField(stack[%d], %s, %s, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			h-1, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor),
			transferCode(pc, m, blockOf), h-1)
		return code, h, false, nil
	case rt.OpPutfield:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nthrown := t.PutField(stack[%d], %s, %s, %s, stack[%d])\nif thrown != nil {\n%s\n}\n}",
			h-2, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor),
			h-1, transferCode(pc, m, blockOf))
		return code, h - 2, false, nil
	case rt.OpGetstatic:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nv, thrown := t.GetStatic(%s, %s, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor),
			transferCode(pc, m, blockOf), h)
		return code, h + 1, false, nil
	case rt.OpPutstatic:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nthrown := t.PutStatic(%s, %s, %s, stack[%d])\nif thrown != nil {\n%s\n}\n}",
			strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor),
			h-1, transferCode(pc, m, blockOf))
		return code, h - 1, false, nil
	case rt.OpAnewarray:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\narr, thrown := t.NewArray(rt.AsInt(stack[%d]), %s)\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(arr)\n}",
			h-1, strconv.Quote(ref.Class), transferCode(pc, m, blockOf), h-1)
		return code, h, false, nil
	case rt.OpAaload:
		code := fmt.Sprintf("{\nv, thrown := t.ArrayLoad(stack[%d], rt.AsInt(stack[%d]))\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			h-2, h-1, transferCode(pc, m, blockOf), h-2)
		return code, h - 1, false, nil
	case rt.OpAastore:
		code := fmt.Sprintf("{\nthrown := t.ArrayStore(stack[%d], rt.AsInt(stack[%d]), stack[%d])\nif thrown != nil {\n%s\n}\n}",
			h-3, h-2, h-1, transferCode(pc, m, blockOf))
		return code, h - 3, false, nil
	case rt.OpArraylength:
		code := fmt.Sprintf("{\nn, thrown := t.ArrayLength(stack[%d])\nif thrown != nil {\n%s\n}\nstack[%d] = rt.IntValue(n)\n}",
			h-1, transferCode(pc, m, blockOf), h-1)
		return code, h, false, nil
	case rt.OpIadd:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) + rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIrem:
		return "", h, false, fmt.Errorf("irem has no AOT lowering (pc %d, bci %d)", pc, pc)
	case rt.OpInvokestatic, rt.OpInvokevirtual, rt.OpInvokespecial, rt.OpInvokeinterface:
		ref := m.Refs[insn.Ref]
		mt, err := rt.ParseMethodType(ref.Descriptor)
		if err != nil {
			return "", h, false, err
		}
		n := mt.ArgCount
		args := fmt.Sprintf("stack[%d:%d]", h-n, h)
		var call string
		var newH int
		switch insn.Op {
		case rt.OpInvokestatic:
			call = fmt.Sprintf("t.InvokeStatic(m, %d, %s, %s, %s, %s)", pc,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH = h - n
		case rt.OpInvokespecial:
			call = fmt.Sprintf("t.InvokeSpecial(m, %d, stack[%d], %s, %s, %s, %s)", pc, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH = h - n - 1
		case rt.OpInvokevirtual:
			call = fmt.Sprintf("t.InvokeVirtual(m, %d, stack[%d], %s, %s, %s, %s)", pc, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH = h - n - 1
		case rt.OpInvokeinterface:
			call = fmt.Sprintf("t.InvokeInterface(m, %d, stack[%d], %s, %s, %s, %s)", pc, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH = h - n - 1
		}
		var body strings.Builder
		body.WriteString("{\nout := " + call + "\nif out.Thrown != nil {\nthrown := out.Thrown\n")
		body.WriteString(transferCode(pc, m, blockOf))
		body.WriteString("\n}\n")
		if mt.Return.Kind != 'V' {
			fmt.Fprintf(&body, "stack[%d] = out.Value\n", newH)
			newH++
		}
		body.WriteString("}")
		return body.String(), newH, false, nil
	case rt.OpReturn:
		return "return rt.Normal()", h, true, nil
	case rt.OpIreturn, rt.OpAreturn:
		return fmt.Sprintf("return rt.OkValue(stack[%d])", h-1), h - 1, true, nil
	case rt.OpAthrow:
		code := fmt.Sprintf("{\nthrown := stack[%d].Ref\nif thrown == nil {\nthrown = t.Throwable(\"java/lang/NullPointerException\")\n}\n%s\n}",
			h-1, transferCode(pc, m, blockOf))
		return code, h - 1, true, nil
	case rt.OpIfeq, rt.OpIfne:
		cond := fmt.Sprintf("rt.AsInt(stack[%d]) == 0", h-1)
		if insn.Op == rt.OpIfne {
			cond = fmt.Sprintf("rt.AsInt(stack[%d]) != 0", h-1)
		}
		target, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		fall := insn.PC + insn.Len
		fallID, err := blockID(fall, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fallID), h - 1, true, nil
	case rt.OpIfIcmpeq, rt.OpIfIcmpne:
		op := "=="
		if insn.Op == rt.OpIfIcmpne {
			op = "!="
		}
		cond := fmt.Sprintf("rt.AsInt(stack[%d]) %s rt.AsInt(stack[%d])", h-2, op, h-1)
		target, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		fallID, err := blockID(insn.PC+insn.Len, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fallID), h - 2, true, nil
	case rt.OpIfAcmpeq, rt.OpIfAcmpne:
		op := "=="
		if insn.Op == rt.OpIfAcmpne {
			op = "!="
		}
		cond := fmt.Sprintf("stack[%d].Ref %s stack[%d].Ref", h-2, op, h-1)
		target, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		fallID, err := blockID(insn.PC+insn.Len, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fallID), h - 2, true, nil
	case rt.OpGoto:
		id, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("block = %d\ncontinue", id), h, true, nil
	default:
		return "", h, false, fmt.Errorf("unsupported opcode 0x%02x", insn.Op)
	}
}

func blockID(pc int32, blockOf map[int32]int) (int, error) {
	id, ok := blockOf[pc]
	if !ok {
		return 0, fmt.Errorf("no generated block at pc %d", pc)
	}
	return id, nil
}

func transferCode(pc int32, m *classfile.Method, blockOf map[int32]int) string {
	if len(m.Handlers) == 0 {
		return "return rt.ThrowOutcome(thrown)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "if hpc, ok := t.HandleThrow(m, %d, thrown); ok {\n", pc)
	fmt.Fprintf(&b, "stack = [%d]rt.Value{}\n", m.MaxStack)
	b.WriteString("stack[0] = rt.RefValue(thrown)\n")
	b.WriteString("switch hpc {\n")
	seen := map[int32]bool{}
	for _, h := range m.Handlers {
		if seen[h.HandlerPC] {
			continue
		}
		seen[h.HandlerPC] = true
		id, ok := blockOf[h.HandlerPC]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "case %d:\nblock = %d\n", h.HandlerPC, id)
	}
	b.WriteString("default:\npanic(\"catty: generated handler transfer has no target block\")\n}\n")
	b.WriteString("continue\n}\n")
	b.WriteString("return rt.ThrowOutcome(thrown)")
	return b.String()
}
