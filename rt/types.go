// Package rt implements the Catty Java runtime shared by generated Go code and
// the bytecode interpreter. It supports one logical Java thread, version 65.0
// class files, the primitive and reference value set published by the Gate 2
// profile, and the bootstrap inventory documented in GATE2.md.
package rt

import (
	"fmt"
	"math"
	"unicode/utf16"
)

// Kind identifies the payload of a Value.
type Kind uint8

const (
	KindVoid Kind = iota
	KindInt
	KindLong
	KindFloat
	KindDouble
	KindRef
)

// Value is the engine-independent Java value ABI. Bits carries the raw payload
// of a primitive value: the low 32 bits for int, all 64 bits for long, the
// IEEE 754 bit pattern in the low 32 bits for float, and the IEEE 754 bit
// pattern for double. Ref carries a GC-visible Java reference for KindRef; nil
// is Java null. Both engines, reflection and host bindings exchange values in
// this representation.
type Value struct {
	Kind Kind    `json:"kind"`
	Bits uint64  `json:"bits,omitempty"`
	Ref  *Object `json:"-"`
}

// IntValue boxes a Java int without losing its two's-complement bits.
func IntValue(v int32) Value { return Value{Kind: KindInt, Bits: uint64(uint32(v))} }

// LongValue boxes a Java long.
func LongValue(v int64) Value { return Value{Kind: KindLong, Bits: uint64(v)} }

// LongBitsValue boxes a Java long from its raw 64-bit pattern. It accepts
// constants that are not representable as a signed int64 in generated code.
func LongBitsValue(bits uint64) Value { return Value{Kind: KindLong, Bits: bits} }

// FloatValue boxes a Java float, preserving its IEEE 754 bits.
func FloatValue(v float32) Value { return Value{Kind: KindFloat, Bits: uint64(math.Float32bits(v))} }

// FloatBitsValue boxes a Java float from its raw IEEE 754 bits.
func FloatBitsValue(bits uint32) Value { return Value{Kind: KindFloat, Bits: uint64(bits)} }

// DoubleValue boxes a Java double, preserving its IEEE 754 bits.
func DoubleValue(v float64) Value { return Value{Kind: KindDouble, Bits: math.Float64bits(v)} }

// DoubleBitsValue boxes a Java double from its raw IEEE 754 bits.
func DoubleBitsValue(bits uint64) Value { return Value{Kind: KindDouble, Bits: bits} }

// RefValue boxes a Java reference.
func RefValue(o *Object) Value { return Value{Kind: KindRef, Ref: o} }

// VoidValue is the value of a void expression.
func VoidValue() Value { return Value{Kind: KindVoid} }

// AsInt reads the Java int payload.
func AsInt(v Value) int32 { return int32(uint32(v.Bits)) }

// AsLong reads the Java long payload.
func AsLong(v Value) int64 { return int64(v.Bits) }

// AsFloatBits reads the raw IEEE 754 bits of a Java float.
func AsFloatBits(v Value) uint32 { return uint32(v.Bits) }

// AsFloat reads the Java float payload.
func AsFloat(v Value) float32 { return math.Float32frombits(uint32(v.Bits)) }

// AsDoubleBits reads the raw IEEE 754 bits of a Java double.
func AsDoubleBits(v Value) uint64 { return v.Bits }

// AsDouble reads the Java double payload.
func AsDouble(v Value) float64 { return math.Float64frombits(v.Bits) }

// Width is the number of Java local/operand slots a value occupies.
func (v Value) Width() int {
	if v.Kind == KindLong || v.Kind == KindDouble {
		return 2
	}
	return 1
}

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

// Object is a Java object header. Fields holds instance slots in a
// GC-visible slice (category-two values occupy one slot). Array objects use
// Array instead. Payload holds runtime-owned data for mirrors, strings, boxed
// values, methods and throwables.
type Object struct {
	Class    *Class
	Fields   []Value
	Array    *ArrayData
	Payload  any
	Identity uint64
	TraceID  int64
}

