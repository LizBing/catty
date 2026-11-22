// Package rt implements the Catty Java runtime shared by generated Go code and
// the bytecode interpreter. Gate 1 supports one logical Java thread, version
// 65.0 class files, integer/reference values, and the small bootstrap inventory
// documented in IMPLEMENTATION.md section 5.
package rt

import (
	"fmt"
	"unicode/utf16"
)

// Runtime profile identity. Both are checked during image installation before
// any initializer or main method runs.
const (
	ABIRevision = 1
	ProfileID   = "catty-jvm21-gate1"
)

// Kind identifies the payload of a Value. Gate 1 has no category-two values.
type Kind uint8

const (
	KindVoid Kind = iota
	KindInt
	KindRef
)

// Value is the engine-independent Java value ABI. Bits carries an int32 for
// KindInt. Ref carries a GC-visible Java reference for KindRef; nil is Java
// null. Both engines, reflection and host bindings exchange values in this
// representation.
type Value struct {
	Kind Kind    `json:"kind"`
	Bits uint64  `json:"bits,omitempty"`
	Ref  *Object `json:"-"`
}

// IntValue boxes a Java int without losing its two's-complement bits.
func IntValue(v int32) Value { return Value{Kind: KindInt, Bits: uint64(uint32(v))} }

// RefValue boxes a Java reference.
func RefValue(o *Object) Value { return Value{Kind: KindRef, Ref: o} }

// VoidValue is the value of a void expression.
func VoidValue() Value { return Value{Kind: KindVoid} }

// AsInt reads the Java int payload.
func AsInt(v Value) int32 { return int32(uint32(v.Bits)) }

// Outcome is a Java return value or a thrown Java throwable. A non-nil Thrown
// always takes precedence over Value.
type Outcome struct {
	Value  Value
	Thrown *Object
}

// Normal is a void return.
func Normal() Outcome { return Outcome{Value: VoidValue()} }

// OkValue is a normal return carrying v.
func OkValue(v Value) Outcome { return Outcome{Value: v} }

// ThrowOutcome reports a Java throwable through the calling convention.
func ThrowOutcome(o *Object) Outcome { return Outcome{Thrown: o} }

// Object is a Java object header. Fields holds reference and primitive slots
// in a GC-visible slice. Array objects use Array instead of Fields. Payload
// holds runtime-owned data for mirrors, strings, boxed values, methods and
// throwables.
type Object struct {
	Class    *Class
	Fields   []Value
	Array    *ArrayData
	Payload  any
	Identity uint64
	TraceID  int64
}

// ArrayData is the storage of a reference array. Gate 1 does not implement
// primitive arrays; the interpreter only creates reference arrays.
type ArrayData struct {
	Component *Class
	Elems     []Value
}

// Type is a supported value type from a JVM field or method descriptor. Gate 1
// admits int, void (returns only) and reference/array types.
type Type struct {
	// Kind is 'V', 'I', 'L' or 'A'. 'A' carries the full array descriptor in
	// Descriptor (for example "[Ljava/lang/String;").
	Kind       byte
	ClassName  string // internal name for 'L'
	Descriptor string // canonical descriptor for all kinds
}

// MethodType is a parsed method descriptor.
type MethodType struct {
	Params     []Type
	Return     Type
	ParamSlots int
	ArgCount   int
}

// FieldSymbol, MethodSymbol and ClassSymbol identify JVM members in images.
type FieldSymbol struct {
	Owner      string
	Name       string
	Descriptor string
}

type MethodSymbol struct {
	Owner      string
	Name       string
	Descriptor string
}

// FieldDef and ClassDef describe immutable image metadata.
type FieldDef struct {
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
	Flags      uint16 `json:"flags,omitempty"`
}

type ClassDef struct {
	Name           string       `json:"name"`
	Super          string       `json:"super,omitempty"`
	Interfaces     []string     `json:"interfaces,omitempty"`
	Flags          uint16       `json:"flags,omitempty"`
	InstanceFields []FieldDef   `json:"instanceFields,omitempty"`
	StaticFields   []FieldDef   `json:"staticFields,omitempty"`
	Methods        []*MethodDef `json:"methods,omitempty"`
}

// RefKind describes one symbolic constant in a retained method body.
type RefKind uint8

const (
	RefInt RefKind = iota
	RefString
	RefClass
	RefField
	RefMethod
)

