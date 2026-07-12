package classfile

// AttributeInfo is the base for all attributes (JVMS §4.7). Most attributes are
// skipped after reading their name/length; only the ones catty cares about have
// a concrete readInfo that extracts fields. Unknown attributes are ignored.
type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

// readAttributes reads the attribute array that closes several class file
// structures (class, field, method, Code). cp resolves attribute names.
func readAttributes(reader *ClassReader, cp *ConstantPool) []AttributeInfo {
	count := reader.ReadUint16()
	attrs := make([]AttributeInfo, count)
	for i := range attrs {
		attrs[i] = readAttribute(reader, cp)
	}
	return attrs
}

func readAttribute(reader *ClassReader, cp *ConstantPool) AttributeInfo {
	nameIndex := reader.ReadUint16()
	length := reader.ReadUint32()
	info := reader.data[:length]
	reader.data = reader.data[length:]

	// Only decode attributes we use; everything else is a nop. The attribute's
	// body slice (info) is consumed by the typed reader over a fresh ClassReader.
	switch cp.UTF8(nameIndex) {
	case "Code":
		return readCodeAttribute(info, cp)
	case "StackMapTable":
		return readStackMapTable(info)
	default:
		return &UnparsedAttribute{}
	}
}

// UnparsedAttribute is the sink for attributes catty does not interpret.
type UnparsedAttribute struct{}

func (a *UnparsedAttribute) readInfo(*ClassReader) {}

// CodeAttribute holds the bytecode of a method (JVMS §4.7.3).
type CodeAttribute struct {
	maxStack       uint16
	maxLocals      uint16
	code           []byte
	exceptionTable []*ExceptionTableEntry
	attributes     []AttributeInfo
}

func (a *CodeAttribute) readInfo(_ *ClassReader) {} // parsed by readCodeAttribute

type ExceptionTableEntry struct {
	startPc   uint16
	endPc     uint16
	handlerPc uint16
	catchType uint16
}

func (e *ExceptionTableEntry) StartPc() uint16   { return e.startPc }
func (e *ExceptionTableEntry) EndPc() uint16     { return e.endPc }
func (e *ExceptionTableEntry) HandlerPc() uint16 { return e.handlerPc }
func (e *ExceptionTableEntry) CatchType() uint16 { return e.catchType }

func readCodeAttribute(info []byte, cp *ConstantPool) *CodeAttribute {
	r := NewClassReader(info)
	a := &CodeAttribute{
		maxStack:  r.ReadUint16(),
		maxLocals: r.ReadUint16(),
		code:      r.ReadBytes(),
	}
	// exception table
	exLen := r.ReadUint16()
	a.exceptionTable = make([]*ExceptionTableEntry, exLen)
	for i := range a.exceptionTable {
		a.exceptionTable[i] = &ExceptionTableEntry{
			startPc:   r.ReadUint16(),
			endPc:     r.ReadUint16(),
			handlerPc: r.ReadUint16(),
			catchType: r.ReadUint16(),
		}
	}
	a.attributes = readAttributes(r, cp)
	return a
}

func (a *CodeAttribute) MaxStack() uint16                       { return a.maxStack }
func (a *CodeAttribute) MaxLocals() uint16                      { return a.maxLocals }
func (a *CodeAttribute) Code() []byte                           { return a.code }
func (a *CodeAttribute) ExceptionTable() []*ExceptionTableEntry { return a.exceptionTable }

// StackMapTable returns the parsed StackMapTable attribute, or nil if the method
// has none (e.g. no branches, or pre-Java-6 class files).
func (a *CodeAttribute) StackMapTable() *StackMapTableAttribute {
	for _, attr := range a.attributes {
		if smt, ok := attr.(*StackMapTableAttribute); ok {
			return smt
		}
	}
	return nil
}

// CodeAttributeOf finds and returns the Code attribute among a method's
// attributes, or nil if absent (native/abstract methods have no Code).
func CodeAttributeOf(attrs []AttributeInfo) *CodeAttribute {
	for _, attr := range attrs {
		if code, ok := attr.(*CodeAttribute); ok {
			return code
		}
	}
	return nil
}
