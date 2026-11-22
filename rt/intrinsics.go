package rt

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// intrinsicObjectInit implements Object.<init> and every bootstrap throwable
// constructor: allocation happened at `new`, so construction is a no-op that
// leaves the receiver as-is.
func intrinsicObjectInit(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return Normal()
}

// intrinsicObjectGetClass implements Object.getClass.
func intrinsicObjectGetClass(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if receiver.Class == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	return OkValue(RefValue(receiver.Class.Mirror))
}

func writeJavaString(s string) {
	_, _ = os.Stdout.WriteString(s)
}

func writeJavaUnits(units []uint16) {
	_, _ = os.Stdout.WriteString(string(utf16Decode(units)))
}

// utf16Decode mirrors string(utf16.Decode(units)) without an extra import in
// every call site.
func utf16Decode(units []uint16) string { return UnitsToString(units) }

// intrinsicPrintlnInt implements PrintStream.println(I)V.
func intrinsicPrintlnInt(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	writeJavaString(strconv.FormatInt(int64(AsInt(args[0])), 10) + "\n")
	return Normal()
}

// intrinsicPrintlnLong implements PrintStream.println(J)V.
func intrinsicPrintlnLong(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	writeJavaString(strconv.FormatInt(AsLong(args[0]), 10) + "\n")
	return Normal()
}

// intrinsicPrintlnBoolean implements PrintStream.println(Z)V.
func intrinsicPrintlnBoolean(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if AsInt(args[0]) != 0 {
		writeJavaString("true\n")
	} else {
		writeJavaString("false\n")
	}
	return Normal()
}

// intrinsicPrintlnChar implements PrintStream.println(C)V.
func intrinsicPrintlnChar(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	writeJavaString(string(rune(uint16(AsInt(args[0])))) + "\n")
	return Normal()
}

// intrinsicPrintlnByte implements PrintStream.println(B)V.
func intrinsicPrintlnByte(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	writeJavaString(strconv.FormatInt(int64(int8(AsInt(args[0]))), 10) + "\n")
	return Normal()
}

// intrinsicPrintlnShort implements PrintStream.println(S)V.
func intrinsicPrintlnShort(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	writeJavaString(strconv.FormatInt(int64(int16(AsInt(args[0]))), 10) + "\n")
	return Normal()
}

// intrinsicPrintlnString implements PrintStream.println(String)V.
func intrinsicPrintlnString(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[0].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[0].Ref == nil {
		writeJavaString("null\n")
		return Normal()
	}
	units, ok := StringUnits(args[0].Ref)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	writeJavaUnits(units)
	writeJavaString("\n")
	return Normal()
}

// intrinsicStringLength implements String.length().
func intrinsicStringLength(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	units, ok := StringUnits(receiver)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	return OkValue(IntValue(int32(len(units))))
}

// intrinsicThrowableInitCause implements Throwable(Throwable).
func intrinsicThrowableInitCause(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 || args[0].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	SetCause(receiver, args[0].Ref)
	return Normal()
}

// intrinsicThrowableInitMessage implements Throwable(String).
func intrinsicThrowableInitMessage(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 || args[0].Kind != KindRef || args[0].Ref == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if text, ok := StringText(args[0].Ref); ok {
		ThrowableOf(receiver).Message = text
	}
	return Normal()
}

// intrinsicThrowableGetCause implements Throwable.getCause() and the
// InvocationTargetException override.
func intrinsicThrowableGetCause(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	return OkValue(RefValue(CauseOf(receiver)))
}

// intrinsicBoxValueOf implements the boxed primitive valueOf methods.
func intrinsicBoxValueOf(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	desc := primitiveDescOfBox(m.Owner.Name)
	if desc == 0 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(RefValue(t.RT.BoxValue(desc, args[0])))
}