// Ref is a symbolic constant. Class names and descriptors use JVM internal
// form. Units preserves UTF-16 for string constants.
type Ref struct {
	Kind       RefKind  `json:"kind"`
	Class      string   `json:"class,omitempty"`
	Owner      string   `json:"owner,omitempty"`
	Name       string   `json:"name,omitempty"`
	Descriptor string   `json:"descriptor,omitempty"`
	Int        int32    `json:"int,omitempty"`
	Units      []uint16 `json:"units,omitempty"`
	CP         int      `json:"cp,omitempty"`
}

// Insn is one retained, verified bytecode instruction. PC is the original
// bytecode offset. Operand interpretation is opcode-specific: A is a local
// index, immediate, or branch target; Ref indexes MethodDef.Refs.
type Insn struct {
	Op  uint8 `json:"op"`
	PC  int32 `json:"pc"`
	A   int32 `json:"a,omitempty"`
	B   int32 `json:"b,omitempty"`
	Ref int32 `json:"ref"`
	Len int32 `json:"len,omitempty"`
}

// Handler is one exception table entry with a half-open bytecode range.
type Handler struct {
	StartPC   int32  `json:"start"`
	EndPC     int32  `json:"end"`
	HandlerPC int32  `json:"handler"`
	CatchType string `json:"catchType,omitempty"`
}

// Entry is the common engine ABI. receiver is separate from args; args are in
// descriptor order. A nil receiver is used by static calls.
type Entry func(t *Thread, method *Method, receiver *Object, args []Value) Outcome

// MethodDef is immutable build output. Engine is one of "aot", "interp",
// "intrinsic" or "abstract". EntryIndex names a compiled Go entry in the
// image's entry table; Entry is attached at load time.
type MethodDef struct {
	Owner      string    `json:"owner"`
	Name       string    `json:"name"`
	Descriptor string    `json:"descriptor"`
	Flags      uint16    `json:"flags,omitempty"`
	Engine     string    `json:"engine"`
	Reason     string    `json:"reason,omitempty"`
	MaxStack   int       `json:"maxStack,omitempty"`
	MaxLocals  int       `json:"maxLocals,omitempty"`
	CodeSize   int       `json:"codeSize,omitempty"`
	Refs       []Ref     `json:"refs,omitempty"`
	Insns      []Insn    `json:"insns,omitempty"`
	Handlers   []Handler `json:"handlers,omitempty"`
	EntryIndex int       `json:"entry"`
	Entry      Entry     `json:"-"`
}

// MethodKey names the application entry point.
type MethodKey struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
}

// Startup is a runtime startup initializer used by the runtime acceptance test
// and by a generated image that needs install-time work.
type Startup func(t *Thread) Outcome

// MainFunc overrides Java main invocation for minimal runtime images.
type MainFunc func(t *Thread, args []Value) Outcome

// Image is the immutable executable image. It is checked for ABI and profile
// compatibility before any initialization.
type Image struct {
	ABI      int         `json:"abi"`
	Profile  string      `json:"profile"`
	Classes  []*ClassDef `json:"classes,omitempty"`
	Main     MethodKey   `json:"main,omitempty"`
	Startups []Startup   `json:"-"`
	MainFunc MainFunc    `json:"-"`
}

// Class is runtime class state. Definitions are immutable after install;
// static storage, mirrors and initialization state are mutable.
type Class struct {
	Name          string
	Super         *Class
	Interfaces    []*Class
	Flags         uint16
	InstanceField []*Field
	StaticField   []*Field
	Methods       []*Method
	MethodIndex   map[string]*Method
	InstanceSize  int
	Statics       []Value
	State         classState
	InitError     *Object
	Mirror        *Object
	IsArray       bool
	IsPrimitive   bool
	Component     *Class
	Descriptor    string
	Loader        *Loader
}

// Field is one resolved field.
type Field struct {
	Owner      *Class
	Name       string
	Descriptor string
	Flags      uint16
	Static     bool
	Offset     int
	Index      int
	Type       Type
}

// Method is one resolved method plus its engine body.
type Method struct {
	Owner      *Class
	Name       string
	Descriptor string
	Flags      uint16
	Engine     string
	Reason     string
	MaxStack   int
	MaxLocals  int
	CodeSize   int
	Refs       []Ref
	Insns      []Insn
	Handlers   []Handler
	Entry      Entry
	Intrinsic  func(t *Thread, method *Method, receiver *Object, args []Value) Outcome
	Type       *MethodType
	ParamClass []*Class
	pcIndex    map[int32]int
}

// Loader is the single application loader for Gate 1. The bootstrap namespace
// is owned by the runtime and never admitted from the application classpath.
type Loader struct {
	Name    string
	Runtime *Runtime
}

type classState uint8

