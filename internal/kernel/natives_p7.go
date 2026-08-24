package kernel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ---- Wrapper type natives (P-0007) ----

func natWrapperInit(primDesc string) NativeFunc {
	return func(ctx *CallContext, recv Value, args []Value) (Value, error) {
		recv.(*Instance).Fields[0] = args[0]
		return nil, nil
	}
}

func natWrapperValue(primDesc string) NativeFunc {
	return func(ctx *CallContext, recv Value, args []Value) (Value, error) {
		return recv.(*Instance).Fields[0], nil
	}
}

func natWrapperHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := recv.(*Instance).Fields[0]
	switch x := v.(type) {
	case int32:
		return x, nil
	case int64:
		return int32(int64(x) ^ (int64(x) >> 32)), nil
	case bool:
		if x {
			return int32(1231), nil
		}
		return int32(1237), nil
	default:
		return int32(0), nil
	}
}

func natWrapperEquals(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if len(args) == 0 || args[0] == nil {
		return boolV(false), nil
	}
	a := recv.(*Instance).Fields[0]
	other, ok := args[0].(*Instance)
	if !ok {
		return boolV(false), nil
	}
	b := other.Fields[0]
	return boolV(a == b), nil
}

func natWrapperToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := recv.(*Instance).Fields[0]
	var s string
	switch x := v.(type) {
	case int32:
		s = strconv.FormatInt(int64(x), 10)
	case int64:
		s = strconv.FormatInt(x, 10)
	case bool:
		s = strconv.FormatBool(x)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return ctx.K.InternGo(s), nil
}

func natWrapperValueOf(clsName string) NativeFunc {
	return func(ctx *CallContext, recv Value, args []Value) (Value, error) {
		c, ok := ctx.K.ClassByName(clsName)
		if !ok {
			return nil, fmt.Errorf("class not found: %s", clsName)
		}
		obj, err := ctx.K.NewInstance(c)
		if err != nil {
			return nil, err
		}
		obj.Fields[0] = args[0]
		return obj, nil
	}
}

func natWrapperParse(clsName string) NativeFunc {
	return func(ctx *CallContext, recv Value, args []Value) (Value, error) {
		if len(args) == 0 || args[0] == nil {
			return nil, ctx.Throw("java/lang/NumberFormatException", "null")
		}
		js, _ := AsJString(args[0])
		s := js.String()
		switch clsName {
		case "java/lang/Long":
			v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				return nil, ctx.Throw("java/lang/NumberFormatException", s)
			}
			return v, nil
		case "java/lang/Short":
			v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 16)
			if err != nil {
				return nil, ctx.Throw("java/lang/NumberFormatException", s)
			}
			return int32(v), nil
		case "java/lang/Byte":
			v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 8)
			if err != nil {
				return nil, ctx.Throw("java/lang/NumberFormatException", s)
			}
			return int32(v), nil
		}
		return nil, ctx.Throw("java/lang/NumberFormatException", s)
	}
}

func natBooleanParse(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if len(args) == 0 || args[0] == nil {
		return boolV(false), nil
	}
	js, _ := AsJString(args[0])
	return boolV(strings.EqualFold(js.String(), "true")), nil
}

// ---- HashMap natives ----

// jkey is the content/identity key for map/set backing stores. Strings key
// by their shared UTF-8 form (no per-op allocation); boxed numerics by
// class+value (Java equals semantics); everything else by identity.
type jkey struct {
	kind byte   // 'S' string-content, 'N' numeric, 'R' identity, 'X' fallback, 'Z' null
	s    string // shared UTF-8 for strings / rendered form for fallback
	cls  string // declaring class for numerics
	num  int64  // numeric payload
	p    any    // identity pointer
}

func hashKey(v Value) jkey {
	switch x := v.(type) {
	case nil:
		return jkey{kind: 'Z'}
	case *JString:
		return jkey{kind: 'S', s: x.Go()}
	case *Instance:
		f := x.Fields[0]
		switch n := f.(type) {
		case int32:
			return jkey{kind: 'N', cls: x.Class.Name, num: int64(n)}
		case int64:
			return jkey{kind: 'N', cls: x.Class.Name, num: n}
		default:
			return jkey{kind: 'R', p: x}
		}
	case *ArrayObj:
		return jkey{kind: 'R', p: x}
	default:
		return jkey{kind: 'X', s: fmt.Sprintf("%v", v)}
	}
}

type mapBuf struct {
	data map[jkey]Value
	keys map[jkey]Value
}

func hashMapOf(recv Value) *mapBuf { return recv.(*Instance).Payload.(*mapBuf) }

func natHashMapInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &mapBuf{
		data: make(map[jkey]Value),
		keys: make(map[jkey]Value),
	}
	return nil, nil
}