// intrinsicUnbox implements the boxed primitive unboxing methods.
func intrinsicUnbox(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	desc := primitiveDescOfBox(m.Owner.Name)
	d, ok := UnboxValue(receiver)
	if !ok || d.Kind != desc {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	return OkValue(d.Value())
}

// intrinsicIntegerValueOf implements Integer.valueOf(I).
func intrinsicIntegerValueOf(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(RefValue(t.RT.BoxedInteger(AsInt(args[0]))))
}

// intrinsicIntegerIntValue implements Integer.intValue()I.
func intrinsicIntegerIntValue(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return intrinsicUnbox(t, m, receiver, args)
}

// intrinsicFloatIntBitsToFloat implements Float.intBitsToFloat(I)F.
func intrinsicFloatIntBitsToFloat(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(FloatBitsValue(uint32(AsInt(args[0]))))
}

// intrinsicFloatFloatToRawIntBits implements Float.floatToRawIntBits(F)I.
func intrinsicFloatFloatToRawIntBits(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(IntValue(int32(math.Float32bits(AsFloat(args[0])))))
}

// intrinsicFloatFloatToIntBits implements Float.floatToIntBits(F)I.
func intrinsicFloatFloatToIntBits(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(IntValue(FloatToIntBits(AsFloat(args[0]))))
}

// intrinsicDoubleLongBitsToDouble implements Double.longBitsToDouble(J)D.
func intrinsicDoubleLongBitsToDouble(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(DoubleBitsValue(uint64(AsLong(args[0]))))
}

// intrinsicDoubleDoubleToRawLongBits implements Double.doubleToRawLongBits(D)J.
func intrinsicDoubleDoubleToRawLongBits(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(LongValue(int64(math.Float64bits(AsDouble(args[0])))))
}

// intrinsicDoubleDoubleToLongBits implements Double.doubleToLongBits(D)J.
func intrinsicDoubleDoubleToLongBits(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	return OkValue(LongValue(DoubleToLongBits(AsDouble(args[0]))))
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
	if args[0].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[0].Ref == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	name, ok := StringText(args[0].Ref)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	var params []*Class
	if args[1].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[1].Ref != nil {
		if args[1].Ref.Array == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		for _, elem := range args[1].Ref.Array.Elems {
			if elem.Kind != KindRef {
				return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
			}
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

// lookupPublicMethod implements Class.getMethod lookup: public methods by name
// and exact parameter class identity, including inherited and overridden
// declarations. Covariant returns and bridge methods are resolved to the most
// specific declaration.
func (r *Runtime) lookupPublicMethod(c *Class, name string, params []*Class) (*Method, *Object) {
	candidates := r.publicMethodCandidates(c, name, params)
	if len(candidates) == 0 {
		return nil, r.bootstrapThrowable("java/lang/NoSuchMethodException")
	}
	return selectMostSpecific(candidates, c), nil
}

func paramsMatch(m *Method, params []*Class) bool {
	if len(m.ParamClass) != len(params) {
		return false
	}
	for i := range params {
		if m.ParamClass[i] != params[i] {
			return false
		}
	}
	return true
}

func (r *Runtime) publicMethodCandidates(c *Class, name string, params []*Class) []*Method {
	var out []*Method
	seen := map[*Method]bool{}
	add := func(m *Method) {
		if m == nil || seen[m] {
			return
		}
		if m.Name != name || m.Name == "<init>" || m.Name == "<clinit>" {
			return
		}
		if m.Flags&AccPublic == 0 {
			return
		}
		if !paramsMatch(m, params) {
			return
		}
		seen[m] = true
		out = append(out, m)
	}
	if c.IsArray {
		// Array classes inherit Object's public methods.
		for m := range r.mustClass("java/lang/Object").MethodIndex {
			_ = m
		}
	}
	if c.IsPrimitive {
		return out
	}
	for cur := c; cur != nil; cur = cur.Super {
		for _, m := range cur.Methods {
			add(m)
		}
	}
	visited := map[*Class]bool{}
	var visit func(iface *Class)
	visit = func(iface *Class) {
		if iface == nil || visited[iface] {
			return
		}
		visited[iface] = true
		for _, m := range iface.Methods {
			add(m)
		}
		for _, in := range iface.Interfaces {
			visit(in)
		}
	}
	visit(c)
	for cur := c; cur != nil; cur = cur.Super {
		for _, in := range cur.Interfaces {
			visit(in)
		}
	}
	return out
}

// classDepth returns the number of superclass edges from c down to target.
func classDepth(c, target *Class) int {
	depth := 0
	for cur := c; cur != nil; cur = cur.Super {
		if cur == target {
			return depth
		}
		depth++
	}
	return 1 << 20
}

// selectMostSpecific chooses the maximally specific candidate: return-type
// specificity wins, then a declaration deeper in the receiver's hierarchy, then
// a non-bridge method with the same descriptor.
func selectMostSpecific(candidates []*Method, receiver *Class) *Method {
	best := candidates[0]
	for _, cand := range candidates[1:] {
		if moreSpecificMethod(cand, best, receiver) {
			best = cand
		}
	}
	return best
}

func moreSpecificMethod(a, b *Method, receiver *Class) bool {
	// Return type specificity: a's return must be assignable to b's return.
	at, aerr := TypeOfReturn(a)
	bt, berr := TypeOfReturn(b)
	if aerr == nil && berr == nil && at != bt {
		if returnAssignable(at, bt) {
			return true
		}
		if returnAssignable(bt, at) {
			return false
		}
	}
	aBridge := a.Flags&AccBridge != 0
	bBridge := b.Flags&AccBridge != 0
	if aBridge != bBridge {
		return !aBridge
	}
	ad, bd := classDepth(receiver, a.Owner), classDepth(receiver, b.Owner)
	if ad != bd {
		return ad < bd
	}
	return a.Owner.Name < b.Owner.Name
}

// TypeOfReturn parses a method's return type.
func TypeOfReturn(m *Method) (Type, error) { return m.Type.Return, nil }

func returnAssignable(from, to Type) bool {
	if from.Kind == to.Kind {
		if from.Kind == 'I' {
			return true
		}
		if from.Kind == 'L' || from.Kind == 'A' {
			return from.Descriptor == to.Descriptor || to.Descriptor == "Ljava/lang/Object;" || to.Kind == 'A'
		}
		return true
	}
	if from.Kind == 'A' && to.Kind == 'L' && to.ClassName == "java/lang/Object" {
		return true
	}
	return false
}

// intrinsicMethodInvoke implements Method.invoke(Object, Object[])Object.
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
	caller := t.CallerMethod()
	if !t.RT.reflectionAccessible(caller, decl) {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessException"))
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
	if len(rawArgs) != len(decl.Type.Params) {
		if !(len(rawArgs) == 0 && len(decl.Type.Params) == 0) {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
	}
	converted, thrown := convertReflectiveArgs(t, decl, rawArgs)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	isStatic := decl.Flags&AccStatic != 0
	if !isStatic {
		if realReceiver == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
		}
		if !t.RT.IsAssignable(realReceiver.Class, decl.Owner) {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
	}
	var target *Method
	if isStatic {
		if thrown := t.EnsureInitialized(decl.Owner); thrown != nil {
			return ThrowOutcome(thrown)
		}
		target = decl
	} else {
		var thrown *Object
		target, thrown = t.RT.selectVirtual(realReceiver.Class, decl)
		if thrown != nil {
			wrapper := t.RT.NewThrowable(t, "java/lang/reflect/InvocationTargetException")
			SetCause(wrapper, thrown)
			return ThrowOutcome(wrapper)
		}
	}
	if tr := t.trace(); tr != nil {
		tr.reflect(t.ID, decl, target)
	}
	var callReceiver *Object
	if !isStatic {
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
	return OkValue(boxReturnValue(t.RT, decl.Type.Return, out.Value))
}

// reflectionAccessible applies the caller-sensitive access check: the member is
// public, and its declaring class must be accessible from the logical caller.
func (r *Runtime) reflectionAccessible(caller *Method, decl *Method) bool {
	if decl.Flags&AccPublic == 0 {
		return false
	}
	var callerClass *Class
	if caller != nil {
		callerClass = caller.Owner
	}
	if callerClass == nil {
		return true
	}
	return r.classAccessible(callerClass, decl.Owner)
}

// boxReturnValue boxes a primitive method result for reflection.
func boxReturnValue(r *Runtime, ret Type, v Value) Value {
	switch ret.Kind {
	case 'V':
		return RefValue(nil)
	case 'L', 'A':
		return v
	case 'I':
		return RefValue(r.BoxValue(ret.Descriptor[0], v))
	case 'J':
		return RefValue(r.BoxValue('J', v))
	case 'F':
		return RefValue(r.BoxValue('F', v))
	case 'D':
		return RefValue(r.BoxValue('D', v))
	}
	return RefValue(nil)
}

// convertReflectiveArgs validates and converts one Method.invoke argument array.
// Boxed primitive arguments may widen according to JLS 5.3.
func convertReflectiveArgs(t *Thread, decl *Method, raw []Value) ([]Value, *Object) {
	if len(raw) != len(decl.Type.Params) {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
	}
	out := make([]Value, len(raw))
	for i, p := range decl.Type.Params {
		switch p.Kind {
		case 'I', 'J', 'F', 'D':
			if raw[i].Kind != KindRef || raw[i].Ref == nil {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			d, ok := UnboxValue(raw[i].Ref)
			if !ok {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			value, ok := widenPrimitive(d, p.Descriptor[0])
			if !ok {
				return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
			}
			out[i] = value
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

// widenPrimitive converts a boxed primitive value to the target primitive
// descriptor, applying JLS 5.3 widening primitive conversion.
func widenPrimitive(d boxData, target byte) (Value, bool) {
	from := d.Kind
	if from == target {
		return d.Value(), true
	}
	intLike := func(v int32) bool {
		switch target {
		case 'I', 'J', 'F', 'D':
			return true
		}
		return false
	}
	switch from {
	case 'Z':
		if target == 'Z' {
			return d.Value(), true
		}
	case 'B', 'S', 'C', 'I':
		v := AsInt(d.Value())
		if !intLike(v) {
			return Value{}, false
		}
		switch target {
		case 'I':
			return IntValue(v), true
		case 'J':
			return LongValue(int64(v)), true
		case 'F':
			return FloatValue(float32(v)), true
		case 'D':
			return DoubleValue(float64(v)), true
		}
	case 'J':
		v := AsLong(d.Value())
		switch target {
		case 'J':
			return LongValue(v), true
		case 'F':
			return FloatValue(float32(v)), true
		case 'D':
			return DoubleValue(float64(v)), true
		}
	case 'F':
		v := AsFloat(d.Value())
		switch target {
		case 'F':
			return FloatValue(v), true
		case 'D':
			return DoubleValue(float64(v)), true
		}
	case 'D':
		if target == 'D' {
			return d.Value(), true
		}
	}
	return Value{}, false
}

func (r *Runtime) StringOfGo(s string) *Object { return r.NewStringGo(s) }

func (r *Runtime) debugMethod(m *Method) string {
	return fmt.Sprintf("%s.%s%s", m.Owner.Name, m.Name, m.Descriptor)
}
