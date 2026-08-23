package emitter

import (
	"fmt"

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
	pushV := func(expr string) { e.p("s%d = %s", e.depth, expr); e.depth++ }
	excAfter := func() { e.excDispatch(pc) }
	pushCat2 := func(expr string) {
		e.p("s%d = %s", e.depth, expr)
		e.depth += 2
	}

	switch op {

	case 0x00: // nop

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
	case 0x3f, 0x40, 0x41, 0x42: // lstore_0..3
		n := int(op - 0x3f)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRef())
	case 0x43, 0x44, 0x45, 0x46: // fstore_0..3
		n := int(op - 0x43)
		if err := e.checkLocalIdx(n); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(n), popRef())
	case 0x47, 0x48, 0x49, 0x4a: // dstore_0..3
		n := int(op - 0x47)
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

	case 0x37: // lstore (indexed)
		idx := int(code[pc+1])
		if err := e.checkLocalIdx(idx); err != nil {
			return fail("%v", err)
		}
		e.p("%s = %s", e.LocalName(idx), popRef())
	case 0x57: // pop
		e.depth--
	case 0x59: // dup
		top := e.peekTop()
		pushV(top)

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
		e.p("%s, exc = genrt.IDiv(%s.(int32), %s.(int32))", dst, a, b)
		excAfter()
		e.depth++
	case 0x70: // irem
		b, a := popRef(), popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.IRem(%s.(int32), %s.(int32))", dst, a, b)
		excAfter()
		e.depth++

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
		opStr := "=="
		if op == 0xa6 {
			opStr = "!="
		}
		e.p("if genrt.RefEq(%s, %s) %s { goto L%d }", a, b, opStr, tgt)
	case 0xc6, 0xc7: // ifnull / ifnonnull
		v := popRef()
		tgt := branchTarget(pc, code)
		wantNil := op == 0xc6
		e.p("if (%s == nil) == %v { goto L%d }", v, wantNil, tgt)
	case 0xa7: // goto
		tgt := branchTarget(pc, code)
		e.p("goto L%d", tgt)

	case 0xac, 0xb0: // ireturn / areturn
		v := popRef()
		e.p("return %s, nil", v)
	case 0xb1: // return
		e.p("return nil, nil")

	case 0xb2: // getstatic
		cls, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		pushV(fmt.Sprintf("genrt.GetStatic(thr, %q, %q, %q)", cls, name, desc))
	case 0xb3: // putstatic
		cls, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		val := popRef()
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
		e.depth++
	case 0xb5: // putfield
		_, name, desc, _, rerr := refAt(e.cf, code, pc)
		if rerr != nil {
			return rerr
		}
		val := popRef()
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
		argsExpr := "[]kernel.Value{" + joinStrings(argNames, ", ") + "}"
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
		dst := "_"
		retSlot := -1
		if hasRet {
			retSlot = e.depth
			dst = fmt.Sprintf("s%d", retSlot)
		}
		var call string
		if op == 0xb8 {
			call = fmt.Sprintf("genrt.CallStatic(thr, %q, %q, %q, %s)", cls, name, desc, argsExpr)
		} else {
			call = fmt.Sprintf("genrt.Call%s(thr, %s, %q, %q, %q, %s)",
				invokeKindName(op), recv, cls, name, desc, argsExpr)
		}
		e.p("%s, exc = %s", dst, call)
		_ = fnName
		excAfter()
		if hasRet {
			e.depth++
		}

	case 0xbb: // new
		clsName, cerr := classAt(e.cf, code, pc)
		if cerr != nil {
			return cerr
		}
		pushV(fmt.Sprintf("genrt.New(%q)", clsName))
	case 0xbc: // newarray
		nm, ok := atypeName(int(code[pc+1]))
		if !ok {
			return fail("newarray bad atype %d", code[pc+1])
		}
		sizeExpr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.NewPrimitiveArray(%q, %s.(int32))", dst, nm, sizeExpr)
		excAfter()
		e.depth++
	case 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35: // array loads
		idx := popRef()
		arr := popRef()
		dst := fmt.Sprintf("s%d", e.depth)
		e.p("%s, exc = genrt.ALoadChecked(%s, %s.(int32))", dst, arr, idx)
		excAfter()
		switch op {
		case 0x33: // baload sign-extend
			e.p("%s = int32(int8(%s.(int32)))", dst, dst)
		case 0x34: // caload zero-extend
			e.p("%s = %s.(int32) & 0xFFFF", dst, dst)
		case 0x35: // saload sign-extend
			e.p("%s = int32(int16(%s.(int32)))", dst, dst)
		}

	case 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56: // array stores
		val := popRef()
		idx := popRef()
		arr := popRef()
		e.p("exc = genrt.AStoreChecked(%s, %s.(int32), %s)", arr, idx, val)
		excAfter()

	case 0xbe: // arraylength
		arr := popRef()
		pushV(fmt.Sprintf("genrt.ArrayLength(%s)", arr))
		excAfter()
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
