package rt

import (
	"fmt"
	"os"
	"strconv"
)

// intrinsicObjectInit implements Object.<init> and every bootstrap throwable
// constructor: allocation happened at `new`, so construction is a no-op that
// leaves the receiver as-is.
func intrinsicObjectInit(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return Normal()
}

// intrinsicPrintlnInt implements PrintStream.println(I)V with signed decimal
// output and a newline.
func intrinsicPrintlnInt(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	_, _ = os.Stdout.WriteString(strconv.FormatInt(int64(AsInt(args[0])), 10) + "\n")
	return Normal()
}

// intrinsicThrowableGetCause implements Throwable.getCause() and the
// InvocationTargetException override. It returns Java null when no cause was
// recorded.
func intrinsicThrowableGetCause(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	return OkValue(RefValue(CauseOf(receiver)))
}

// intrinsicIntegerValueOf implements Integer.valueOf(I) with the profile's
// supported boxing identity guarantee.
func intrinsicIntegerValueOf(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(RefValue(t.RT.BoxedInteger(AsInt(args[0]))))
}

// intrinsicIntegerIntValue implements Integer.intValue()I.
func intrinsicIntegerIntValue(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	v, ok := IntegerValue(receiver)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	return OkValue(IntValue(v))
}

// intrinsicClassGetMethod implements
// Class.getMethod(String, Class[])Method.
func intrinsicClassGetMethod(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	c, ok := receiver.Payload.(*Class)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	if len(args) != 2 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	name, ok := StringText(args[0].Ref)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	var params []*Class
	if args[1].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[1].Ref != nil {
		arr := args[1].Ref.Array
		if arr == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		for _, elem := range arr.Elems {
			if elem.Ref == nil {
				return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
			}
			pc, ok := elem.Ref.Payload.(*Class)
			if !ok {
				return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
			}
			params = append(params, pc)
		}
	}
	found, thrown := t.RT.lookupPublicMethod(c, name, params)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	methodObj := t.RT.NewInstance(t.RT.mustClass("java/lang/reflect/Method"))
	methodObj.Payload = found
	return OkValue(RefValue(methodObj))
}

// lookupPublicMethod searches public declarations, then inherited public
// methods. Constructors and class initializers are excluded.
func (r *Runtime) lookupPublicMethod(c *Class, name string, params []*Class) (*Method, *Object) {
	if m := matchPublicMethod(c, name, params); m != nil {
		return m, nil
	}
	for cur := c.Super; cur != nil; cur = cur.Super {
		if m := matchPublicMethod(cur, name, params); m != nil {
			return m, nil
		}
	}
	if m := matchInterfaceMethod(c, name, params); m != nil {
		return m, nil
	}
	return nil, r.bootstrapThrowable("java/lang/NoSuchMethodException")
}

func matchPublicMethod(c *Class, name string, params []*Class) *Method {
	for _, m := range c.Methods {
		if m.Name != name || m.Name == "<init>" || m.Name == "<clinit>" {
			continue
		}
		if m.Flags&AccPublic == 0 {
			continue
		}
		if len(m.ParamClass) != len(params) {
			continue
		}
		ok := true
		for i := range params {
			if m.ParamClass[i] != params[i] {
				ok = false
				break
			}
		}
		if ok {
			return m
		}
	}
	return nil
}

func matchInterfaceMethod(c *Class, name string, params []*Class) *Method {
	seen := map[*Class]bool{}
	var visit func(*Class) *Method
	visit = func(c *Class) *Method {
		if c == nil || seen[c] {
			return nil
		}
		seen[c] = true
		for _, in := range c.Interfaces {
			if m := matchPublicMethod(in, name, params); m != nil {
				return m
			}
			if m := visit(in); m != nil {
				return m
			}
		}
		return visit(c.Super)
	}
	return visit(c)
}

// intrinsicMethodInvoke implements Method.invoke(Object, Object[])Object,
// including argument conversion, virtual selection, boxing and the
// InvocationTargetException wrapper.
func intrinsicMethodInvoke(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	decl, ok := receiver.Payload.(*Method)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	if len(args) != 2 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	realReceiver := args[0].Ref
	var rawArgs []Value
	if args[1].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[1].Ref != nil {
		if args[1].Ref.Array == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		rawArgs = args[1].Ref.Array.Elems
	}
	converted, thrown := convertReflectiveArgs(t, decl, rawArgs)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if decl.Flags&AccStatic == 0 {
		if realReceiver == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
		}
		if !t.RT.IsAssignable(realReceiver.Class, decl.Owner) {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
	}
	var target *Method
	if decl.Flags&AccStatic != 0 {
		if thrown := t.EnsureInitialized(decl.Owner); thrown != nil {
			return ThrowOutcome(thrown)
		}
		target = decl
	} else {
		target, thrown = t.RT.selectVirtual(realReceiver.Class, decl)
		if thrown != nil {
			return ThrowOutcome(thrown)
		}
	}
	if tr := t.trace(); tr != nil {
		tr.reflect(t.ID, decl, target)
	}
	var callReceiver *Object
	if decl.Flags&AccStatic == 0 {
		callReceiver = realReceiver
	}
	out := t.Call(target, callReceiver, converted)
	if out.Thrown != nil {
		wrapper := t.RT.NewThrowable(t, "java/lang/reflect/InvocationTargetException")
		SetCause(wrapper, out.Thrown)
		if tr := t.trace(); tr != nil {
			tr.wrap(t.ID, wrapper, out.Thrown)
		}
		return ThrowOutcome(wrapper)
	}
	switch decl.Type.Return.Kind {
	case 'V':
		return OkValue(RefValue(nil))
	case 'I':
		return OkValue(RefValue(t.RT.BoxedInteger(AsInt(out.Value))))
	case 'L', 'A':
		return OkValue(out.Value)
	default:
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
}

func convertReflectiveArgs(t *Thread, decl *Method, raw []Value) ([]Value, *Object) {
	if len(raw) != len(decl.Type.Params) {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
	}
	out := make([]Value, len(raw))
	for i, p := range decl.Type.Params {
		switch p.Kind {
		case 'I':
			v, ok := IntegerValue(raw[i].Ref)
			if !ok {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			out[i] = IntValue(v)
		case 'L', 'A':
			if raw[i].Kind != KindRef {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			if raw[i].Ref == nil {
				out[i] = RefValue(nil)
				continue
			}
			pc := decl.ParamClass[i]
			if pc == nil || !t.RT.IsAssignable(raw[i].Ref.Class, pc) {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			out[i] = raw[i]
		default:
			return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
		}
	}
	return out, nil
}

func (r *Runtime) StringOfGo(s string) *Object { return r.NewStringGo(s) }

func (r *Runtime) debugMethod(m *Method) string {
	return fmt.Sprintf("%s.%s%s", m.Owner.Name, m.Name, m.Descriptor)
}
