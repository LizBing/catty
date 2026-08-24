// Package interop is the productized Java↔Go binding surface (ADR-0011):
// embedders declare Go functions with plain signatures; Bind registers
// them as static methods of a synthesized Java class, converting values
// across the boundary.
//
// Bind fails fast at REGISTRATION time on anything unsupported (non-func
// entries, unmapped parameter/result types, malformed shapes). Runtime
// conversion failures surface as java.lang.RuntimeException with an
// "interop:" prefix.
package interop

import (
	"fmt"
	"reflect"
	"strings"

	"catty/internal/kernel"
)

// Spec describes one synthesized bridge class.
type Spec struct {
	// Class is the internal name, e.g. "com/app/GoBridge". The class is
	// born initialized and shadows any classpath twin (registry lookup
	// precedes the resolver).
	Class string

	// Funcs maps Java method names to Go functions. Supported shapes:
	//
	//	func(params...) T            // never throws (T ≠ error)
	//	func(params...) error        // throws RuntimeException(msg) on error
	//	func(params...) (T, error)   // throws on non-nil error
	//
	// Mapped param/result types: string, int32, int64, int (as J),
	// float32, float64, bool, and kernel.Value (opaque passthrough,
	// descriptor Ljava/lang/Object;).
	Funcs map[string]any
}

var (
	stringT = reflect.TypeOf("")
	int32T  = reflect.TypeOf(int32(0))
	int64T  = reflect.TypeOf(int64(0))
	goIntT  = reflect.TypeOf(int(0))
	f32T    = reflect.TypeOf(float32(0))
	f64T    = reflect.TypeOf(float64(0))
	boolT   = reflect.TypeOf(false)
	valueT  = reflect.TypeOf((*kernel.Value)(nil)).Elem()
	errT    = reflect.TypeOf((*error)(nil)).Elem()
)

func descFor(t reflect.Type) (string, error) {
	switch t {
	case stringT:
		return "Ljava/lang/String;", nil
	case int32T:
		return "I", nil
	case int64T, goIntT:
		return "J", nil // Go int maps to J; platform-word caveat documented
	case f32T:
		return "F", nil
	case f64T:
		return "D", nil
	case boolT:
		return "Z", nil
	case valueT:
		return "Ljava/lang/Object;", nil // opaque handle passthrough
	}
	return "", fmt.Errorf("interop: unsupported type %s (supported: string, int32, int64, int, float32, float64, bool, kernel.Value)", t)
}

func fromJava(v kernel.Value, t reflect.Type, argIdx int) (reflect.Value, error) {
	if v == nil {
		if t == valueT {
			// A typed nil interface element: legal Call input, arrives as
			// untyped null inside the Go function.
			return reflect.New(t).Elem(), nil
		}
		return reflect.Value{}, fmt.Errorf("interop: null argument %d for %s", argIdx, t)
	}
	switch t {
	case stringT:
		js, ok := v.(*kernel.JString)
		if !ok || js == nil {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want String", argIdx, v)
		}
		return reflect.ValueOf(js.Go()), nil
	case int32T:
		n, ok := v.(int32)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want int", argIdx, v)
		}
		return reflect.ValueOf(n), nil
	case int64T:
		n, ok := v.(int64)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want long", argIdx, v)
		}
		return reflect.ValueOf(n), nil
	case goIntT:
		n, ok := v.(int64)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want long", argIdx, v)
		}
		return reflect.ValueOf(int(n)), nil
	case f32T:
		n, ok := v.(float32)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want float", argIdx, v)
		}
		return reflect.ValueOf(n), nil
	case f64T:
		n, ok := v.(float64)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want double", argIdx, v)
		}
		return reflect.ValueOf(n), nil
	case boolT:
		n, ok := v.(int32)
		if !ok {
			return reflect.Value{}, fmt.Errorf("interop: argument %d is %T, want boolean", argIdx, v)
		}
		return reflect.ValueOf(n != 0), nil
	case valueT:
		return reflect.ValueOf(v), nil
	}
	return reflect.Value{}, fmt.Errorf("interop: unsupported parameter type %s", t)
}

