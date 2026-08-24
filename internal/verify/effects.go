package verify

import (
	"fmt"
	"os"

	"catty/internal/classfile"
)

// Table-driven instruction typing. Each opcode maps to a simFunc that pops
// expectations, pushes results and reports successors. Irregular opcodes
// (ldc family, fields, invokes, switches, wide) have dedicated functions;
// everything else is generated from families so each semantics lives once.

type simFunc func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error)

type cursor struct {
	code []byte
	pc   int // position AFTER the opcode byte
}

func (cu *cursor) u1() uint8 {
	v := cu.code[cu.pc]
	cu.pc++
	return v
}
func (cu *cursor) u2() uint16 {
	v := uint16(cu.code[cu.pc])<<8 | uint16(cu.code[cu.pc+1])
	cu.pc += 2
	return v
}
func (cu *cursor) s2() int32 { return int32(int16(cu.u2())) }
func (cu *cursor) align4() {
	for cu.pc%4 != 0 {
		cu.pc++
	}
}

func (cu *cursor) s4() int32 {
	v := uint32(cu.code[cu.pc])<<24 | uint32(cu.code[cu.pc+1])<<16 |
		uint32(cu.code[cu.pc+2])<<8 | uint32(cu.code[cu.pc+3])
	cu.pc += 4
	return int32(v)
}

var effects map[byte]simFunc

func init() { buildEffects() }

func fall(st *frame, nextPC int) []successor {
	return []successor{{pc: nextPC, st: st.clone()}}
}

// --- building blocks ---------------------------------------------------------

func constK(v vtype) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		st.push(v)
		return fall(st, pc+1), nil
	}
}

func nopSim(pc int, st *frame) ([]successor, error) { return fall(st, pc+1), nil }

func popK(k vkind, what string) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if _, err := st.popExpect(k, what); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		return fall(st, pc+1), nil
	}
}

func fixedIdx(n int) func(*cursor) int { return func(*cursor) int { return n } }
func byteIdx(cu *cursor) int           { return int(cu.u1()) }
func shortIdx(cu *cursor) int          { return int(cu.s2()) }

func loadLocal(idxOf func(*cursor) int, want vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := cursor{code: m.Code.Code, pc: pc + 1}
		idx := idxOf(&cu)
		if idx >= len(st.locals) {
			return nil, verr(c.mn, pc, "local[%d] out of range (%d locals)", idx, len(st.locals))
		}
		got := st.locals[idx]
		switch want {
		case kInt, kFloat, kLong, kDouble:
			if got.kind != want && got.kind != kUnknown {
				return nil, verr(c.mn, pc, "load %s from local holding %s",
					vtype{kind: want}, got)
			}
			got = vtype{kind: want}
		case kObj:
			if !isRef(got) {
				return nil, verr(c.mn, pc, "aload of %s", got)
			}
		}
		st.push(got)
		return fall(st, cu.pc), nil
	}
}

func storeLocal(idxOf func(*cursor) int, want vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := cursor{code: m.Code.Code, pc: pc + 1}
		idx := idxOf(&cu)
		var v vtype
		var err error
		if want == kLong || want == kDouble {
			v, err = st.popCat2()
		} else {
			v, err = st.popExpect(want, "store value")
			if err == nil && want == kObj && !isRef(v) {
				err = fmt.Errorf("astore of non-reference %s", v)
			}
		}
		if err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		setLocalPair(st, idx, v)
		return fall(st, cu.pc), nil
	}
}

func binary(kind vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if kind == kLong || kind == kDouble {
			// Category-2 operands sit as [value, top] pairs; plain
			// popExpect would trip on the sentinel (found via Bench
			// lsub — Route-C gap surfaced early).
			if _, err := st.popCat2(); err != nil {
				return nil, verr(c.mn, pc, "rhs: %v", err)
			}
			if _, err := st.popCat2(); err != nil {
				return nil, verr(c.mn, pc, "lhs: %v", err)
			}
			st.push(vtype{kind: kind})
			return fall(st, pc+1), nil
		}
		if _, err := st.popExpect(kind, "rhs"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		if _, err := st.popExpect(kind, "lhs"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		st.push(vtype{kind: kind})
		return fall(st, pc+1), nil
	}
}

// divRem adds ArithmeticException handler edges on top of binary typing.
func divRem(kind vkind) simFunc {
	base := binary(kind)
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		succs, err := base(c, m, pc, st)
		if err != nil {
			return nil, err
		}
		return append(succs, c.throwSuccs(m, pc, st,
			tObj("java/lang/ArithmeticException"))...), nil
	}
}

// shift: value kind result; count is always int.
func shift(value vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if _, err := st.popExpect(kInt, "shift count"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		v, err := st.popExpect(value, "shift value")
		if err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		st.push(v)
		return fall(st, pc+1), nil
	}
}

