// Package classfile parses JVM class files (JVMS chapter 4).
//
// M0 scope: versions 45..52 (JDK 8 baseline is 52). All constant pool
// tags are parsed (so higher-level structures can be skipped safely),
// but invokedynamic-based features are rejected at execution time by
// the VM, not here.
package classfile

import "fmt"

// Magic is the class file magic number (0xCAFEBABE).
const Magic uint32 = 0xCAFEBABE

// MaxSupportedMajor is the highest class file version this parser accepts.
const MaxSupportedMajor uint16 = 52

// ConstTag is a constant pool entry tag (JVMS §4.4).
type ConstTag uint8

// Constant pool tags.
const (
	CUtf8               ConstTag = 1
	CInteger            ConstTag = 3
	CFloat              ConstTag = 4
	CLong               ConstTag = 5
	CDouble             ConstTag = 6
	CClass              ConstTag = 7
	CString             ConstTag = 8
	CFieldref           ConstTag = 9
	CMethodref          ConstTag = 10
	CInterfaceMethodref ConstTag = 11
	CNameAndType        ConstTag = 12
	CMethodHandle       ConstTag = 15
	CMethodType         ConstTag = 16
	CDynamic            ConstTag = 17
	CInvokeDynamic      ConstTag = 18
	CModule             ConstTag = 19
	CPackage            ConstTag = 20
)

func (t ConstTag) String() string {
	switch t {
	case CUtf8:
		return "Utf8"
	case CInteger:
		return "Integer"
	case CFloat:
		return "Float"
	case CLong:
		return "Long"
	case CDouble:
		return "Double"
	case CClass:
		return "Class"
	case CString:
		return "String"
	case CFieldref:
		return "Fieldref"
	case CMethodref:
		return "Methodref"
	case CInterfaceMethodref:
		return "InterfaceMethodref"
	case CNameAndType:
		return "NameAndType"
	case CMethodHandle:
		return "MethodHandle"
	case CMethodType:
		return "MethodType"
	case CDynamic:
		return "Dynamic"
	case CInvokeDynamic:
		return "InvokeDynamic"
	case CModule:
		return "Module"
	case CPackage:
		return "Package"
	default:
		return fmt.Sprintf("tag(%d)", uint8(t))
	}
}

// Access flags shared by classes and members (JVMS §4.1, §4.5, §4.6).
const (
	AccPublic       uint16 = 0x0001
	AccPrivate      uint16 = 0x0002
	AccProtected    uint16 = 0x0004
	AccStatic       uint16 = 0x0008
	AccFinal        uint16 = 0x0010
	AccSuper        uint16 = 0x0020 // class only
	AccSynchronized uint16 = 0x0020 // method only
	AccVolatile     uint16 = 0x0040 // field only
	AccBridge       uint16 = 0x0040 // method only
	AccTransient    uint16 = 0x0080 // field only
	AccVarargs      uint16 = 0x0080 // method only
	AccNative       uint16 = 0x0100
	AccInterface    uint16 = 0x0200
	AccAbstract     uint16 = 0x0400
	AccSynthetic    uint16 = 0x1000
)

// ConstEntry is a decoded constant pool entry. Which fields are meaningful
// depends on Tag (mirrors JVMS §4.4 struct layout; Idx1/Idx2 are the raw u2s).
type ConstEntry struct {
	Tag       ConstTag
	Str       string  // CUtf8
	IntVal    int32   // CInteger
	LongVal   int64   // CLong
	FloatVal  float32 // CFloat
	DoubleVal float64 // CDouble
	Idx1      uint16  // see tag docs below
	Idx2      uint16  // see tag docs below

	// Idx1/Idx2 meaning per tag:
	//   CClass, CString, CMethodType, CModule, CPackage: Idx1 = utf8/class idx
	//   CFieldref/CMethodref/CInterfaceMethodref: Idx1=classIdx, Idx2=nameAndTypeIdx
	//   CNameAndType: Idx1=nameIdx(utf8), Idx2=descIdx(utf8)
	//   CMethodHandle: Idx1=reference_kind(u1), Idx2=refIdx
	//   CInvokeDynamic/CDynamic: Idx1=bootstrapAttrIdx, Idx2=nameAndTypeIdx
}

// Attribute is a raw attribute (name resolved, data unparsed).
type Attribute struct {
	Name string
	Data []byte
}

// ExceptionHandler is one entry of a Code exception table (JVMS §4.7.3).
// CatchType == 0 means catch-all (finally).
type ExceptionHandler struct {
	StartPc   uint16
	EndPc     uint16 // exclusive
	HandlerPc uint16
	CatchType uint16 // constant pool index of CClass; 0 = any
}

// Code is the parsed Code attribute of a method.
type Code struct {
	MaxStack  uint16
	MaxLocals uint16
	Code      []byte
	Handlers  []ExceptionHandler
	StackMaps []StackMapFrame // decoded StackMapTable (may be nil)

	// LineNumbers is the LineNumberTable (JVMS §4.7.12), ascending by
	// StartPc. Nil when the attribute is absent.
	LineNumbers []LineNum
}

// LineNum maps a bytecode range start to a source line.
type LineNum struct {
	StartPc uint16
	Line    uint16
}