func toJava(k *kernel.Kernel, r reflect.Value) (kernel.Value, error) {
	switch r.Type() {
	case stringT:
		return k.MakeJStringFromGo(r.String()), nil
	case int32T:
		return int32(r.Int()), nil
	case int64T:
		return r.Int(), nil
	case goIntT:
		return int64(r.Int()), nil
	case f32T:
		return float32(r.Float()), nil
	case f64T:
		return r.Float(), nil
	case boolT:
		if r.Bool() {
			return int32(1), nil
		}
		return int32(0), nil
	case valueT:
		if r.IsNil() {
			return nil, nil
		}
		return r.Interface().(kernel.Value), nil
	}
	if r.Type().Implements(valueT) {
		if r.IsNil() {
			return nil, nil
		}
		return r.Interface().(kernel.Value), nil
	}
	return nil, fmt.Errorf("interop: unsupported result type %s", r.Type())
}

// buildNative compiles one Go function into a kernel.NativeFunc wrapper.
// All reflection happens here (registration); the returned closure only
// converts values and calls.
func buildNative(k *kernel.Kernel, name string, fn any) (string, kernel.NativeFunc, error) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return "", nil, fmt.Errorf("interop: %q is %T, want func", name, fn)
	}
	ft := fv.Type()
	if ft.IsVariadic() {
		return "", nil, fmt.Errorf("interop: %q is variadic — not supported", name)
	}

	var argsDesc strings.Builder
	for i := 0; i < ft.NumIn(); i++ {
		d, err := descFor(ft.In(i))
		if err != nil {
			return "", nil, fmt.Errorf("interop: %q param %d: %w", name, i, err)
		}
		argsDesc.WriteString(d)
	}
	hasErr := ft.NumOut() > 0 && ft.Out(ft.NumOut()-1) == errT
	valOuts := ft.NumOut()
	if hasErr {
		valOuts--
	}
	if valOuts > 1 {
		return "", nil, fmt.Errorf("interop: %q returns %d values + error — at most one", name, valOuts)
	}
	retDesc := "V"
	if valOuts == 1 {
		d, err := descFor(ft.Out(0))
		if err != nil {
			return "", nil, fmt.Errorf("interop: %q result: %w", name, err)
		}
		retDesc = d
	}
	desc := "(" + argsDesc.String() + ")" + retDesc

	native := func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
		if len(args) != ft.NumIn() {
			return nil, ctx.Throw("java/lang/RuntimeException",
				fmt.Sprintf("interop: %s expects %d arguments, got %d", name, ft.NumIn(), len(args)))
		}
		ins := make([]reflect.Value, ft.NumIn())
		for i := range ins {
			rv, err := fromJava(args[i], ft.In(i), i)
			if err != nil {
				return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
			}
			ins[i] = rv
		}
		outs := fv.Call(ins)
		if hasErr {
			if e := outs[len(outs)-1]; !e.IsNil() {
				return nil, ctx.Throw("java/lang/RuntimeException", e.Interface().(error).Error())
			}
		}
		if valOuts == 1 {
			v, err := toJava(k, outs[0])
			if err != nil {
				return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
			}
			return v, nil
		}
		return nil, nil
	}
	return desc, native, nil
}

// Bind registers spec.Funcs as public static methods of a synthesized
// class. Registration is fail-fast: the first unsupported shape aborts
// with an error naming the function. Call before first dispatch (same
// window as gen.Install).
func Bind(k *kernel.Kernel, spec Spec) error {
	if spec.Class == "" {
		return fmt.Errorf("interop: Spec.Class is empty")
	}
	if len(spec.Funcs) == 0 {
		return fmt.Errorf("interop: Spec.Funcs is empty")
	}
	methods := make([]kernel.MethodDef, 0, len(spec.Funcs))
	names := make([]string, 0, len(spec.Funcs))
	for n := range spec.Funcs {
		names = append(names, n)
	}
	sortStrings(names)
	for _, n := range names {
		desc, native, err := buildNative(k, n, spec.Funcs[n])
		if err != nil {
			return err
		}
		methods = append(methods, kernel.MethodDef{
			Name: n, Desc: desc,
			Flags:  0x0001 | 0x0008, // public static
			Native: native,
		})
	}
	_, err := k.DefineClass(&kernel.ClassDef{
		Name:    spec.Class,
		Methods: methods,
	})
	if err != nil {
		return fmt.Errorf("interop: define %s: %w", spec.Class, err)
	}
	return nil
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
