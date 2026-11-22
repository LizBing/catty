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
// heights are the verified logical operand-stack depths and shapes the logical
// value widths at each instruction offset.
func Generate(symbol string, m *classfile.Method, heights map[int32]int, shapes map[int32][]int) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "func %s(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {\n", symbol)
	fmt.Fprintf(&b, "\tvar locals [%d]rt.Value\n", m.MaxLocals)
	fmt.Fprintf(&b, "\tvar stack [%d]rt.Value\n", m.MaxStack)
	b.WriteString("\t_ = locals\n\t_ = stack\n")
	mt, err := rt.ParseMethodType(m.Descriptor)
	if err != nil {
		return "", err
	}
	base := 0
	if m.Flags&rt.AccStatic == 0 {
		b.WriteString("\tlocals[0] = rt.RefValue(receiver)\n")
		base = 1
	}
	for i := range mt.Params {
		fmt.Fprintf(&b, "\tlocals[%d] = args[%d]\n", base+mt.ParamSlot[i], i)
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
			widths := shapes[insn.PC]
			code, nextH, term, err := emitInsn(insn, m, h, widths, heights, blockOf)
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
			case rt.OpGoto, rt.OpGotoW:
			default:
				next := insn.PC + insn.Len
				if int(next) < m.CodeLength {
					leaders[next] = true
				}
			}
		}
		if _, targets, ok := classfile.SwitchTargets(insn); ok {
			for _, target := range targets {
				leaders[target] = true
			}
			leaders[insn.A] = true
			next := insn.PC + insn.Len
			if int(next) < m.CodeLength {
				leaders[next] = true
			}
		}
		switch insn.Op {
		case rt.OpReturn, rt.OpIreturn, rt.OpLreturn, rt.OpFreturn, rt.OpDreturn,
			rt.OpAreturn, rt.OpAthrow, rt.OpGoto, rt.OpGotoW:
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

func blockID(pc int32, blockOf map[int32]int) (int, error) {
	id, ok := blockOf[pc]
	if !ok {
		return 0, fmt.Errorf("no generated block at pc %d", pc)
	}
	return id, nil
}

func booleanCode(op uint8, a, b string) string {
	switch op {
	case rt.OpIfeq:
		return a + " == 0"
	case rt.OpIfne:
		return a + " != 0"
	case rt.OpIflt:
		return a + " < 0"
	case rt.OpIfge:
		return a + " >= 0"
	case rt.OpIfgt:
		return a + " > 0"
	case rt.OpIfle:
		return a + " <= 0"
	case rt.OpIfIcmpeq:
		return a + " == " + b
	case rt.OpIfIcmpne:
		return a + " != " + b
	case rt.OpIfIcmplt:
		return a + " < " + b
	case rt.OpIfIcmpge:
		return a + " >= " + b
	case rt.OpIfIcmpgt:
		return a + " > " + b
	case rt.OpIfIcmple:
		return a + " <= " + b
	}
	return ""
}

func emitInsn(insn rt.Insn, m *classfile.Method, h int, widths []int, heights map[int32]int, blockOf map[int32]int) (string, int, bool, error) {
	pc := insn.PC
	width := func(index int) int {
		if index < 0 || index >= len(widths) {
			return 1
		}
		return widths[index]
	}
	transfer := transferCode(pc, m, blockOf)

	switch insn.Op {
	case rt.OpNop:
		return "", h, false, nil
	case rt.OpAconstNull:
		return fmt.Sprintf("stack[%d] = rt.RefValue(nil)", h), h + 1, false, nil
	case rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2, rt.OpIconst3, rt.OpIconst4, rt.OpIconst5:
		return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, int32(insn.Op)-3), h + 1, false, nil
	case rt.OpLconst0, rt.OpLconst1:
		return fmt.Sprintf("stack[%d] = rt.LongValue(%d)", h, int64(insn.Op-rt.OpLconst0)), h + 1, false, nil
	case rt.OpFconst0, rt.OpFconst1, rt.OpFconst2:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(%s)", h, strconv.FormatFloat(float64(insn.Op-rt.OpFconst0), 'g', -1, 64)), h + 1, false, nil
	case rt.OpDconst0, rt.OpDconst1:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(%s)", h, strconv.FormatFloat(float64(insn.Op-rt.OpDconst0), 'g', -1, 64)), h + 1, false, nil
	case rt.OpBipush, rt.OpSipush:
		return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, insn.A), h + 1, false, nil
	case rt.OpLdc, rt.OpLdcW, rt.OpLdc2W:
		ref := m.Refs[insn.Ref]
		switch ref.Kind {
		case rt.RefInt:
			return fmt.Sprintf("stack[%d] = rt.IntValue(%d)", h, ref.Int), h + 1, false, nil
		case rt.RefFloat:
			return fmt.Sprintf("stack[%d] = rt.FloatBitsValue(0x%08x)", h, uint32(ref.Bits)), h + 1, false, nil
		case rt.RefLong:
			return fmt.Sprintf("stack[%d] = rt.LongBitsValue(0x%016x)", h, ref.Bits), h + 1, false, nil
		case rt.RefDouble:
			return fmt.Sprintf("stack[%d] = rt.DoubleBitsValue(0x%016x)", h, ref.Bits), h + 1, false, nil
		case rt.RefClass:
			code := fmt.Sprintf("{\nv, thrown := t.ClassConstant(m, %d, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
				ref.CP, strconv.Quote(ref.Class), transfer, h)
			return code, h + 1, false, nil
		case rt.RefString:
			var units []string
			for _, u := range ref.Units {
				units = append(units, strconv.Itoa(int(u)))
			}
			return fmt.Sprintf("stack[%d] = t.StringConstant([]uint16{%s})", h, strings.Join(units, ", ")), h + 1, false, nil
		default:
			return "", h, false, fmt.Errorf("unsupported ldc reference")
		}
	case rt.OpIload, rt.OpLload, rt.OpFload, rt.OpDload, rt.OpAload:
		return fmt.Sprintf("stack[%d] = locals[%d]", h, insn.A), h + 1, false, nil
	case rt.OpIstore, rt.OpLstore, rt.OpFstore, rt.OpDstore, rt.OpAstore:
		return fmt.Sprintf("locals[%d] = stack[%d]", insn.A, h-1), h - 1, false, nil
	case rt.OpIinc:
		return fmt.Sprintf("locals[%d] = rt.IntValue(rt.AsInt(locals[%d]) + %d)", insn.A, insn.A, insn.B), h, false, nil
	case rt.OpIaload, rt.OpLaload, rt.OpFaload, rt.OpDaload, rt.OpAaload, rt.OpBaload, rt.OpCaload, rt.OpSaload:
		var kind byte
		switch insn.Op {
		case rt.OpIaload:
			kind = 'I'
		case rt.OpLaload:
			kind = 'J'
		case rt.OpFaload:
			kind = 'F'
		case rt.OpDaload:
			kind = 'D'
		case rt.OpAaload:
			kind = 'L'
		case rt.OpBaload:
			kind = 'B'
		case rt.OpCaload:
			kind = 'C'
		case rt.OpSaload:
			kind = 'S'
		}
		code := fmt.Sprintf("{\nv, thrown := t.ArrayLoad(stack[%d], rt.AsInt(stack[%d]), %q)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			h-2, h-1, kind, transfer, h-2)
		return code, h - 1, false, nil
	case rt.OpIastore, rt.OpLastore, rt.OpFastore, rt.OpDastore, rt.OpAastore, rt.OpBastore, rt.OpCastore, rt.OpSastore:
		var kind byte
		switch insn.Op {
		case rt.OpIastore:
			kind = 'I'
		case rt.OpLastore:
			kind = 'J'
		case rt.OpFastore:
			kind = 'F'
		case rt.OpDastore:
			kind = 'D'
		case rt.OpAastore:
			kind = 'L'
		case rt.OpBastore:
			kind = 'B'
		case rt.OpCastore:
			kind = 'C'
		case rt.OpSastore:
			kind = 'S'
		}
		code := fmt.Sprintf("{\nthrown := t.ArrayStore(stack[%d], rt.AsInt(stack[%d]), stack[%d], %q)\nif thrown != nil {\n%s\n}\n}",
			h-3, h-2, h-1, kind, transfer)
		return code, h - 3, false, nil
	case rt.OpPop:
		return "", h - 1, false, nil
	case rt.OpPop2:
		if width(len(widths)-1) == 2 {
			return "", h - 1, false, nil
		}
		return "", h - 2, false, nil
	case rt.OpDup:
		return fmt.Sprintf("stack[%d] = stack[%d]", h, h-1), h + 1, false, nil
	case rt.OpDupX1:
		code := fmt.Sprintf("{\nv1 := stack[%d]\nv2 := stack[%d]\nstack[%d] = v1\nstack[%d] = v2\nstack[%d] = v1\n}", h-1, h-2, h, h-1, h-2)
		return code, h + 1, false, nil
	case rt.OpDupX2:
		if width(len(widths)-2) == 2 {
			code := fmt.Sprintf("{\nv1 := stack[%d]\nv2 := stack[%d]\nstack[%d] = v1\nstack[%d] = v2\nstack[%d] = v1\n}", h-1, h-2, h, h-1, h-2)
			return code, h + 1, false, nil
		}
		code := fmt.Sprintf("{\nv1 := stack[%d]\nv3 := stack[%d]\nv2 := stack[%d]\nstack[%d] = v1\nstack[%d] = v2\nstack[%d] = v3\nstack[%d] = v1\n}", h-1, h-3, h-2, h, h-1, h-2, h-3)
		return code, h + 1, false, nil
	case rt.OpDup2:
		if width(len(widths)-1) == 2 {
			return fmt.Sprintf("stack[%d] = stack[%d]", h, h-1), h + 1, false, nil
		}
		code := fmt.Sprintf("{\nv2 := stack[%d]\nv1 := stack[%d]\nstack[%d] = v2\nstack[%d] = v1\n}", h-2, h-1, h, h+1)
		return code, h + 2, false, nil
	case rt.OpDup2X1:
		if width(len(widths)-1) == 2 {
			code := fmt.Sprintf("{\nv1 := stack[%d]\nv2 := stack[%d]\nstack[%d] = v1\nstack[%d] = v2\nstack[%d] = v1\n}", h-1, h-2, h, h-1, h-2)
			return code, h + 1, false, nil
		}
		code := fmt.Sprintf("{\nv3 := stack[%d]\nv2 := stack[%d]\nv1 := stack[%d]\nstack[%d] = v2\nstack[%d] = v1\nstack[%d] = v3\nstack[%d] = v2\nstack[%d] = v1\n}",
			h-3, h-2, h-1, h-3, h-2, h-1, h, h+1)
		return code, h + 2, false, nil
	case rt.OpDup2X2:
		switch {
		case width(len(widths)-1) == 2 && width(len(widths)-2) == 2:
			code := fmt.Sprintf("{\nv1 := stack[%d]\nv2 := stack[%d]\nstack[%d] = v1\nstack[%d] = v2\nstack[%d] = v1\n}", h-1, h-2, h, h-1, h-2)
			return code, h + 1, false, nil
		case width(len(widths)-1) == 2:
			code := fmt.Sprintf("{\nv3 := stack[%d]\nv2 := stack[%d]\nv1 := stack[%d]\nstack[%d] = v1\nstack[%d] = v3\nstack[%d] = v2\nstack[%d] = v1\n}", h-3, h-2, h-1, h-3, h-2, h-1, h)
			return code, h + 1, false, nil
		case width(len(widths)-3) == 2:
			code := fmt.Sprintf("{\nv3 := stack[%d]\nv2 := stack[%d]\nv1 := stack[%d]\nstack[%d] = v2\nstack[%d] = v1\nstack[%d] = v3\nstack[%d] = v2\nstack[%d] = v1\n}", h-3, h-2, h-1, h-3, h-2, h-1, h, h+1)
			return code, h + 2, false, nil
		default:
			code := fmt.Sprintf("{\nv4 := stack[%d]\nv3 := stack[%d]\nv2 := stack[%d]\nv1 := stack[%d]\nstack[%d] = v2\nstack[%d] = v1\nstack[%d] = v4\nstack[%d] = v3\nstack[%d] = v2\nstack[%d] = v1\n}",
				h-4, h-3, h-2, h-1, h-4, h-3, h-2, h-1, h, h+1)
			return code, h + 2, false, nil
		}
	case rt.OpSwap:
		return fmt.Sprintf("stack[%d], stack[%d] = stack[%d], stack[%d]", h-1, h-2, h-2, h-1), h, false, nil
	case rt.OpIadd:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) + rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLadd:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) + rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFadd:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(rt.AsFloat(stack[%d]) + rt.AsFloat(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDadd:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(rt.AsDouble(stack[%d]) + rt.AsDouble(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIsub:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) - rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLsub:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) - rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFsub:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(rt.AsFloat(stack[%d]) - rt.AsFloat(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDsub:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(rt.AsDouble(stack[%d]) - rt.AsDouble(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpImul:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) * rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLmul:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) * rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFmul:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(rt.AsFloat(stack[%d]) * rt.AsFloat(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDmul:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(rt.AsDouble(stack[%d]) * rt.AsDouble(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIdiv:
		code := fmt.Sprintf("{\nd := rt.AsInt(stack[%d])\nif d == 0 {\nthrown := t.Throwable(\"java/lang/ArithmeticException\")\n%s\n}\nstack[%d] = rt.IntValue(rt.AsInt(stack[%d]) / d)\n}",
			h-1, transfer, h-2, h-2)
		return code, h - 1, false, nil
	case rt.OpLdiv:
		code := fmt.Sprintf("{\nd := rt.AsLong(stack[%d])\nif d == 0 {\nthrown := t.Throwable(\"java/lang/ArithmeticException\")\n%s\n}\nstack[%d] = rt.LongValue(rt.AsLong(stack[%d]) / d)\n}",
			h-1, transfer, h-2, h-2)
		return code, h - 1, false, nil
	case rt.OpIrem:
		return "", h, false, fmt.Errorf("irem has no AOT lowering (pc %d, bci %d)", pc, pc)
	case rt.OpLrem:
		code := fmt.Sprintf("{\nd := rt.AsLong(stack[%d])\nif d == 0 {\nthrown := t.Throwable(\"java/lang/ArithmeticException\")\n%s\n}\nstack[%d] = rt.LongValue(rt.AsLong(stack[%d]) %% d)\n}",
			h-1, transfer, h-2, h-2)
		return code, h - 1, false, nil
	case rt.OpFdiv:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(rt.AsFloat(stack[%d]) / rt.AsFloat(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDdiv:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(rt.AsDouble(stack[%d]) / rt.AsDouble(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFrem:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(rt.FRem(rt.AsFloat(stack[%d]), rt.AsFloat(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDrem:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(rt.DRem(rt.AsDouble(stack[%d]), rt.AsDouble(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIneg:
		return fmt.Sprintf("stack[%d] = rt.IntValue(-rt.AsInt(stack[%d]))", h-1, h-1), h, false, nil
	case rt.OpLneg:
		return fmt.Sprintf("stack[%d] = rt.LongValue(-rt.AsLong(stack[%d]))", h-1, h-1), h, false, nil
	case rt.OpFneg:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(-rt.AsFloat(stack[%d]))", h-1, h-1), h, false, nil
	case rt.OpDneg:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(-rt.AsDouble(stack[%d]))", h-1, h-1), h, false, nil
	case rt.OpIshl:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) << rt.IntShiftDistance(rt.AsInt(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIshr:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) >> rt.IntShiftDistance(rt.AsInt(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIushr:
		return fmt.Sprintf("stack[%d] = rt.IntValue(int32(uint32(rt.AsInt(stack[%d])) >> rt.IntShiftDistance(rt.AsInt(stack[%d]))))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLshl:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) << rt.LongShiftDistance(rt.AsInt(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLshr:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) >> rt.LongShiftDistance(rt.AsInt(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLushr:
		return fmt.Sprintf("stack[%d] = rt.LongValue(int64(uint64(rt.AsLong(stack[%d])) >> rt.LongShiftDistance(rt.AsInt(stack[%d]))))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIand:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) & rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIor:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) | rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIxor:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.AsInt(stack[%d]) ^ rt.AsInt(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLand:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) & rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLor:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) | rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpLxor:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.AsLong(stack[%d]) ^ rt.AsLong(stack[%d]))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpI2l:
		return fmt.Sprintf("stack[%d] = rt.LongValue(int64(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpI2f:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(float32(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpI2d:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(float64(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpL2i:
		return fmt.Sprintf("stack[%d] = rt.IntValue(int32(rt.AsLong(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpL2f:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(float32(rt.AsLong(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpL2d:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(float64(rt.AsLong(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpF2i:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.F2I(rt.AsFloat(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpF2l:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.F2L(rt.AsFloat(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpF2d:
		return fmt.Sprintf("stack[%d] = rt.DoubleValue(float64(rt.AsFloat(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpD2i:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.D2I(rt.AsDouble(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpD2l:
		return fmt.Sprintf("stack[%d] = rt.LongValue(rt.D2L(rt.AsDouble(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpD2f:
		return fmt.Sprintf("stack[%d] = rt.FloatValue(float32(rt.AsDouble(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpI2b:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.I2B(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpI2c:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.I2C(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpI2s:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.I2S(rt.AsInt(stack[%d])))", h-1, h-1), h, false, nil
	case rt.OpLcmp:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[%d]), rt.AsLong(stack[%d])))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFcmpl:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[%d]), rt.AsFloat(stack[%d]), -1))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpFcmpg:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[%d]), rt.AsFloat(stack[%d]), 1))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDcmpl:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[%d]), rt.AsDouble(stack[%d]), -1))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpDcmpg:
		return fmt.Sprintf("stack[%d] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[%d]), rt.AsDouble(stack[%d]), 1))", h-2, h-2, h-1), h - 1, false, nil
	case rt.OpIfeq, rt.OpIfne, rt.OpIflt, rt.OpIfge, rt.OpIfgt, rt.OpIfle,
		rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfIcmplt, rt.OpIfIcmpge, rt.OpIfIcmpgt, rt.OpIfIcmple:
		var cond string
		if insn.Op >= rt.OpIfIcmpeq {
			cond = booleanCode(insn.Op, fmt.Sprintf("rt.AsInt(stack[%d])", h-2), fmt.Sprintf("rt.AsInt(stack[%d])", h-1))
		} else {
			cond = booleanCode(insn.Op, fmt.Sprintf("rt.AsInt(stack[%d])", h-1), "")
		}
		target, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		fall, err := blockID(pc+insn.Len, blockOf)
		if err != nil {
			return "", h, false, err
		}
		pops := 1
		if insn.Op >= rt.OpIfIcmpeq {
			pops = 2
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fall), h - pops, true, nil
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
		fall, err := blockID(pc+insn.Len, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fall), h - 2, true, nil
	case rt.OpIfnull, rt.OpIfnonnull:
		op := "=="
		if insn.Op == rt.OpIfnonnull {
			op = "!="
		}
		cond := fmt.Sprintf("stack[%d].Ref %s nil", h-1, op)
		target, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		fall, err := blockID(pc+insn.Len, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("if %s {\nblock = %d\ncontinue\n}\nblock = %d\ncontinue", cond, target, fall), h - 1, true, nil
	case rt.OpGoto, rt.OpGotoW:
		id, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		return fmt.Sprintf("block = %d\ncontinue", id), h, true, nil
	case rt.OpTableswitch, rt.OpLookupswitch:
		def, err := blockID(insn.A, blockOf)
		if err != nil {
			return "", h, false, err
		}
		var body strings.Builder
		fmt.Fprintf(&body, "switch rt.AsInt(stack[%d]) {\n", h-1)
		for i, key := range insn.Keys {
			id, err := blockID(insn.Targets[i], blockOf)
			if err != nil {
				return "", h, false, err
			}
			fmt.Fprintf(&body, "case %d:\nblock = %d\ncontinue\n", key, id)
		}
		fmt.Fprintf(&body, "default:\nblock = %d\ncontinue\n}", def)
		return body.String(), h - 1, true, nil
	case rt.OpReturn:
		return "return rt.Normal()", h, true, nil
	case rt.OpIreturn, rt.OpLreturn, rt.OpFreturn, rt.OpDreturn, rt.OpAreturn:
		return fmt.Sprintf("return rt.OkValue(stack[%d])", h-1), h - 1, true, nil
	case rt.OpGetstatic:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nv, thrown := t.GetStatic(m, %d, %s, %s, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			ref.CP, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), transfer, h)
		return code, h + 1, false, nil
	case rt.OpPutstatic:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nthrown := t.PutStatic(m, %d, %s, %s, %s, stack[%d])\nif thrown != nil {\n%s\n}\n}",
			ref.CP, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), h-1, transfer)
		return code, h - 1, false, nil
	case rt.OpGetfield:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nv, thrown := t.GetField(m, %d, stack[%d], %s, %s, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			ref.CP, h-1, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), transfer, h-1)
		return code, h, false, nil
	case rt.OpPutfield:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nthrown := t.PutField(m, %d, stack[%d], %s, %s, %s, stack[%d])\nif thrown != nil {\n%s\n}\n}",
			ref.CP, h-2, strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), h-1, transfer)
		return code, h - 2, false, nil
	case rt.OpNew:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nobj, thrown := t.NewObject(m, %d, %s)\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(obj)\n}",
			ref.CP, strconv.Quote(ref.Class), transfer, h)
		return code, h + 1, false, nil
	case rt.OpNewarray:
		desc, _ := rt.NewarrayDescriptor(insn.A)
		_ = desc
		code := fmt.Sprintf("{\nobj, thrown := t.NewPrimitiveArray(m, %d, rt.AsInt(stack[%d]))\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(obj)\n}",
			insn.A, h-1, transfer, h-1)
		return code, h, false, nil
	case rt.OpAnewarray:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nobj, thrown := t.NewReferenceArray(m, %d, %s, rt.AsInt(stack[%d]))\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(obj)\n}",
			ref.CP, strconv.Quote(ref.Class), h-1, transfer, h-1)
		return code, h, false, nil
	case rt.OpMultianewarray:
		ref := m.Refs[insn.Ref]
		dims := int(insn.B)
		var items []string
		for i := 0; i < dims; i++ {
			items = append(items, fmt.Sprintf("rt.AsInt(stack[%d])", h-dims+i))
		}
		code := fmt.Sprintf("{\ncounts := []int32{%s}\nobj, thrown := t.MultiANewArray(m, %d, %s, %d, counts)\nif thrown != nil {\n%s\n}\nstack[%d] = rt.RefValue(obj)\n}",
			strings.Join(items, ", "), ref.CP, strconv.Quote(ref.Class), dims, transfer, h-dims)
		return code, h - dims + 1, false, nil
	case rt.OpArraylength:
		code := fmt.Sprintf("{\nn, thrown := t.ArrayLength(stack[%d])\nif thrown != nil {\n%s\n}\nstack[%d] = rt.IntValue(n)\n}",
			h-1, transfer, h-1)
		return code, h, false, nil
	case rt.OpAthrow:
		code := fmt.Sprintf("{\nthrown := stack[%d].Ref\nif thrown == nil {\nthrown = t.Throwable(\"java/lang/NullPointerException\")\n}\n%s\n}",
			h-1, transfer)
		return code, h - 1, true, nil
	case rt.OpCheckcast:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nv, thrown := t.CheckCast(m, %d, %s, stack[%d])\nif thrown != nil {\n%s\n}\nstack[%d] = v\n}",
			ref.CP, strconv.Quote(ref.Class), h-1, transfer, h-1)
		return code, h, false, nil
	case rt.OpInstanceof:
		ref := m.Refs[insn.Ref]
		code := fmt.Sprintf("{\nn, thrown := t.InstanceOf(m, %d, %s, stack[%d])\nif thrown != nil {\n%s\n}\nstack[%d] = rt.IntValue(n)\n}",
			ref.CP, strconv.Quote(ref.Class), h-1, transfer, h-1)
		return code, h, false, nil
	case rt.OpInvokestatic, rt.OpInvokevirtual, rt.OpInvokespecial, rt.OpInvokeinterface:
		ref := m.Refs[insn.Ref]
		mt, err := rt.ParseMethodType(ref.Descriptor)
		if err != nil {
			return "", h, false, err
		}
		n := mt.ArgCount
		args := fmt.Sprintf("stack[%d:%d]", h-n, h)
		var call string
		newH := h - n
		switch insn.Op {
		case rt.OpInvokestatic:
			call = fmt.Sprintf("t.InvokeStatic(m, %d, %d, %s, %s, %s, %s)", pc, ref.CP,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
		case rt.OpInvokespecial:
			call = fmt.Sprintf("t.InvokeSpecial(m, %d, %d, stack[%d], %s, %s, %s, %s)", pc, ref.CP, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH--
		case rt.OpInvokevirtual:
			call = fmt.Sprintf("t.InvokeVirtual(m, %d, %d, stack[%d], %s, %s, %s, %s)", pc, ref.CP, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH--
		case rt.OpInvokeinterface:
			call = fmt.Sprintf("t.InvokeInterface(m, %d, %d, stack[%d], %s, %s, %s, %s)", pc, ref.CP, h-n-1,
				strconv.Quote(ref.Owner), strconv.Quote(ref.Name), strconv.Quote(ref.Descriptor), args)
			newH--
		}
		var body strings.Builder
		body.WriteString("{\nout := " + call + "\nif out.Thrown != nil {\nthrown := out.Thrown\n")
		body.WriteString(transfer)
		body.WriteString("\n}\n")
		if mt.Return.Kind != 'V' {
			fmt.Fprintf(&body, "stack[%d] = out.Value\n", newH)
			newH++
		}
		body.WriteString("}")
		return body.String(), newH, false, nil
	default:
		return "", h, false, fmt.Errorf("unsupported opcode 0x%02x", insn.Op)
	}
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