// ArrayData is the storage of a Java array. Reference arrays retain
// GC-visible Value slots; primitive arrays use their typed backing slice.
type ArrayData struct {
	Component *Class
	Kind      byte // 'L' reference array; otherwise the primitive descriptor letter
	Elems     []Value
	Bytes     []int8
	Chars     []uint16
	Shorts    []int16
	Ints      []int32
	Longs     []int64
	Floats    []float32
	Doubles   []float64
}

// Len returns the Java array length.
func (a *ArrayData) Len() int {
	if a == nil {
		return 0
	}
	switch a.Kind {
	case 'L':
		return len(a.Elems)
	case 'Z', 'B':
		return len(a.Bytes)
	case 'C':
		return len(a.Chars)
	case 'S':
		return len(a.Shorts)
	case 'I':
		return len(a.Ints)
	case 'J':
		return len(a.Longs)
	case 'F':
		return len(a.Floats)
	case 'D':
		return len(a.Doubles)
	}
	return 0
}

// Type is a supported value type from a JVM field or method descriptor. Kind
// is 'V', 'I' (the int family: boolean, byte, char, short, int), 'J', 'F',
// 'D', 'L' or 'A'. 'A' carries the full array descriptor in Descriptor.
type Type struct {
	Kind       byte
	ClassName  string // internal name for 'L'
	Descriptor string // canonical descriptor for all kinds
}

// Slots is the number of Java local variable or operand stack slots.
func (t Type) Slots() int {
	if t.Kind == 'J' || t.Kind == 'D' {
		return 2
	}
	if t.Kind == 'V' {
		return 0
	}
	return 1
}

// IsReference reports whether the type is a reference type.
func (t Type) IsReference() bool { return t.Kind == 'L' || t.Kind == 'A' }