const (
	classUninitialized classState = iota
	classInitializing
	classInitialized
	classErroneous
)

// ParseFieldType parses a supported field descriptor.
func ParseFieldType(desc string) (Type, error) {
	t, n, err := parseType(desc, false)
	if err != nil {
		return Type{}, err
	}
	if n != len(desc) {
		return Type{}, fmt.Errorf("trailing characters in field descriptor %q", desc)
	}
	return t, nil
}

// ParseMethodType parses a supported method descriptor.
func ParseMethodType(desc string) (*MethodType, error) {
	if len(desc) == 0 || desc[0] != '(' {
		return nil, fmt.Errorf("malformed method descriptor %q", desc)
	}
	mt := &MethodType{}
	i := 1
	for i < len(desc) && desc[i] != ')' {
		t, n, err := parseType(desc[i:], false)
		if err != nil {
			return nil, err
		}
		mt.Params = append(mt.Params, t)
		mt.ParamSlots++
		i += n
	}
	if i >= len(desc) || desc[i] != ')' {
		return nil, fmt.Errorf("malformed method descriptor %q", desc)
	}
	i++
	if i >= len(desc) {
		return nil, fmt.Errorf("malformed method descriptor %q", desc)
	}
	ret, n, err := parseType(desc[i:], true)
	if err != nil {
		return nil, err
	}
	mt.Return = ret
	i += n
	if i != len(desc) {
		return nil, fmt.Errorf("trailing characters in method descriptor %q", desc)
	}
	mt.ArgCount = len(mt.Params)
	return mt, nil
}

func parseType(desc string, allowVoid bool) (Type, int, error) {
	if len(desc) == 0 {
		return Type{}, 0, fmt.Errorf("empty type descriptor")
	}
	switch desc[0] {
	case 'V':
		if !allowVoid {
			return Type{}, 0, fmt.Errorf("void is not a value type")
		}
		return Type{Kind: 'V', Descriptor: "V"}, 1, nil
	case 'I':
		return Type{Kind: 'I', Descriptor: "I"}, 1, nil
	case 'Z', 'B', 'C', 'S':
		// All integer-family category-one values share the int operand
		// representation in the Gate 1 profile.
		return Type{Kind: 'I', Descriptor: string(desc[0])}, 1, nil
	case 'L':
		end := 1
		for end < len(desc) && desc[end] != ';' {
			end++
		}
		if end >= len(desc) || end == 1 {
			return Type{}, 0, fmt.Errorf("malformed reference descriptor %q", desc)
		}
		name := desc[1:end]
		return Type{Kind: 'L', ClassName: name, Descriptor: desc[:end+1]}, end + 1, nil
	case '[':
		dims := 0
		for dims < len(desc) && desc[dims] == '[' {
			dims++
		}
		if dims == len(desc) {
			return Type{}, 0, fmt.Errorf("malformed array descriptor %q", desc)
		}
		if desc[dims] == 'V' {
			return Type{}, 0, fmt.Errorf("void array element in %q", desc)
		}
		// Gate 1 only creates reference arrays. Primitive arrays are outside
		// the profile.
		if desc[dims] == 'I' {
			return Type{}, 0, fmt.Errorf("primitive array descriptor %q is outside the Gate 1 profile", desc)
		}
		_, n, err := parseType(desc[dims:], false)
		if err != nil {
			return Type{}, 0, err
		}
		total := dims + n
		return Type{Kind: 'A', Descriptor: desc[:total]}, total, nil
	default:
		return Type{}, 0, fmt.Errorf("unsupported type descriptor %q", desc)
	}
}

// DescriptorOf constructs a field descriptor from a class name as used by
// `new`, `anewarray` and class constants.
func DescriptorOf(name string) string {
	if len(name) > 0 && name[0] == '[' {
		return name
	}
	return "L" + name + ";"
}

// UnitsToString converts Java UTF-16 code units to Go text. Isolated
// surrogates are replaced with U+FFFD; class/member names are restricted to
// ASCII by the profile.
func UnitsToString(units []uint16) string {
	return string(utf16.Decode(units))
}

// StringToUnits converts Go text to Java UTF-16 code units.
func StringToUnits(s string) []uint16 {
	return utf16.Encode([]rune(s))
}

// Flags used by class, field and method definitions.
const (
	AccPublic       = 0x0001
	AccPrivate      = 0x0002
	AccProtected    = 0x0004
	AccStatic       = 0x0008
	AccFinal        = 0x0010
	AccSynchronized = 0x0020
	AccNative       = 0x0100
	AccInterface    = 0x0200
	AccAbstract     = 0x0400
)