// ParsedAnnotation is one decoded runtime annotation.
type ParsedAnnotation struct {
	TypeDesc string            // e.g. "Lcom/google/gson/annotations/SerializedName;"
	Elements map[string]ElementValue
}

// ElementValue is one element_value in an annotation (JVMS §4.7.16.1).
type ElementValue struct {
	Tag       byte   // 'B','C','D','I','J','S','Z','s','e','c','@','['
	ConstIdx  uint16 // for primitive/String tags: const pool idx
	EnumType  string // for 'e': type descriptor
	EnumName  string // for 'e': constant name
	ClassDesc string // for 'c': return descriptor
	Array     []ElementValue // for '['
}

// FieldInfo is a parsed field.
type FieldInfo struct {
	AccessFlags uint16
	Name        string
	Desc        string
	Annotations []ParsedAnnotation
	Attributes  []Attribute
}

// ConstantValue returns the ConstantValue attribute's pool index, or 0.
func (f *FieldInfo) ConstantValue() uint16 {
	for _, a := range f.Attributes {
		if a.Name == "ConstantValue" && len(a.Data) == 2 {
			return uint16(a.Data[0])<<8 | uint16(a.Data[1])
		}
	}
	return 0
}

// MethodInfo is a parsed method.
type MethodInfo struct {
	AccessFlags uint16
	Name        string
	Annotations []ParsedAnnotation
	Desc        string
	Code        *Code // nil for abstract/native methods
	Attributes  []Attribute
}

// ClassFile is a parsed class file.
type ClassFile struct {
	MinorVersion uint16
	MajorVersion uint16
	Constants    []ConstEntry // 1-based; index 0 is the zero value
	AccessFlags  uint16
	ThisClass    string
	SuperClass   string // "" for java/lang/Object itself
	Annotations  []ParsedAnnotation
	Signature    string // generic signature
	SourceFile   string // SourceFile attribute, "" when absent
	Interfaces   []string
	Fields       []FieldInfo
	Methods      []MethodInfo
	Attributes   []Attribute
}

// Entry returns the constant pool entry at 1-based index i.
func (f *ClassFile) Entry(i uint16) (*ConstEntry, error) {
	if i == 0 || int(i) >= len(f.Constants) {
		return nil, fmt.Errorf("classfile %s: constant pool index %d out of range [1,%d)", f.ThisClass, i, len(f.Constants))
	}
	return &f.Constants[i], nil
}

// UTF8 returns the UTF8 string at pool index i.
func (f *ClassFile) UTF8(i uint16) (string, error) {
	e, err := f.Entry(i)
	if err != nil {
		return "", err
	}
	if e.Tag != CUtf8 {
		return "", fmt.Errorf("classfile %s: pool[%d] is %s, want Utf8", f.ThisClass, i, e.Tag)
	}
	return e.Str, nil
}

// ClassName resolves a CClass entry at pool index i to its name.
func (f *ClassFile) ClassName(i uint16) (string, error) {
	e, err := f.Entry(i)
	if err != nil {
		return "", err
	}
	if e.Tag != CClass {
		return "", fmt.Errorf("classfile %s: pool[%d] is %s, want Class", f.ThisClass, i, e.Tag)
	}
	return f.UTF8(e.Idx1)
}

// NameAndType resolves a CNameAndType entry to (name, descriptor).
func (f *ClassFile) NameAndType(i uint16) (name, desc string, err error) {
	e, err := f.Entry(i)
	if err != nil {
		return "", "", err
	}
	if e.Tag != CNameAndType {
		return "", "", fmt.Errorf("classfile %s: pool[%d] is %s, want NameAndType", f.ThisClass, i, e.Tag)
	}
	if name, err = f.UTF8(e.Idx1); err != nil {
		return "", "", err
	}
	if desc, err = f.UTF8(e.Idx2); err != nil {
		return "", "", err
	}
	return name, desc, nil
}

// Ref resolves a CFieldref / CMethodref / CInterfaceMethodref into
// (class name, member name, member descriptor).
func (f *ClassFile) Ref(i uint16) (cls, name, desc string, tag ConstTag, err error) {
	e, err := f.Entry(i)
	if err != nil {
		return "", "", "", 0, err
	}
	switch e.Tag {
	case CFieldref, CMethodref, CInterfaceMethodref:
	default:
		return "", "", "", 0, fmt.Errorf("classfile %s: pool[%d] is %s, want a member ref", f.ThisClass, i, e.Tag)
	}
	if cls, err = f.ClassName(e.Idx1); err != nil {
		return "", "", "", 0, err
	}
	if name, desc, err = f.NameAndType(e.Idx2); err != nil {
		return "", "", "", 0, err
	}
	return cls, name, desc, e.Tag, nil
}

// FindMethod returns the method with matching name+desc, or nil.
func (f *ClassFile) FindMethod(name, desc string) *MethodInfo {
	for i := range f.Methods {
		m := &f.Methods[i]
		if m.Name == name && m.Desc == desc {
			return m
		}
	}
	return nil
}

// FindField returns the field with matching name+desc, or nil.
func (f *ClassFile) FindField(name, desc string) *FieldInfo {
	for i := range f.Fields {
		fd := &f.Fields[i]
		if fd.Name == name && fd.Desc == desc {
			return fd
		}
	}
	return nil
}