// MethodType is a parsed method descriptor. ParamSlot holds the local slot
// offset of each parameter (excluding any receiver); ParamSlots is the total
// slot width of all parameters; ArgCount is the logical argument count.
type MethodType struct {
	Params     []Type
	Return     Type
	ParamSlots int
	ArgCount   int
	ParamSlot  []int
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

// Constant is a ConstantValue attribute payload. Kind is one of
// "int", "long", "float", "double" or "string".
type Constant struct {
	Kind  string   `json:"kind"`
	Bits  uint64   `json:"bits,omitempty"`
	Units []uint16 `json:"units,omitempty"`
}

// FieldDef and ClassDef describe immutable image metadata.
type FieldDef struct {
	Name       string    `json:"name"`
	Descriptor string    `json:"descriptor"`
	Flags      uint16    `json:"flags,omitempty"`
	Constant   *Constant `json:"constant,omitempty"`
}

type ClassDef struct {
	Name           string       `json:"name"`
	Super          string       `json:"super,omitempty"`
	Interfaces     []string     `json:"interfaces,omitempty"`
	Flags          uint16       `json:"flags,omitempty"`
	NestHost       string       `json:"nestHost,omitempty"`
	NestMembers    []string     `json:"nestMembers,omitempty"`
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
	RefLong
	RefFloat
	RefDouble
)

// Ref is a symbolic constant. Class names and descriptors use JVM internal
// form. Bits carries long, float and double payload bits; Units preserves
// UTF-16 for string constants.
type Ref struct {
	Kind       RefKind  `json:"kind"`
	Class      string   `json:"class,omitempty"`
	Owner      string   `json:"owner,omitempty"`
	Name       string   `json:"name,omitempty"`
	Descriptor string   `json:"descriptor,omitempty"`
	Int        int32    `json:"int,omitempty"`
	Bits       uint64   `json:"bits,omitempty"`
	Units      []uint16 `json:"units,omitempty"`
	CP         int      `json:"cp,omitempty"`
}

// Insn is one retained, verified bytecode instruction. PC is the original
// bytecode offset. Operand interpretation is opcode-specific: A is a local
// index, immediate, switch default target or branch target; B is an iinc
// increment, invokeinterface slot count or multianewarray dimension count; Ref
// indexes MethodDef.Refs. Keys and Targets carry switch arms.
type Insn struct {
	Op      uint8   `json:"op"`
	PC      int32   `json:"pc"`
	A       int32   `json:"a,omitempty"`
	B       int32   `json:"b,omitempty"`
	Ref     int32   `json:"ref"`
	Len     int32   `json:"len,omitempty"`
	Keys    []int32 `json:"keys,omitempty"`
	Targets []int32 `json:"targets,omitempty"`
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

// Image is the immutable executable image. It is checked for ABI, profile and
// opcode inventory compatibility before any initialization.
type Image struct {
	ABI          int           `json:"abi"`
	Profile      string        `json:"profile"`
	Opcodes      []string      `json:"opcodes,omitempty"`
	Capabilities []string      `json:"capabilities,omitempty"`
	Plugins      bool          `json:"plugins,omitempty"`
	Limits       *PluginLimits `json:"pluginLimits,omitempty"`
	Classes      []*ClassDef   `json:"classes,omitempty"`
	Main         MethodKey     `json:"main,omitempty"`
	Startups     []Startup     `json:"-"`
	MainFunc     MainFunc      `json:"-"`
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
	InstanceKinds []string
	Statics       []Value
	State         classState
	InitError     *Object
	Mirror        *Object
	IsArray       bool
	IsPrimitive   bool
	Component     *Class
	Descriptor    string
	NestHost      string
	NestMembers   []string
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
	Constant   *Constant
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

// Loader is one class-defining loader. The bootstrap loader owns the fixed
// bootstrap inventory; the application loader owns the image classes; plugin
// loaders own admitted JAR definitions. Lookup is parent-first and a parent
// never searches its children.
type Loader struct {
	ID      int
	Name    string
	Runtime *Runtime
	Parent  *Loader
	Classes map[string]*Class
	Plugin  bool
	Digest  string
	Mirror  *Object
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
		mt.ParamSlot = append(mt.ParamSlot, mt.ParamSlots)
		mt.ParamSlots += t.Slots()
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
		// The int family shares the int operand representation.
		return Type{Kind: 'I', Descriptor: string(desc[0])}, 1, nil
	case 'J':
		return Type{Kind: 'J', Descriptor: "J"}, 1, nil
	case 'F':
		return Type{Kind: 'F', Descriptor: "F"}, 1, nil
	case 'D':
		return Type{Kind: 'D', Descriptor: "D"}, 1, nil
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
	AccBridge       = 0x0040
	AccVolatile     = 0x0040
	AccNative       = 0x0100
	AccInterface    = 0x0200
	AccAbstract     = 0x0400
	AccSynthetic    = 0x1000
)

// PrimitiveClassInfo is the bootstrap metadata for one primitive class.
type PrimitiveClassInfo struct {
	Name       string
	Descriptor string
}

// PrimitiveClasses lists the canonical primitive classes in deterministic
// order.
var PrimitiveClasses = []PrimitiveClassInfo{
	{"boolean", "Z"},
	{"byte", "B"},
	{"char", "C"},
	{"short", "S"},
	{"int", "I"},
	{"long", "J"},
	{"float", "F"},
	{"double", "D"},
	{"void", "V"},
}

// PrimitiveName returns the Java name for a primitive descriptor letter.
func PrimitiveName(desc byte) string {
	switch desc {
	case 'Z':
		return "boolean"
	case 'B':
		return "byte"
	case 'C':
		return "char"
	case 'S':
		return "short"
	case 'I':
		return "int"
	case 'J':
		return "long"
	case 'F':
		return "float"
	case 'D':
		return "double"
	case 'V':
		return "void"
	}
	return ""
}

// PrimitiveDescriptor returns the descriptor letter for a Java primitive
// class name.
func PrimitiveDescriptor(name string) (byte, bool) {
	switch name {
	case "boolean":
		return 'Z', true
	case "byte":
		return 'B', true
	case "char":
		return 'C', true
	case "short":
		return 'S', true
	case "int":
		return 'I', true
	case "long":
		return 'J', true
	case "float":
		return 'F', true
	case "double":
		return 'D', true
	case "void":
		return 'V', true
	}
	return 0, false
}
