package rtda

import "catty/classfile"

// Method is the runtime representation of a class method. Non-native methods
// carry their bytecode; native methods carry a Go implementation invoked in lieu
// of interpretation.
type Method struct {
	owner          *Class
	name           string
	descriptor     string
	accessFlags    uint16
	maxStack       uint
	maxLocals      uint
	code           []byte
	exceptionTable []exceptionEntry
	stackMap       *classfile.StackMapTableAttribute
	argSlotCount   uint
	native         bool
	nativeFunc     func(*Frame)
}

type exceptionEntry struct {
	startPc   int
	endPc     int
	handlerPc int
	catchType string // internal class name, "" for catch-all (finally via any)
}

// NativeMethod builds a Method backed by a Go function. Used for the synthetic
// core classes (java.lang.Object/System/...) that catty implements natively.
// argSlotCount counts parameters only; the interpreter adds 1 for `this` on
// instance methods. maxStack of 2 covers any category-2 return value.
func NativeMethod(owner *Class, name, descriptor string, fn func(*Frame)) *Method {
	md := ParseMethodDescriptor(descriptor)
	argSlots := uint(md.ArgSlots())
	return &Method{
		owner:        owner,
		name:         name,
		descriptor:   descriptor,
		argSlotCount: argSlots,
		maxLocals:    argSlots + 1,
		maxStack:     2,
		native:       true,
		nativeFunc:   fn,
	}
}

// InterpretedMethod builds a Method from parsed class-file data.
func InterpretedMethod(owner *Class, name, descriptor string, access uint16,
	maxStack, maxLocals uint, code []byte, exTable []exceptionEntry) *Method {
	md := ParseMethodDescriptor(descriptor)
	m := &Method{
		owner:         owner,
		name:          name,
		descriptor:    descriptor,
		accessFlags:   access,
		maxStack:      maxStack,
		maxLocals:     maxLocals,
		code:          code,
		exceptionTable: exTable,
		argSlotCount:  uint(md.ArgSlots()),
	}
	if access&accNative != 0 {
		m.native = true // no nativeFunc yet; resolved against the registry at load
	}
	return m
}

func (m *Method) Owner() *Class    { return m.owner }
func (m *Method) Name() string     { return m.name }
func (m *Method) Descriptor() string { return m.descriptor }
func (m *Method) AccessFlags() uint16 { return m.accessFlags }
func (m *Method) IsStatic() bool      { return m.accessFlags&accStatic != 0 }
func (m *Method) IsNative() bool            { return m.native }
func (m *Method) NativeFunc() func(*Frame)  { return m.nativeFunc }
func (m *Method) ArgSlotCount() uint        { return m.argSlotCount }

// ReturnType returns the descriptor of the return type (the part after ')').
func (m *Method) ReturnType() string {
	for i := 0; i < len(m.descriptor); i++ {
		if m.descriptor[i] == ')' {
			return m.descriptor[i+1:]
		}
	}
	return "V"
}

// ExceptionTable exposes parsed handlers for try/catch (used post-MVP).
func (m *Method) ExceptionTable() []exceptionEntry { return m.exceptionTable }
func (m *Method) MaxStack() uint   { return m.maxStack }
func (m *Method) MaxLocals() uint  { return m.maxLocals }
func (m *Method) Code() []byte     { return m.code }

// StackMap returns the parsed StackMapTable (for type tracking) or nil.
func (m *Method) StackMap() *classfile.StackMapTableAttribute { return m.stackMap }

// SetStackMap attaches the parsed StackMapTable (called during class loading).
func (m *Method) SetStackMap(smt *classfile.StackMapTableAttribute) { m.stackMap = smt }

// JVM access flags (JVMS §4.6 / §4.5), reused by Method, Field, Class.
const (
	accPublic     uint16 = 0x0001
	accPrivate    uint16 = 0x0002
	accProtected  uint16 = 0x0004
	accStatic     uint16 = 0x0008
	accFinal      uint16 = 0x0010
	accNative     uint16 = 0x0100
	accInterface  uint16 = 0x0200
	accAbstract   uint16 = 0x0400
)
