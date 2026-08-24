package emitter

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/classfile"
)

// emitOne translates one instruction at pc.
func (e *methodEmitter) emitOne(pc int) error {
	code := e.code
	op := code[pc]

	fail := func(format string, args ...any) error {
		return fmt.Errorf("pc=%d: %s", pc, fmt.Sprintf(format, args...))
	}
	popRef := func() string { e.depth--; return fmt.Sprintf("s%d", e.depth) }
	popRefCat2 := func() string { e.depth -= 2; return fmt.Sprintf("s%d", e.depth) }
	pushV := func(expr string) { e.p("s%d = %s", e.depth, expr); e.depth++ }
	excAfter := func() { e.excDispatch(pc) }
	pushCat2 := func(expr string) {
		e.p("s%d = %s", e.depth, expr)
		e.depth += 2
	}

	switch op {

	case 0x00: // nop

	case 0x0e, 0x0f: // dconst_0, dconst_1
		pushCat2(fmt.Sprintf("float64(%d)", int(op)-0x0e))
	case 0x09, 0x0a: // lconst_0, lconst_1 (cat2)
		pushCat2(fmt.Sprintf("int64(%d)", int(op)-0x09))
	case 0x0b, 0x0c, 0x0d: // fconst_0..2
		pushV(fmt.Sprintf("float32(%d)", int(op)-0x0b))
	case 0x01:
		pushV("nil")
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
		pushV(fmt.Sprintf("int32(%d)", int32(op)-3))
	case 0x10: // bipush
		pushV(fmt.Sprintf("int32(%d)", int32(int8(code[pc+1]))))
	case 0x11: // sipush
		v := int32(int16(uint16(code[pc+1])<<8 | uint16(code[pc+2])))
		pushV(fmt.Sprintf("int32(%d)", v))

	case 0x12, 0x13: // ldc / ldc_w
		var idx uint16
		if op == 0x12 {
			idx = uint16(code[pc+1])
		} else {
			idx = uint16(code[pc+1])<<8 | uint16(code[pc+2])
		}
		cnst := &e.cf.Constants[idx]
		switch cnst.Tag {
		case classfile.CInteger:
			pushV(fmt.Sprintf("int32(%d)", cnst.IntVal))
		case classfile.CFloat:
			pushV(fmt.Sprintf("float32(%v)", cnst.FloatVal))
		case classfile.CString:
			sv, serr := e.cf.UTF8(cnst.Idx1)
			if serr != nil {
				return fail("%v", serr)
			}
			pushV(fmt.Sprintf("genrt.Str(%q)", sv))
		default:
			return fail("ldc on %s unsupported in v1", cnst.Tag)
		}

	case 0x14: // ldc2_w (3 bytes; category-2 result)
		idx := uint16(code[pc+1])<<8 | uint16(code[pc+2])
		cnst := &e.cf.Constants[idx]
		switch cnst.Tag {
		case classfile.CLong:
			pushCat2(fmt.Sprintf("int64(%d)", cnst.LongVal))
		case classfile.CDouble:
			pushCat2(fmt.Sprintf("%v", cnst.DoubleVal))
		default:
			return fail("ldc2_w on %s unsupported", cnst.Tag)
		}

	case 0x15, 0x19: // iload / aload
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		pushV(e.LocalName(idx))
	case 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d:
		n := int(op - 0x1a)
		if op >= 0x2a {
			n = int(op - 0x2a)
		}
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		pushV(e.LocalName(n))

	case 0x17: // fload (indexed)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		pushV(e.LocalName(idx))
	case 0x18: // dload (indexed, cat2)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		pushCat2(e.LocalName(idx))
	case 0x22, 0x23, 0x24, 0x25: // fload_0..3
		n := int(op - 0x22)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		pushV(e.LocalName(n))
	case 0x26, 0x27, 0x28, 0x29: // dload_0..3 (cat2)
		n := int(op - 0x26)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		pushCat2(e.LocalName(n))
	case 0x16: // lload
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		pushCat2(e.LocalName(idx))
	case 0x1e, 0x1f, 0x20, 0x21: // lload_0..3
		n := int(op - 0x1e)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		pushCat2(e.LocalName(n))
	case 0x39: // dstore (indexed, cat2)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(idx), popRefCat2())
	case 0x47, 0x48, 0x49, 0x4a: // dstore_0..3 (cat2)
		n := int(op - 0x47)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRefCat2())
	case 0x36, 0x3a: // istore / astore
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(idx), popRef())
	case 0x3b, 0x3c, 0x3d, 0x3e:
		n := int(op - 0x3b)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRef())
	case 0x3f, 0x40, 0x41, 0x42: // lstore_0..3 (cat2)
		n := int(op - 0x3f)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRefCat2())
	case 0x43, 0x44, 0x45, 0x46: // fstore_0..3
		n := int(op - 0x43)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRef())
	case 0x4b, 0x4c, 0x4d, 0x4e: // astore_0..3
		n := int(op - 0x4b)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRef())

	case 0x38: // fstore (indexed)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(idx), popRef())
	case 0x62, 0x66, 0x6a, 0x6e: // dadd,dsub,dmul,ddiv (cat2)
		b, a := popRefCat2(), popRefCat2()
		opStr := map[byte]string{0x62: "+", 0x66: "-", 0x6a: "*", 0x6e: "/"}[op]
		pushCat2(fmt.Sprintf("(%s.(float64)) %s (%s.(float64))", a, opStr, b))
	case 0x63, 0x67, 0x6b, 0x6f: // fastore? no: fadd,fsub,fmul,fdiv (cat1)
		b, a := popRef(), popRef()
		opStr := map[byte]string{0x63: "+", 0x67: "-", 0x6b: "*", 0x6f: "/"}[op]
		pushV(fmt.Sprintf("(%s.(float32)) %s (%s.(float32))", a, opStr, b))
	case 0x76: // fneg
		a := popRef()
		pushV(fmt.Sprintf("-(%s.(float32))", a))
	case 0x77: // dneg
		a := popRefCat2()
		pushCat2(fmt.Sprintf("-(%s.(float64))", a))

	case 0x37: // lstore (indexed, cat2)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(idx), popRefCat2())
	case 0x57: // pop
		e.depth--
	case 0x58: // pop2
		e.depth -= 2
	case 0x59: // dup
		top := e.peekTop()
		pushV(top)
	case 0x5a: // dup_x1  [B,T] -> [T,B,T]
		d := e.depth
		e.p("s%d = s%d", d, d-1)
		e.p("s%d = s%d", d-1, d-2)
		e.p("s%d = s%d", d-2, d)
		e.depth++
	case 0x5b: // dup_x2 [C,B,T] -> [T,C,B,T]
		d := e.depth
		e.p("s%d = s%d", d, d-1)
		e.p("s%d = s%d", d-1, d-2)
		e.p("s%d = s%d", d-2, d-3)
		e.p("s%d = s%d", d-3, d)
		e.depth++
	case 0x5c: // dup2 [B,T] -> [B,T,B,T]
		d := e.depth
		e.p("s%d = s%d", d, d-2)
		e.p("s%d = s%d", d+1, d-1)
		e.depth += 2
	case 0x5d: // dup2_x1 [C,B,T] -> [B,T,C,B,T]
		d := e.depth
		e.p("s%d = s%d", d, d-2)
		e.p("s%d = s%d", d+1, d-1)
		e.p("s%d = s%d", d-1, d-3)
		e.p("s%d = s%d", d-2, d)
		e.p("s%d = s%d", d-3, d+1)
		e.depth += 2
	case 0x5e: // dup2_x2 [D,C,B,T] -> [B,T,D,C,B,T]
		d := e.depth
		e.p("s%d = s%d", d, d-2)
		e.p("s%d = s%d", d+1, d-1)
		e.p("s%d = s%d", d-1, d-4)
		e.p("s%d = s%d", d-2, d-3)
		e.p("s%d = s%d", d-3, d)
		e.p("s%d = s%d", d-4, d+1)
		e.depth += 2

	case 0x60: // iadd
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) + (%s.(int32))", a, b))
	case 0x64:
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) - (%s.(int32))", a, b))
	case 0x68:
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) * (%s.(int32))", a, b))
	case 0x74:
		v := popRef()
		pushV(fmt.Sprintf("-(%s.(int32))", v))

	case 0x6c: // idiv
		b, a := popRef(), popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.IDiv(thr, %s.(int32), %s.(int32))", dst, a, b)
		excAfter()
		e.depth++
	case 0x70: // irem
		b, a := popRef(), popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.IRem(thr, %s.(int32), %s.(int32))", dst, a, b)
		excAfter()
		e.depth++

	case 0x78: // ishl
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32) << (%s.(int32) & 31))", a, b))
	case 0x79: // lshl (long value, int shift)
		b := popRef()
		a := popRefCat2()
		pushCat2(fmt.Sprintf("((%s.(int64)) << (%s.(int32) & 63))", a, b))
	case 0x7a: // ishr
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32) >> (%s.(int32) & 31))", a, b))
	case 0x7b: // lshr
		b := popRef()
		a := popRefCat2()
		pushCat2(fmt.Sprintf("((%s.(int64)) >> (%s.(int32) & 63))", a, b))
	case 0x7c: // iushr
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("int32(uint32(%s.(int32)) >> (uint32(%s.(int32)) & 31))", a, b))
	case 0x7d: // lushr
		b := popRef()
		a := popRefCat2()
		pushCat2(fmt.Sprintf("int64(uint64(%s.(int64)) >> (uint32(%s.(int32)) & 63))", a, b))

	case 0x61: // ladd
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) + (%s.(int64))", a, b))
	case 0x65: // lsub
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) - (%s.(int64))", a, b))
	case 0x69: // lmul
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) * (%s.(int64))", a, b))
	case 0x6d: // ldiv
		b, a := popRefCat2(), popRefCat2()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.LDiv(thr, %s.(int64), %s.(int64))", dst, a, b)
		excAfter()
		e.depth += 2
	case 0x71: // lrem
		b, a := popRefCat2(), popRefCat2()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.LRem(thr, %s.(int64), %s.(int64))", dst, a, b)
		excAfter()
		e.depth += 2
	case 0x75: // lneg
		a := popRefCat2()
		pushCat2(fmt.Sprintf("-(%s.(int64))", a))
	case 0x7f: // land
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) & (%s.(int64))", a, b))
	case 0x81: // lor
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) | (%s.(int64))", a, b))
	case 0x83: // lxor
		b, a := popRefCat2(), popRefCat2()
		pushCat2(fmt.Sprintf("(%s.(int64)) ^ (%s.(int64))", a, b))

	case 0x84: // iinc (non-wide: u1 index + s1 const)
		idx := int(code[pc+1])
		cst := int(int8(code[pc+2]))
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s.(int32) + %d", e.LocalName(idx), e.LocalName(idx), cst)

	case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e: // ifeq..ifle
		v := popRef()
		tgt := branchTarget(pc, code)
		cond := map[byte]string{
			0x99: "== 0", 0x9a: "!= 0", 0x9b: "< 0",
			0x9c: ">= 0", 0x9d: "> 0", 0x9e: "<= 0",
		}[op]
		e.p("if %s.(int32) %s { goto L%d }", v, cond, tgt)
	case 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4: // if_icmpXX
		b, a := popRef(), popRef()
		tgt := branchTarget(pc, code)
		cond := map[byte]string{
			0x9f: "==", 0xa0: "!=", 0xa1: "<",
			0xa2: ">=", 0xa3: ">", 0xa4: "<=",
		}[op]
		e.p("if %s.(int32) %s %s.(int32) { goto L%d }", a, cond, b, tgt)
	case 0xa5, 0xa6: // if_acmpeq / ne
		b, a := popRef(), popRef()
		tgt := branchTarget(pc, code)
		if op == 0xa6 {
			e.p("if !genrt.RefEq(%s, %s) { goto L%d }", a, b, tgt)
		} else {
			e.p("if genrt.RefEq(%s, %s) { goto L%d }", a, b, tgt)
		}
	case 0xc6, 0xc7: // ifnull / ifnonnull
		v := popRef()
		tgt := branchTarget(pc, code)
		wantNil := op == 0xc6
		e.p("if (%s == nil) == %v { goto L%d }", v, wantNil, tgt)

	case 0xa7: // goto
		tgt := branchTarget(pc, code)
		e.p("goto L%d", tgt)

	case 0xaa: // tableswitch
		p := pc + 1
		for p%4 != 0 {
			p++
		}
		rd32 := func(i int) int32 {
			return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
				uint32(code[i+2])<<8 | uint32(code[i+3]))
		}
		def := int(rd32(p))
		low := int(rd32(p + 4))
		high := int(rd32(p + 8))
		p += 12
		k := popRef()
		e.p("switch %s.(int32) {", k)
		for v := low; v <= high; v++ {
			tgt := pc + int(rd32(p))
			e.p("case %d:", v)
			e.p("	goto L%d", tgt)
			p += 4
		}
		e.p("default:")
		e.p("	goto L%d", pc+def)
		e.p("}")
	case 0xab: // lookupswitch
		p := pc + 1
		for p%4 != 0 {
			p++
		}
		rd32 := func(i int) int32 {
			return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
				uint32(code[i+2])<<8 | uint32(code[i+3]))
		}
		def := int(rd32(p))
		npairs := int(rd32(p + 4))
		p += 8
		k := popRef()
		e.p("switch %s.(int32) {", k)
		for j := 0; j < npairs; j++ {
			match := int(rd32(p))
			tgt := pc + int(rd32(p+4))
			e.p("case %d:", match)
			e.p("	goto L%d", tgt)
			p += 8
		}
		e.p("default:")
		e.p("	goto L%d", pc+def)
		e.p("}")

	case 0x7e: // iand
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) & (%s.(int32))", a, b))
	case 0x80: // ior
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) | (%s.(int32))", a, b))
	case 0x82: // ixor
		b, a := popRef(), popRef()
		pushV(fmt.Sprintf("(%s.(int32)) ^ (%s.(int32))", a, b))

	case 0x94: // lcmp (two cat2 -> int)
		b, a := popRefCat2(), popRefCat2()
		pushV(fmt.Sprintf("genrt.LCmp(%s.(int64), %s.(int64))", a, b))
	case 0x95, 0x96: // fcmpl, fcmpg
		b, a := popRef(), popRef()
		fn := "genrt.FCmpl"
		if op == 0x96 {
			fn = "genrt.FCmpg"
		}
		pushV(fmt.Sprintf("%s(%s.(float32), %s.(float32))", fn, a, b))
	case 0x97, 0x98: // dcmpl, dcmpg
		b, a := popRefCat2(), popRefCat2()
		fn := "genrt.DCmpl"
		if op == 0x98 {
			fn = "genrt.DCmpg"
		}
		pushV(fmt.Sprintf("%s(%s.(float64), %s.(float64))", fn, a, b))

	case 0x85: // i2l
		a := popRef()
		pushCat2(fmt.Sprintf("int64(%s.(int32))", a))
	case 0x86: // i2f
		a := popRef()
		pushV(fmt.Sprintf("float32(%s.(int32))", a))
	case 0x87: // i2d
		a := popRef()
		pushCat2(fmt.Sprintf("float64(%s.(int32))", a))
	case 0x88: // l2i
		a := popRefCat2()
		pushV(fmt.Sprintf("int32(%s.(int64))", a))
	case 0x89: // l2f
		a := popRefCat2()
		pushV(fmt.Sprintf("float32(%s.(int64))", a))
	case 0x8a: // l2d
		a := popRefCat2()
		pushCat2(fmt.Sprintf("float64(%s.(int64))", a))
	case 0x8b: // f2i
		a := popRef()
		pushV(fmt.Sprintf("int32(%s.(float32))", a))
	case 0x8c: // f2l
		a := popRef()
		pushCat2(fmt.Sprintf("int64(%s.(float32))", a))
	case 0x8d: // f2d
		a := popRef()
		pushCat2(fmt.Sprintf("float64(%s.(float32))", a))
	case 0x8e: // d2i
		a := popRefCat2()
		pushV(fmt.Sprintf("int32(%s.(float64))", a))
	case 0x8f: // d2l
		a := popRefCat2()
		pushCat2(fmt.Sprintf("int64(%s.(float64))", a))
	case 0x90: // d2f
		a := popRefCat2()
		pushV(fmt.Sprintf("float32(%s.(float64))", a))
	case 0x91: // i2b
		a := popRef()
		pushV(fmt.Sprintf("int32(int8(%s.(int32)))", a))
	case 0x92: // i2c
		a := popRef()
		pushV(fmt.Sprintf("%s.(int32) & 0xFFFF", a))
	case 0x93: // i2s
		a := popRef()
		pushV(fmt.Sprintf("int32(int16(%s.(int32)))", a))

	case 0xac, 0xae, 0xb0: // ireturn/freturn/areturn (cat1)
		v := popRef()
		e.p("return %s, nil", v)
	case 0xad, 0xaf: // lreturn/dreturn (cat2)
		e.p("return %s, nil", popRefCat2())
	case 0xb1: // return
		e.p("return nil, nil")

	case 0xb2: // getstatic
		cls, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		if descSlots(desc) == 2 {
			pushCat2(fmt.Sprintf("genrt.GetStatic(thr, %q, %q, %q)", cls, name, desc))
		} else {
			pushV(fmt.Sprintf("genrt.GetStatic(thr, %q, %q, %q)", cls, name, desc))
		}
	case 0xb3: // putstatic
		cls, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		var val string
		if descSlots(desc) == 2 {
			val = popRefCat2()
		} else {
			val = popRef()
		}
		e.p("genrt.SetStatic(thr, %q, %q, %q, %s)", cls, name, desc, val)
	case 0xb4: // getfield
		_, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		obj := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.GetFieldChecked(thr, %s, %q, %q)", dst, obj, name, desc)
		excAfter()
		e.depth += descSlots(desc)
	case 0xb5: // putfield
		_, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		var val string
		if descSlots(desc) == 2 {
			val = popRefCat2() // cat2 value occupies two raw slots
		} else {
			val = popRef()
		}
		obj := popRef()
		e.p("exc = genrt.SetFieldChecked(thr, %s, %q, %q, %s)", obj, name, desc, val)
		excAfter()

	case 0xb6, 0xb7, 0xb8, 0xb9: // invokes
		cls, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		if op == 0xb9 {
		}
		argDescs, _, derr := splitMethodDesc(desc)
		if derr != nil {
			return derr
		}
		total := 0
		starts := make([]int, len(argDescs))
		for i, ad := range argDescs {
			starts[i] = total
			total += descSlots(ad)
		}
		base := e.depth - total
		argNames := make([]string, len(argDescs))
		for i := range argDescs {
			argNames[i] = fmt.Sprintf("s%d", base+starts[i])
		}
		e.depth -= total
		var recv string
		if op != 0xb8 {
			recv = popRef()
		}
		fnName := map[byte]string{
			0xb6: "genrt.CallVirtual",
			0xb7: "genrt.CallSpecial",
			0xb8: "genrt.CallStatic",
			0xb9: "genrt.CallVirtual",
		}[op]
		hasRet, rerr2 := returnsValue(desc)
		if rerr2 != nil {
			return rerr2
		}
		retDesc := ""
		if idx := strings.LastIndex(desc, ")"); idx >= 0 {
			retDesc = desc[idx+1:]
		}
		retSlots := descSlots(retDesc) // 0 for void, 2 for J/D, else 1
		dst := "_"
		retSlot := -1
		if hasRet {
			retSlot = e.depth
			dst = fmt.Sprintf("s%d", retSlot)
		}
		// Args expression: reusable frame buffer when sized for it
		// (kills per-call heap slices); heap literal otherwise.
		argsExpr := "[]kernel.Value{" + joinStrings(argNames, ", ") + "}"
		if e.maxArgs > 0 && len(argNames) <= e.maxArgs {
			for i, an := range argNames {
				e.p("abuf[%d] = %s", i, an)
			}
			argsExpr = fmt.Sprintf("abuf[:%d]", len(argNames))
		}
		var call string
		if op == 0xb8 {
			call = fmt.Sprintf("genrt.CallStatic(thr, %q, %q, %q, %s)", cls, name, desc, argsExpr)
		} else if op == 0xb6 || op == 0xb9 {
			// Virtual/interface dispatch goes through the monomorphic
			// inline cache; the slot index is baked at emission time.
			// genrt.ICSlot must stay identical to icHash (pinned by
			// TestEmitterICSlotAgreement).
			call = fmt.Sprintf("genrt.CallVirtualIC(%d, thr, %s, %q, %q, %q, %s)",
				icSlotFor(cls, name, desc), recv, cls, name, desc, argsExpr)
		} else {
			call = fmt.Sprintf("genrt.Call%s(thr, %s, %q, %q, %q, %s)",
				invokeKindName(op), recv, cls, name, desc, argsExpr)
		}
		if os.Getenv("CATTY_EMITDBG") != "" && strings.Contains(name, "isNaN") {
			fmt.Fprintf(os.Stderr, "[emit-isa] depth=%d base=%d argNames=%v\n", e.depth, base, argNames)
		}
		e.p("%s, exc = %s", dst, call)
		_ = fnName
		excAfter()
		if hasRet {
			e.depth += retSlots
		}

	case 0xbb: // new
		clsName, cerr := classAt(e.cf, code, pc)
		if cerr != nil {
			return cerr
		}
		pushV(fmt.Sprintf("genrt.New(thr, %q)", clsName))
	case 0xbc: // newarray
		nm, ok := atypeName(int(code[pc+1]))
		if !ok {
			return fail("newarray bad atype %d", code[pc+1])
		}
		sizeExpr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.NewPrimitiveArray(thr, %q, %s.(int32))", dst, nm, sizeExpr)
		excAfter()
		e.depth++
	case 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35: // array loads
		idx := popRef()
		arr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.ALoadChecked(thr, %s, %s.(int32))", dst, arr, idx)
		excAfter()
		switch op {
		case 0x2f, 0x31: // laload, daload: cat2 element
			e.depth++
		case 0x33: // baload sign-extend
			e.p("%s = int32(int8(%s.(int32)))", dst, dst)
		case 0x34: // caload zero-extend
			e.p("%s = %s.(int32) & 0xFFFF", dst, dst)
		case 0x35: // saload sign-extend
			e.p("%s = int32(int16(%s.(int32)))", dst, dst)
		}
		e.depth++

	case 0x4f, 0x51, 0x53, 0x54, 0x55, 0x56: // cat1 array stores
		val := popRef()
		idx := popRef()
		arr := popRef()
		e.p("exc = genrt.AStoreChecked(thr, %s, %s.(int32), %s)", arr, idx, val)
		excAfter()
	case 0x50, 0x52: // lastore, dastore (cat2 value)
		val := popRefCat2()
		idx := popRef()
		arr := popRef()
		e.p("exc = genrt.AStoreChecked(thr, %s, %s.(int32), %s)", arr, idx, val)
		excAfter()

	case 0xbd: // anewarray
		clsName, cerr := classAt(e.cf, code, pc)
		if cerr != nil {
			return cerr
		}
		sizeExpr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.NewRefArray(thr, %q, %s.(int32))", dst, clsName, sizeExpr)
		excAfter()
		e.depth++

	case 0xbe: // arraylength
		arr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.ArrayLengthChecked(thr, %s)", dst, arr)
		excAfter()
		e.depth++
	case 0xbf: // athrow
		v := popRef()
		e.p("exc = &kernel.Thrown{Obj: %s.(*kernel.Instance)}", v)
		excAfter()
	case 0xc0: // checkcast
		clsName, cerr := classAt(e.cf, code, pc)
		if cerr != nil {
			return cerr
		}
		top := e.peekTop()
		e.p("%s = genrt.CheckCast(thr, %s, %q)", top, top, clsName)
		excAfter()
	case 0xc1: // instanceof
		clsName, cerr := classAt(e.cf, code, pc)
		if cerr != nil {
			return cerr
		}
		v := popRef()
		pushV(fmt.Sprintf("genrt.BoolValue(genrt.InstanceOf(%s, %q))", v, clsName))
	case 0xc2, 0xc3: // monitorenter/exit
		obj := popRef()
		helper := "genrt.MonitorEnter"
		if op == 0xc3 {
			helper = "genrt.MonitorExit"
		}
		e.p("exc = %s(thr, %s)", helper, obj)
		excAfter()

	default:
		return fail("opcode %#x unsupported by emitter v1", op)
	}

	return nil
}

func invokeKindName(op byte) string {
	switch op {
	case 0xb6:
		return "Virtual"
	case 0xb7:
		return "Special"
	case 0xb9:
		return "Interface"
	}
	return "Static"
}

// icSlotFor mirrors genrt.icHash (FNV-1a over cls/name/desc mod table
// size). The two MUST stay identical — pinned cross-package by
// TestEmitterICSlotAgreement in internal/genrt.
func icSlotFor(cls, name, desc string) uint32 {
	h := uint32(2166136261)
	for _, s := range [3]string{cls, name, desc} {
		for i := 0; i < len(s); i++ {
			h ^= uint32(s[i])
			h *= 16777619
		}
	}
	return h % 1024
}
