package verify

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/classfile"
)

// checker carries per-class verification context.
type checker struct {
	cf *classfile.ClassFile
	r  RefResolver
	mn string // method name+desc for diagnostics
}

// throwSuccs enqueues exception-handler entries covering pc. The pushed
// exception type is the catch type when known (handler frames pin it);
// catch-all handlers receive the concrete thrown type.
func (c *checker) throwSuccs(m *classfile.MethodInfo, pc int, st *frame, exc vtype) []successor {
	var out []successor
	for _, h := range m.Code.Handlers {
		if pc < int(h.StartPc) || pc >= int(h.EndPc) {
			continue
		}
		excType := exc
		if h.CatchType != 0 {
			catchName, err := c.cf.ClassName(h.CatchType)
			if err != nil {
				continue
			}
			if exc.kind == kObj && !c.refCompatible(exc.class, catchName) {
				continue // cannot catch
			}
			excType = tObj(catchName)
		}
		hs := st.clone()
		hs.clearStack()
		hs.push(excType)
		out = append(out, successor{pc: int(h.HandlerPc), st: hs})
	}
	return out
}

// simulate routes to the table effect or an irregular handler.
func (c *checker) simulate(m *classfile.MethodInfo, pc int, st *frame) ([]successor, error) {
	op := m.Code.Code[pc]
	if os.Getenv("CATTY_VDBG") != "" && strings.Contains(c.mn, "readArray") && pc >= 10 && pc <= 22 {
		fmt.Fprintf(os.Stderr, "[stk] pc=%d op=%#x len=%d\n", pc, op, len(st.stack))
	}
	if f, ok := effects[op]; ok {
		return f(c, m, pc, st)
	}
	return c.simulateIrregular(m, pc, st, op)
}