func neg(kind vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if _, err := st.popExpect(kind, "operand"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		st.push(vtype{kind: kind})
		return fall(st, pc+1), nil
	}
}

func conv(from, to vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if isCat2Kind(from) {
			if _, err := st.popCat2(); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
		} else if _, err := st.popExpect(from, "conversion input"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		st.push(vtype{kind: to})
		return fall(st, pc+1), nil
	}
}

func cmp2(a, b vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if _, err := st.popExpect(b, "cmp rhs"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		if _, err := st.popExpect(a, "cmp lhs"); err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		st.push(tInt)
		return fall(st, pc+1), nil
	}
}

// branch: both edges feasible at type level; validates popped operands.
func branch(pops ...vkind) simFunc {
	return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		for i := len(pops) - 1; i >= 0; i-- {
			if _, err := st.popExpect(pops[i], "branch operand"); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
		}
		cu := cursor{code: m.Code.Code, pc: pc + 1}
		tgt := pc + int(cu.s2())
		out := []successor{{pc: tgt, st: st.clone()}, {pc: pc + 3, st: st.clone()}}
		return out, nil
	}
}

func gotoSim(m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
	cu := cursor{code: m.Code.Code, pc: pc + 1}
	return []successor{{pc: pc + int(cu.s2()), st: st.clone()}}, nil
}

// --- registration ------------------------------------------------------------

