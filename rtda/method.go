package rtda

import "catty/classfile"

// Method is the runtime representation of a class method. Non-native methods
// carry their bytecode; native methods carry a Go implementation invoked in lieu
// of interpretation. A native method may be declared but not yet have a
// registered implementation; HasNativeImplementation() distinguishes the two.
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
	hasNativeImpl  bool // true if nativeFunc is a real implementation (not nil)
}

type exceptionEntry struct {
	startPc   int
	endPc     int
	handlerPc int
	catchType string // internal class name, "" for catch-all (finally via any)
}

// ExceptionEntry is the exported form for the interpreter's exception handler.
type ExceptionEntry = exceptionEntry

func (e *exceptionEntry) StartPc() int      { return e.startPc }
func (e *exceptionEntry) EndPc() int        { return e.endPc }
func (e *exceptionEntry) HandlerPc() int    { return e.handlerPc }
func (e *exceptionEntry) CatchType() string { return e.catchType }

// NativeMethod builds a Method backed by a Go function. Used for the synthetic
// core classes (java.lang.Object/System/...) that catty implements natively.
// argSlotCount counts parameters only; the interpreter adds 1 for `this` on
// instance methods. maxStack of 2 covers any category-2 return value.
// By default the method is public (instance); call SetStatic() for static methods.
func NativeMethod(owner *Class, name, descriptor string, fn func(*Frame)) *Method {
	md := ParseMethodDescriptor(descriptor)
	argSlots := uint(md.ArgSlots())
	return &Method{
		owner:         owner,
		name:          name,
		descriptor:    descriptor,
		accessFlags:   accPublic,
		argSlotCount:  argSlots,
		maxLocals:     argSlots + 1,
		maxStack:      2,
		native:        true,
		nativeFunc:    fn,
		hasNativeImpl: true,
	}
}

// SetStatic marks a NativeMethod as static (adds ACC_STATIC, adjusts maxLocals).
func (m *Method) SetStatic() {
	m.accessFlags |= accStatic
	m.maxLocals = m.argSlotCount
}

// InterpretedMethod builds a Method from parsed class-file data.
func InterpretedMethod(owner *Class, name, descriptor string, access uint16,
	maxStack, maxLocals uint, code []byte, exTable []exceptionEntry) *Method {
	md := ParseMethodDescriptor(descriptor)
	m := &Method{
		owner:          owner,
		name:           name,
		descriptor:     descriptor,
		accessFlags:    access,
		maxStack:       maxStack,
		maxLocals:      maxLocals,
		code:           code,
		exceptionTable: exTable,
		argSlotCount:   uint(md.ArgSlots()),
	}
	if access&accNative != 0 {
		m.native = true
		// Native methods from class files start without an implementation.
		// resolveNativeMethods in the classloader attaches one from the global
		// registry. Methods that remain unresolved throw UnsatisfiedLinkError
		// at call time (see HasNativeImplementation).
		//
		// Native methods in class files have maxLocals=0 (no bytecode),
		// but we need at least enough locals to hold the arguments + `this`.
		minLocals := m.argSlotCount
		if access&accStatic == 0 {
			minLocals++
		}
		if m.maxLocals < minLocals {
			m.maxLocals = minLocals
		}
		if m.maxStack < 2 {
			m.maxStack = 2 // cover cat-2 return value
		}
	}
	return m
}

// SetNativeFunc replaces the native method's Go implementation. Called by the
// classloader when a real Go implementation is available (e.g. System.arraycopy).
func (m *Method) SetNativeFunc(fn func(*Frame)) {
	m.nativeFunc = fn
	m.hasNativeImpl = fn != nil
}

func (m *Method) Owner() *Class                 { return m.owner }
func (m *Method) Name() string                  { return m.name }
func (m *Method) Descriptor() string            { return m.descriptor }
func (m *Method) AccessFlags() uint16           { return m.accessFlags }
func (m *Method) IsStatic() bool                { return m.accessFlags&accStatic != 0 }
func (m *Method) IsNative() bool                { return m.native }
func (m *Method) HasNativeImplementation() bool { return m.hasNativeImpl }
func (m *Method) NativeFunc() func(*Frame)      { return m.nativeFunc }
func (m *Method) ArgSlotCount() uint            { return m.argSlotCount }

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
func (m *Method) MaxStack() uint                   { return m.maxStack }
func (m *Method) MaxLocals() uint                  { return m.maxLocals }
func (m *Method) Code() []byte                     { return m.code }

// StackMap returns the parsed StackMapTable (for type tracking) or nil.
func (m *Method) StackMap() *classfile.StackMapTableAttribute { return m.stackMap }

// SetStackMap attaches the parsed StackMapTable (called during class loading).
func (m *Method) SetStackMap(smt *classfile.StackMapTableAttribute) { m.stackMap = smt }

// JVM access flags (JVMS §4.6 / §4.5), reused by Method, Field, Class.
const (
	accPublic    uint16 = 0x0001
	accPrivate   uint16 = 0x0002
	accProtected uint16 = 0x0004
	accStatic    uint16 = 0x0008
	accFinal     uint16 = 0x0010
	accNative    uint16 = 0x0100
	accInterface uint16 = 0x0200
	accAbstract  uint16 = 0x0400
)