func natHashMapGet(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	val, ok := m.data[hashKey(args[0])]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func natHashMapPut(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	key := hashKey(args[0])
	old, existed := m.data[key]
	m.data[key] = args[1]
	m.keys[key] = args[0]
	if !existed {
		return nil, nil
	}
	return old, nil
}

func natHashMapRemove(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	key := hashKey(args[0])
	old, existed := m.data[key]
	delete(m.data, key)
	delete(m.keys, key)
	if !existed {
		return nil, nil
	}
	return old, nil
}

func natHashMapContainsKey(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	_, ok := m.data[hashKey(args[0])]
	return boolV(ok), nil
}

func natHashMapSize(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	return int32(len(m.data)), nil
}

func natHashMapIsEmpty(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	return boolV(len(m.data) == 0), nil
}

func natHashMapClear(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m := hashMapOf(recv)
	m.data = make(map[jkey]Value)
	m.keys = make(map[jkey]Value)
	return nil, nil
}

// ---- HashSet natives ----

type setBuf struct {
	data map[jkey]bool
	keys map[jkey]Value
}

func setOf(recv Value) *setBuf { return recv.(*Instance).Payload.(*setBuf) }

func natHashSetInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &setBuf{
		data: make(map[jkey]bool),
		keys: make(map[jkey]Value),
	}
	return nil, nil
}

func natHashSetAdd(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s := setOf(recv)
	key := hashKey(args[0])
	_, existed := s.data[key]
	s.data[key] = true
	s.keys[key] = args[0]
	return boolV(!existed), nil
}

func natHashSetContains(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s := setOf(recv)
	_, ok := s.data[hashKey(args[0])]
	return boolV(ok), nil
}

func natHashSetRemove(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s := setOf(recv)
	key := hashKey(args[0])
	_, existed := s.data[key]
	delete(s.data, key)
	delete(s.keys, key)
	return boolV(existed), nil
}

func natHashSetSize(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s := setOf(recv)
	return int32(len(s.data)), nil
}

func natHashSetIsEmpty(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s := setOf(recv)
	return boolV(len(s.data) == 0), nil
}

// ---- String wide-method natives ----

func javaStr(v Value) string {
	js, err := AsJString(v)
	if err != nil {
		return ""
	}
	return js.String()
}

func natStringSubstring(ctx *CallContext, recv Value, args []Value) (Value, error) {
	r := []rune(javaStr(recv))
	begin := argI(args, 0)
	if int(begin) < 0 || int(begin) > len(r) {
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException",
			fmt.Sprintf("begin %d, length %d", begin, len(r)))
	}
	return ctx.K.InternGo(string(r[begin:])), nil
}

func natStringSubstring2(ctx *CallContext, recv Value, args []Value) (Value, error) {
	r := []rune(javaStr(recv))
	begin := argI(args, 0)
	end := argI(args, 1)
	if int(begin) < 0 || int(end) > len(r) || begin > end {
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException",
			fmt.Sprintf("begin %d, end %d, length %d", begin, end, len(r)))
	}
	return ctx.K.InternGo(string(r[begin:end])), nil
}

func natStringTrim(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(strings.TrimSpace(javaStr(recv))), nil
}

func natStringToLowerCase(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(strings.ToLower(javaStr(recv))), nil
}

func natStringToUpperCase(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(strings.ToUpper(javaStr(recv))), nil
}

func natStringContains(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(strings.Contains(javaStr(recv), javaStr(args[0]))), nil
}

func natStringIsEmpty(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(len(javaStr(recv)) == 0), nil
}

func natStringSplit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	pattern := javaStr(args[0])
	s := javaStr(recv)

	var parts []string
	if re, err := regexp.Compile(pattern); err == nil {
		parts = re.Split(s, -1) // Java-regex subset via Go regexp
	} else {
		parts = strings.Split(s, pattern) // literal fallback
	}

	arr := &ArrayObj{Elems: make([]Value, len(parts))}
	for i, p := range parts {
		arr.Elems[i] = ctx.K.InternGo(p)
	}
	return arr, nil
}

func natStringReplace(ctx *CallContext, recv Value, args []Value) (Value, error) {
	oldC := rune(argI(args, 0))
	newC := rune(argI(args, 1))
	result := strings.Map(func(r rune) rune {
		if r == oldC {
			return newC
		}
		return r
	}, javaStr(recv))
	return ctx.K.InternGo(result), nil
}

func natStringStartsWith(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(strings.HasPrefix(javaStr(recv), javaStr(args[0]))), nil
}

func natStringEndsWith(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(strings.HasSuffix(javaStr(recv), javaStr(args[0]))), nil
}