func buildEffects() {
	effects = make(map[byte]simFunc, 224)

	nop := func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		return fall(st, pc+1), nil
	}
	set := func(o byte, f simFunc) { effects[o] = f }

	set(0x00, nop)
	set(0x01, constK(tNull))
	for o := byte(0x02); o <= 0x08; o++ {
		set(o, constK(tInt))
	}
	set(0x09, constK(tLong))
	set(0x0a, constK(tLong))
	for o := byte(0x0b); o <= 0x0d; o++ {
		set(o, constK(tFloat))
	}
	set(0x0e, constK(tDouble))
	set(0x0f, constK(tDouble))

	set(0x10, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := &cursor{code: m.Code.Code, pc: pc + 1}
		cu.u1()
		st.push(tInt)
		return fall(st, cu.pc), nil
	})
	set(0x11, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := &cursor{code: m.Code.Code, pc: pc + 1}
		cu.u2()
		st.push(tInt)
		return fall(st, cu.pc), nil
	})

	for i := 0; i < 4; i++ {
		set(byte(0x1a+i), loadLocal(fixedIdx(i), kInt))
		set(byte(0x22+i), loadLocal(fixedIdx(i), kFloat))
		set(byte(0x2a+i), loadLocal(fixedIdx(i), kObj))
		set(byte(0x1e+i), loadLocal(fixedIdx(i), kLong))
		set(byte(0x26+i), loadLocal(fixedIdx(i), kDouble))
		set(byte(0x3b+i), storeLocal(fixedIdx(i), kInt))
		set(byte(0x43+i), storeLocal(fixedIdx(i), kFloat))
		set(byte(0x4b+i), storeLocal(fixedIdx(i), kObj))
		set(byte(0x3f+i), storeLocal(fixedIdx(i), kLong))
		set(byte(0x47+i), storeLocal(fixedIdx(i), kDouble))
	}
	set(0x15, loadLocal(byteIdx, kInt))
	set(0x17, loadLocal(byteIdx, kFloat))
	set(0x19, loadLocal(byteIdx, kObj))
	set(0x16, loadLocal(byteIdx, kLong))
	set(0x18, loadLocal(byteIdx, kDouble))
	set(0x36, storeLocal(byteIdx, kInt))
	set(0x38, storeLocal(byteIdx, kFloat))
	set(0x3a, storeLocal(byteIdx, kObj))
	set(0x37, storeLocal(byteIdx, kLong))
	set(0x39, storeLocal(byteIdx, kDouble))

	arrLoad := func(elem vkind) simFunc {
		return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
			if _, err := st.popExpect(kInt, "array index"); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			arr, err := st.popExpect(kObj, "arrayref")
			if err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			succs := c.throwSuccs(m, pc, st, tObj("java/lang/NullPointerException"))
			if arr.kind == kObj && len(arr.class) > 0 && arr.class[0] == '[' {
				el := arrayElem(arr.class)
				if el.kind != elem && el.kind != kUnknown {
					return nil, verr(c.mn, pc, "%s load on %s",
						vtype{kind: elem}, arr.class)
				}
			}
			st.push(vtype{kind: elem})
			next := fall(st, pc+1)
			return append(succs, next...), nil
		}
	}
	set(0x2e, arrLoad(kInt))
	set(0x33, arrLoad(kInt)) // baload ([B/[Z)
	set(0x34, arrLoad(kInt)) // caload
	set(0x35, arrLoad(kInt)) // saload
	set(0x2f, arrLoad(kLong))
	set(0x30, arrLoad(kFloat))
	set(0x31, arrLoad(kDouble))
	set(0x32, arrLoad(kObj))

	arrStore := func(elem vkind) simFunc {
		return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
			var err error
			if elem == kLong || elem == kDouble {
				_, err = st.popCat2()
			} else {
				_, err = st.popExpect(elem, "array store value")
			}
			if err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			if _, err := st.popExpect(kInt, "array index"); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			arr, aerr := st.popExpect(kObj, "arrayref")
			if aerr != nil {
				return nil, verr(c.mn, pc, "%v", aerr)
			}
			succs := c.throwSuccs(m, pc, st, tObj("java/lang/NullPointerException"))
			if arr.kind == kObj && len(arr.class) > 0 && arr.class[0] == '[' {
				el := arrayElem(arr.class)
				if el.kind != elem && el.kind != kUnknown {
					return nil, verr(c.mn, pc, "store %s into %s",
						vtype{kind: elem}, arr.class)
				}
			}
			next := fall(st, pc+1)
			return append(succs, next...), nil
		}
	}
	set(0x4f, arrStore(kInt))
	set(0x54, arrStore(kInt)) // bastore
	set(0x55, arrStore(kInt)) // castore
	set(0x56, arrStore(kInt)) // sastore
	set(0x50, arrStore(kLong))
	set(0x51, arrStore(kFloat))
	set(0x52, arrStore(kDouble))
	set(0x53, arrStore(kObj))

	set(0x57, popK(kUnknown, "pop"))
	set(0x58, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if len(st.stack) < 2 {
			return nil, verr(c.mn, pc, "pop2 underflow")
		}
		top := st.stack[len(st.stack)-1]
		if top.kind == kTop {
			if _, err := st.popCat2(); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
		} else {
			if _, err := st.pop(); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			if _, err := st.pop(); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
		}
		return fall(st, pc+1), nil
	})
	dupRaw := func(n int) simFunc {
		return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
			if len(st.stack) < n {
				return nil, verr(c.mn, pc, "dup underflow")
			}
			tail := append([]vtype(nil), st.stack[len(st.stack)-n:]...)
			st.stack = append(st.stack, tail...)
			return fall(st, pc+1), nil
		}
	}
	dupBelow := func(n, below int) simFunc {
		return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
			total := n + below
			if len(st.stack) < total {
				if os.Getenv("CATTY_VDBG") != "" {
					fmt.Fprintf(os.Stderr, "[dupx] %s pc=%d need=%d have=%d stack=%v\n",
						c.mn, pc, total, len(st.stack), st.stack)
				}
				return nil, verr(c.mn, pc, "dup underflow")
			}
			n2 := n + below
		head := append([]vtype(nil), st.stack[len(st.stack)-n:len(st.stack)]...) // duplicated block
			mid := append([]vtype(nil), st.stack[len(st.stack)-n2:len(st.stack)-n]...) // skipped slots
			prefix := append([]vtype(nil), st.stack[:len(st.stack)-n2]...)
			out := make([]vtype, 0, len(prefix)+n2+n)
			out = append(out, prefix...)
			out = append(out, head...)
			out = append(out, mid...)
			out = append(out, head...)
			st.stack = out
			return fall(st, pc+1), nil
		}
	}
	set(0x59, dupRaw(1))       //      -> v1
	set(0x5a, dupBelow(1, 1)) // v2,v1 -> v1,v2,v1
	set(0x5b, dupBelow(1, 2)) // v3,v2,v1 -> v1,v3,v2,v1
	set(0x5c, dupRaw(2))      // v2,v1 -> v2,v1,v2,v1
	set(0x5d, dupBelow(2, 1)) // v3,v2,v1 -> v2,v1,v3,v2,v1
	set(0x5e, dupBelow(2, 2))
	set(0x5f, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		if len(st.stack) < 2 {
			return nil, verr(c.mn, pc, "swap underflow")
		}
		a := st.stack[len(st.stack)-1]
		b := st.stack[len(st.stack)-2]
		if isCat2(a) || isCat2(b) {
			return nil, verr(c.mn, pc, "swap on category-2")
		}
		st.stack[len(st.stack)-1], st.stack[len(st.stack)-2] = b, a
		return fall(st, pc+1), nil
	})

	for _, o := range []byte{0x60, 0x64, 0x68, 0x7e, 0x80, 0x82} {
		set(o, binary(kInt))
	}
	for _, o := range []byte{0x61, 0x65, 0x69, 0x7f, 0x81, 0x83} {
		set(o, binary(kLong))
	}
	for _, o := range []byte{0x62, 0x66, 0x6a} {
		set(o, binary(kFloat))
	}
	for _, o := range []byte{0x63, 0x67, 0x6b} {
		set(o, binary(kDouble))
	}
	set(0x6c, divRem(kInt))
	set(0x70, divRem(kInt))
	set(0x6d, divRem(kLong))
	set(0x71, divRem(kLong))
	set(0x6e, binary(kFloat))
	set(0x72, binary(kFloat))
	set(0x6f, binary(kDouble))
	set(0x73, binary(kDouble))
	set(0x74, neg(kInt))
	set(0x75, neg(kLong))
	set(0x76, neg(kFloat))
	set(0x77, neg(kDouble))
	set(0x78, shift(kInt))
	set(0x79, shift(kLong))
	set(0x7a, shift(kInt))
	set(0x7b, shift(kLong))
	set(0x7c, shift(kInt))
	set(0x7d, shift(kLong))

	convPairs := []struct {
		op       byte
		from, to vkind
	}{
		{0x85, kInt, kLong}, {0x86, kInt, kFloat}, {0x87, kInt, kDouble},
		{0x88, kLong, kInt}, {0x89, kLong, kFloat}, {0x8a, kLong, kDouble},
		{0x8b, kFloat, kInt}, {0x8c, kFloat, kLong}, {0x8d, kFloat, kDouble},
		{0x8e, kDouble, kInt}, {0x8f, kDouble, kLong}, {0x90, kDouble, kFloat},
		{0x91, kInt, kInt}, {0x92, kInt, kInt}, {0x93, kInt, kInt},
	}
	for _, p := range convPairs {
		set(p.op, conv(p.from, p.to))
	}

	set(0x84, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := &cursor{code: m.Code.Code, pc: pc + 1}
		idx := int(cu.u1())
		cu.s2() // signed constant (typing unaffected)
		if idx >= len(st.locals) {
			return nil, verr(c.mn, pc, "iinc local[%d] out of range", idx)
		}
		got := st.locals[idx]
		if got.kind != kInt && got.kind != kUnknown {
			return nil, verr(c.mn, pc, "iinc on %s local", got)
		}
		return fall(st, pc+3), nil
	})

	set(0x94, cmp2(kLong, kLong))
	set(0x95, cmp2(kFloat, kFloat))
	set(0x96, cmp2(kFloat, kFloat))
	set(0x97, cmp2(kDouble, kDouble))
	set(0x98, cmp2(kDouble, kDouble))

	intBranches := []byte{0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e}
	for _, o := range intBranches {
		set(o, branch(kInt))
	}
	pairBranches := []byte{0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4}
	for _, o := range pairBranches {
		set(o, branch(kInt, kInt))
	}
	set(0xa5, branch(kObj, kObj))
	set(0xa6, branch(kObj, kObj))
	set(0xc6, branch(kObj))
	set(0xc7, branch(kObj))
	retK := func(k vkind, cat2 bool) simFunc {
		return func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
			want, hasRet, err := retVType(m.Desc)
			if err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			if !hasRet {
				return nil, verr(c.mn, pc, "return with value in void method")
			}
			if want.kind != k && k != kUnknown &&
				!(k == kObj && isRef(want)) && !(want.kind == kObj && isRef(vtype{kind: k})) {
				return nil, verr(c.mn, pc, "return kind mismatch: %s vs %s",
					vtype{kind: k}, want)
			}
			if cat2 {
				if _, err := st.popCat2(); err != nil {
					return nil, verr(c.mn, pc, "%v", err)
				}
				return nil, nil
			}
			if _, err := st.pop(); err != nil {
				return nil, verr(c.mn, pc, "%v", err)
			}
			return nil, nil
		}
	}
	set(0xac, retK(kInt, false))
	set(0xad, retK(kLong, true))
	set(0xae, retK(kFloat, false))
	set(0xaf, retK(kDouble, true))
	set(0xb0, retK(kObj, false))
	set(0xb1, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		_, hasRet, err := retVType(m.Desc)
		if err != nil {
			return nil, verr(c.mn, pc, "%v", err)
		}
		if hasRet {
			return nil, verr(c.mn, pc, "void return in value-returning method")
		}
		return nil, nil
	})

	gotoFn := func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := &cursor{code: m.Code.Code, pc: pc + 1}
		return []successor{{pc: pc + int(cu.s2()), st: st.clone()}}, nil
	}
	set(0xa7, gotoFn)
	set(0xc8, func(c *checker, m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
		cu := &cursor{code: m.Code.Code, pc: pc + 1}
		off := cu.s4()
		return []successor{{pc: pc + int(off), st: st.clone()}}, nil
	})
}
