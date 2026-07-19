package classfile

// attrLocation records where an attribute appears in the classfile structure.
// It is used to reject attributes that appear in the wrong location (e.g.
// BootstrapMethods on a field).
type attrLocation int

const (
	locClass  attrLocation = iota // ClassFile attributes
	locField                      // field_info attributes
	locMethod                     // method_info attributes
	locCode                       // Code attributes
)

// AttributeInfo is the base for all attributes (JVMS §4.7). Most attributes are
// skipped after reading their name/length; only the ones catty cares about have
// a concrete readInfo that extracts fields. Unknown attributes are ignored.
type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

// BootstrapMethodsAttr holds the BootstrapMethods attribute (JVMS §4.7.23).
// Each entry contains a bootstrap method handle and its static arguments.
// The full validation is completed in Phase 4.
type BootstrapMethodsAttr struct {
	entries []BootstrapMethodEntry
}

// BootstrapMethodEntry is one entry in the bootstrap method table.
type BootstrapMethodEntry struct {
	MethodRef uint16   // CONSTANT_MethodHandle index
	Arguments []uint16 // constant-pool argument indexes
}

func (a *BootstrapMethodsAttr) readInfo(_ *ClassReader) {}

// NumEntries returns the number of bootstrap method entries.
func (a *BootstrapMethodsAttr) NumEntries() int { return len(a.entries) }

// Entry returns the i-th bootstrap method entry. It returns a copy of the
// argument slice so callers cannot mutate parser-owned metadata.
func (a *BootstrapMethodsAttr) Entry(i int) BootstrapMethodEntry {
	e := a.entries[i]
	if len(e.Arguments) > 0 {
		args := make([]uint16, len(e.Arguments))
		copy(args, e.Arguments)
		e.Arguments = args
	}
	return e
}

// findBootstrapMethods extracts the BootstrapMethods attribute from the class
// attributes. It returns nil if no BootstrapMethods attribute is present.
// Duplicate BootstrapMethods attributes are rejected.
func findBootstrapMethods(attrs []AttributeInfo) *BootstrapMethodsAttr {
	var found *BootstrapMethodsAttr
	for _, a := range attrs {
		if bm, ok := a.(*BootstrapMethodsAttr); ok {
			if found != nil {
				panic(parsePanic{err: &FormatError{
					Op:  "BootstrapMethods",
					Msg: "duplicate BootstrapMethods attribute",
				}})
			}
			found = bm
		}
	}
	return found
}

// readAttributes reads the attribute array that closes several class file
// structures (class, field, method, Code). cp resolves attribute names; loc
// constrains which attributes are allowed at this position.
func readAttributes(reader *ClassReader, cp *ConstantPool, loc attrLocation) []AttributeInfo {
	count := reader.ReadUint16()
	attrs := make([]AttributeInfo, count)
	for i := range attrs {
		attrs[i] = readAttribute(reader, cp, loc)
	}
	return attrs
}

func readAttribute(reader *ClassReader, cp *ConstantPool, loc attrLocation) AttributeInfo {
	nameIndex := reader.ReadUint16()
	length := reader.ReadUint32()

	// Safely consume the attribute body. ReadSlice checks bounds so a
	// truncated length cannot escape as a slice-bounds panic.
	info := reader.ReadSlice(length)

	// Resolve the attribute name with strict UTF8 lookup. A malformed
	// name_index (0, OOB, second slot, non-UTF8 tag) is a *FormatError, not a
	// silently-ignored unknown attribute.
	attrName := cp.utf8Checked(nameIndex)

	// Only decode attributes we use; everything else is a nop. The attribute's
	// body slice (info) is consumed by the typed reader over a fresh ClassReader.
	switch attrName {
	case "Code":
		if loc != locMethod {
			panicf("Code attribute", "found at non-method location")
		}
		return readCodeAttribute(info, cp)
	case "StackMapTable":
		if loc != locCode {
			panicf("StackMapTable", "found outside Code attribute")
		}
		return readStackMapTable(info)
	case "BootstrapMethods":
		if loc != locClass {
			panicf("BootstrapMethods", "found at non-class location")
		}
		return readBootstrapMethods(info)
	default:
		return &UnparsedAttribute{}
	}
}

// UnparsedAttribute is the sink for attributes catty does not interpret.
// Its body has already been safely consumed by ReadSlice.
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
	a.attributes = readAttributes(r, cp, locCode)
	if r.Len() > 0 {
		panicf("Code", "trailing bytes after Code attribute body (%d remaining)", r.Len())
	}
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

// readBootstrapMethods parses the BootstrapMethods attribute (JVMS §4.7.23).
// Structural parsing only — cross-reference validation happens in
// validateDynamicPool.
func readBootstrapMethods(info []byte) AttributeInfo {
	r := NewClassReader(info)
	count := r.ReadUint16()
	entries := make([]BootstrapMethodEntry, count)
	for i := range entries {
		bsRef := r.ReadUint16() // bootstrap_method_ref
		argc := r.ReadUint16()
		args := make([]uint16, argc)
		for j := range args {
			args[j] = r.ReadUint16()
		}
		entries[i] = BootstrapMethodEntry{
			MethodRef: bsRef,
			Arguments: args,
		}
	}
	if r.Len() > 0 {
		panicf("BootstrapMethods", "trailing bytes after BootstrapMethods attribute body (%d remaining)", r.Len())
	}
	return &BootstrapMethodsAttr{entries: entries}
}