func (c *checker) simulateIrregular(m *classfile.MethodInfo, pc int, st *frame, op byte) ([]successor, error) {
	code := m.Code.Code
	cu := cursor{code: code, pc: pc + 1}
	fail := func(format string, args ...any) error { return verr(c.mn, pc, format, args...) }
	next := func() []successor { return fall(st, cu.pc) }

	switch op {

	case 0x12, 0x13: // ldc / ldc_w
		var idx uint16
		if op == 0x12 {
			idx = uint16(cu.u1())
		} else {
			idx = cu.u2()
		}
		e, err := c.cf.Entry(idx)
		if err != nil {
			return nil, fail("%v", err)
		}
		switch e.Tag {
		case classfile.CInteger:
			st.push(tInt)
		case classfile.CFloat:
			st.push(tFloat)
		case classfile.CString:
			st.push(tObj("java/lang/String"))
		case classfile.CClass:
			st.push(tObj("java/lang/Class"))
		default:
			return nil, fail("ldc on %s", e.Tag)
		}
		return next(), nil

	case 0x14: // ldc2_w
		e, err := c.cf.Entry(cu.u2())
		if err != nil {
			return nil, err
		}
		switch e.Tag {
		case classfile.CLong:
			st.push(tLong)
		case classfile.CDouble:
			st.push(tDouble)
		default:
			return nil, fail("ldc2_w on %s", e.Tag)
		}
		return next(), nil

	case 0xb2: // getstatic
		_, _, desc, _, err := c.cf.Ref(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		st.push(vtypeForDesc(desc))
		return next(), nil

	case 0xb3: // putstatic
		_, _, desc, _, err := c.cf.Ref(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		if vtypeForDesc(desc).kind == kLong || vtypeForDesc(desc).kind == kDouble {
			if _, err := st.popCat2(); err != nil {
				return nil, fail("%v", err)
			}
		} else if _, err := st.pop(); err != nil {
			return nil, fail("%v", err)
		}
		return next(), nil

	case 0xb4: // getfield
		_, _, desc, _, err := c.cf.Ref(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		if _, rerr := st.popExpect(kObj, "getfield receiver"); rerr != nil {
			return nil, fail("%v", rerr)
		}
		st.push(vtypeForDesc(desc))
		return next(), nil

	case 0xb5: // putfield
		_, _, desc, _, err := c.cf.Ref(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		want := vtypeForDesc(desc)
		if want.kind == kLong || want.kind == kDouble {
			if _, perr := st.popCat2(); perr != nil {
				return nil, fail("%v", perr)
			}
		} else if _, perr := st.pop(); perr != nil {
			return nil, fail("%v", perr)
		}
		if _, rerr := st.popExpect(kObj, "putfield receiver"); rerr != nil {
			return nil, fail("%v", rerr)
		}
		return next(), nil

	case 0xb6, 0xb7, 0xb8, 0xb9: // invokevirtual/special/static/interface
		clsName, name, desc, _, err := c.cf.Ref(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		if op == 0xb9 {
			cu.u1() // count
			cu.u1() // zero
		}
		argDescs, _, err := splitDesc(desc)
		if err != nil {
			return nil, fail("%v", err)
		}
		// pop args reversed
		for i := len(argDescs) - 1; i >= 0; i-- {
			want := vtypeForDesc(argDescs[i])
			if want.kind == kLong || want.kind == kDouble {
				if _, perr := st.popCat2(); perr != nil {
					return nil, fail("arg[%d]: %v", i, perr)
				}
			} else {
				v, perr := st.pop()
				if perr != nil {
					return nil, fail("arg[%d]: %v", i, perr)
				}
				if want.kind == kObj && !isRef(v) {
					return nil, fail("arg[%d]: want reference, got %s", i, v)
				}
				if want.kind != kObj && v.kind != want.kind && v.kind != kUnknown {
					return nil, fail("arg[%d]: want %s, got %s", i, want, v)
				}
			}
		}
		isStatic := op == 0xb8
		isCtor := op == 0xb7 && name == "<init>"
		var recv vtype
		if !isStatic {
			rv, rerr := st.pop()
			if rerr != nil {
				return nil, fail("receiver: %v", rerr)
			}
			recv = rv
			if !isRef(recv) {
				return nil, fail("receiver %s is not a reference", recv)
			}
		}
		// receiver assignability for virtual/interface calls (skip for
		// invokespecial: super/private/ctor rules are looser)
		if op == 0xb6 || op == 0xb9 {
			if recv.kind == kObj && !c.refCompatible(recv.class, clsName) {
				return nil, fail("receiver %s not a %s", recv.class, clsName)
			}
		}
		// <init> completion: initialize uninit markers. JVMS §4.10.1.6:
		// uninitialized(this) becomes the class UNDER VERIFICATION (not the
		// super constructor's declaring class) — minimal-json's
		// JsonParser.<init> exposed this via a loop merge on local 0.
		if isCtor {
			holder := tObj(c.cf.ThisClass)
			if recv.kind == kUninitThis {
				replaceUninitThis(st, holder)
			} else if recv.kind == kUninit {
				holder = tObj(clsName) // new X(...) sites keep their ref class
				replaceUninitAt(st, recv.off, holder)
			}
			if recv.kind == kUninitThis {
				replaceUninitThis(st, holder)
			} else if recv.kind == kUninit {
				replaceUninitAt(st, recv.off, holder)
			}
		}
		ret, hasRet, err := retVType(desc)
		if err != nil {
			return nil, fail("%v", err)
		}
		if hasRet {
			st.push(ret)
		}
		return next(), nil

	case 0xbb: // new
		cu.u2()
		st.push(tUninit(uint16(pc)))
		return fall(st, pc+3), nil

	case 0xbc: // newarray
		at := cu.u1()
		if _, ok := atypeName(at); !ok {
			return nil, fail("bad newarray atype %d", at)
		}
		if _, err := st.popExpect(kInt, "array size"); err != nil {
			return nil, fail("%v", err)
		}
		nm, _ := atypeName(at)
		st.push(tObj("[" + nm))
		succs := c.throwSuccs(m, pc, st, tObj("java/lang/NegativeArraySizeException"))
		return append(succs, fall(st, pc+2)...), nil

	case 0xbd: // anewarray
		compCls, err := c.cf.ClassName(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		if _, err := st.popExpect(kInt, "array size"); err != nil {
			return nil, fail("%v", err)
		}
		st.push(tObj(arrayDescOf(compCls)))
		succs := c.throwSuccs(m, pc, st, tObj("java/lang/NegativeArraySizeException"))
		return append(succs, fall(st, pc+3)...), nil

	case 0xbe: // arraylength
		if _, err := st.popExpect(kObj, "arrayref"); err != nil {
			return nil, fail("%v", err)
		}
		succs := c.throwSuccs(m, pc, st, tObj("java/lang/NullPointerException"))
		st.push(tInt)
		next := fall(st, pc+1)
		return append(succs, next...), nil

	case 0xbf: // athrow
		v, err := st.popExpect(kObj, "athrow operand")
		if err != nil {
			return nil, fail("%v", err)
		}
		if v.kind == kObj && !c.refCompatible(v.class, "java/lang/Throwable") {
			return nil, fail("athrow of non-Throwable %s", v.class)
		}
		exc := v
		if exc.kind != kObj && exc.kind != kUnknown {
			exc = tObj("java/lang/Throwable")
		}
		return c.throwSuccs(m, pc, st, exc), nil

	case 0xc0: // checkcast
		clsName, err := c.cf.ClassName(cu.u2())
		if err != nil {
			return nil, fail("%v", err)
		}
		v := st.stack[len(st.stack)-1]
		if !isRef(v) {
			return nil, fail("checkcast on %s", v)
		}
		if v.kind == kObj {
			from := v.class
			// A cast is rejected only when provably impossible: both sides
			// known, unrelated in BOTH directions, target a class (casting
			// to an interface can always succeed for some subclass).
			if c.r.Known(from) && c.r.Known(clsName) && !c.r.IsInterface(clsName) &&
				!c.refCompatible(from, clsName) && !c.refCompatible(clsName, from) {
				return nil, fail("checkcast %s → %s impossible", from, clsName)
			}
		}
		// Successful checkcast refines the operand type (JVMS §4.10.1).
		st.stack[len(st.stack)-1] = vtype{kind: kObj, class: clsName}
		return next(), nil

	case 0xc1: // instanceof
		if _, err := c.cf.ClassName(cu.u2()); err != nil {
			return nil, fail("%v", err)
		}
		if _, err := st.popExpect(kObj, "instanceof operand"); err != nil {
			return nil, fail("%v", err)
		}
		st.push(tInt)
		return next(), nil

	case 0xc2, 0xc3: // monitorenter/exit
		if _, err := st.popExpect(kObj, "monitor operand"); err != nil {
			return nil, fail("%v", err)
		}
		succs := c.throwSuccs(m, pc, st, tObj("java/lang/NullPointerException"))
		next := fall(st, pc+1)
		return append(succs, next...), nil

	case 0xc5: // multianewarray
		clsName, err := c.cf.ClassName(cu.u2())
		if err != nil {
			return nil, err
		}
		dims := int(cu.u1())
		if dims < 1 || dims > 255 {
			return nil, fail("multianewarray dims %d", dims)
		}
		for i := 0; i < dims; i++ {
			if _, err := st.popExpect(kInt, "dim"); err != nil {
				return nil, fail("%v", err)
			}
		}
		st.push(tObj(clsName))
		succs := c.throwSuccs(m, pc, st, tObj("java/lang/NegativeArraySizeException"))
		return append(succs, fall(st, cu.pc)...), nil

	case 0xaa: // tableswitch
		cu.align4()
		def := cu.s4()
		low := cu.s4()
		high := cu.s4()
		if high < low {
			return nil, fail("tableswitch low>high")
		}
		if _, err := st.popExpect(kInt, "switch operand"); err != nil {
			return nil, fail("%v", err)
		}
		targets := []successor{}
		addT := func(off int32) {
			targets = append(targets, successor{pc: pc + int(off), st: st.clone()})
		}
		addT(def)
		for v := low; v <= high; v++ {
			addT(cu.s4())
		}
		return targets, nil

	case 0xab: // lookupswitch
		cu.align4()
		def := cu.s4()
		npairs := int(cu.s4())
		if _, err := st.popExpect(kInt, "switch operand"); err != nil {
			return nil, fail("%v", err)
		}
		targets := []successor{}
		addT := func(off int32) {
			targets = append(targets, successor{pc: pc + int(off), st: st.clone()})
		}
		addT(def)
		for j := 0; j < npairs; j++ {
			cu.s4() // match value
			addT(cu.s4())
		}
		return targets, nil

	case 0xc4: // wide
		wop := cu.u1()
		idx := int(cu.u2())
		switch wop {
		case 0x15, 0x17, 0x19:
			if idx >= len(st.locals) {
				return nil, fail("wide load local[%d] out of range", idx)
			}
			v := st.locals[idx]
			if !isRef(v) && v.kind != kInt && v.kind != kFloat &&
				v.kind != kLong && v.kind != kDouble && v.kind != kUnknown {
				return nil, fail("wide load of %s", v)
			}
			st.push(v)
			return next(), nil
		case 0x36, 0x38, 0x3a:
			v, perr := st.pop()
			if perr != nil {
				return nil, fail("%v", perr)
			}
			setLocalPair(st, idx, v)
			return next(), nil
		case 0x37, 0x39:
			v, perr := st.popCat2()
			if perr != nil {
				return nil, fail("%v", perr)
			}
			setLocalPair(st, idx, v)
			return next(), nil
		case 0x84: // wide iinc: const is s2 (consumed; typing unaffected)
			cu.s2()
			if idx >= len(st.locals) {
				return nil, fail("wide iinc local[%d] out of range", idx)
			}
			if st.locals[idx].kind != kInt && st.locals[idx].kind != kUnknown {
				return nil, fail("wide iinc on %s", st.locals[idx])
			}
			return next(), nil
		case 0xa9:
			return nil, fail("jsr/ret unsupported (v51+)")
		default:
			return nil, fail("bad wide opcode %#x", wop)
		}

	case 0xa8, 0xa9: // jsr / ret
		return nil, fail("jsr/ret unsupported (prohibited in v51+)")

	case 0xba: // invokedynamic
		return nil, fail("invokedynamic unsupported (ADR-0002 build-time desugaring)")

	default:
		return nil, fail("unknown opcode %#x", op)
	}
}

func replaceUninitThis(st *frame, with vtype) {
	for i := range st.locals {
		if st.locals[i].kind == kUninitThis {
			st.locals[i] = with
		}
	}
	for i := range st.stack {
		if st.stack[i].kind == kUninitThis {
			st.stack[i] = with
		}
	}
}

func replaceUninitAt(st *frame, off uint16, with vtype) {
	for i := range st.locals {
		if st.locals[i].kind == kUninit && st.locals[i].off == off {
			st.locals[i] = with
		}
	}
	for i := range st.stack {
		if st.stack[i].kind == kUninit && st.stack[i].off == off {
			st.stack[i] = with
		}
	}
}

func atypeName(at byte) (string, bool) {
	switch at {
	case 4:
		return "Z", true
	case 5:
		return "C", true
	case 6:
		return "F", true
	case 7:
		return "D", true
	case 8:
		return "B", true
	case 9:
		return "S", true
	case 10:
		return "I", true
	case 11:
		return "J", true
	}
	return "", false
}

var _ = fmt.Sprintf

// arrayDescOf converts a component class name (slash form) to the array
// descriptor: L...; for references, single [ for each dimension already
// present, and primitive letters as-is.
func arrayDescOf(comp string) string {
	switch comp {
	case "int":
		return "[I"
	case "long":
		return "[J"
	case "float":
		return "[F"
	case "double":
		return "[D"
	case "boolean":
		return "[Z"
	case "byte":
		return "[B"
	case "char":
		return "[C"
	case "short":
		return "[S"
	}
	if len(comp) > 0 && comp[0] == '[' {
		return "[" + comp // component is itself an array type descriptor
	}
	return "[L" + comp + ";"
}
