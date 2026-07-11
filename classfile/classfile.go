package classfile

// ClassFile is the parsed representation of a .class file (JVMS §4.1).
// Parsing is total and upfront: every method's Code attribute is decoded so the
// interpreter can execute straight from this struct without re-parsing.
type ClassFile struct {
	// magic is 0xCAFEBABE; checked in Parse.
	minorVersion uint16
	majorVersion uint16
	constantPool *ConstantPool
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []*MemberInfo
	methods      []*MemberInfo
	attributes   []AttributeInfo
}

// Parse decodes a full class file from raw bytes.
func Parse(data []byte) (*ClassFile, error) {
	r := NewClassReader(data)
	cf := &ClassFile{}
	cf.read(r)
	return cf, nil
}

func (cf *ClassFile) read(r *ClassReader) {
	cf.readAndCheckMagic(r)
	cf.readAndCheckVersions(r)
	cf.constantPool = readConstantPool(r)
	cf.accessFlags = r.ReadUint16()
	cf.thisClass = r.ReadUint16()
	cf.superClass = r.ReadUint16()
	cf.interfaces = r.ReadUint16s()
	cf.fields = readMembers(r, cf.constantPool)
	cf.methods = readMembers(r, cf.constantPool)
	cf.attributes = readAttributes(r, cf.constantPool)
}

func (cf *ClassFile) readAndCheckMagic(r *ClassReader) {
	magic := r.ReadUint32()
	if magic != 0xCAFEBABE {
		panic("catty: bad class file magic")
	}
}

func (cf *ClassFile) readAndCheckVersions(r *ClassReader) {
	cf.minorVersion = r.ReadUint16()
	cf.majorVersion = r.ReadUint16()
	// JDK 25 is major 69. We do not enforce a ceiling here: catty supports the
	// bytecode subset it supports regardless of class version. javac may emit
	// features beyond our subset; we fail loudly if encountered at run time.
}

// --- Accessors used by the class loader and interpreter ---

func (cf *ClassFile) MajorVersion() uint16   { return cf.majorVersion }
func (cf *ClassFile) ConstantPool() *ConstantPool { return cf.constantPool }
func (cf *ClassFile) AccessFlags() uint16    { return cf.accessFlags }

// ClassName returns the internal name of the class defined by this file
// (e.g. "java/lang/Object").
func (cf *ClassFile) ClassName() string {
	return cf.constantPool.ClassName(cf.thisClass)
}

// SuperClassName returns the internal name of the superclass, or "" for
// java.lang.Object (whose super_class index is 0).
func (cf *ClassFile) SuperClassName() string {
	if cf.superClass == 0 {
		return ""
	}
	return cf.constantPool.ClassName(cf.superClass)
}

func (cf *ClassFile) InterfaceNames() []string {
	names := make([]string, len(cf.interfaces))
	for i, idx := range cf.interfaces {
		names[i] = cf.constantPool.ClassName(idx)
	}
	return names
}

func (cf *ClassFile) Fields() []*MemberInfo  { return cf.fields }
func (cf *ClassFile) Methods() []*MemberInfo { return cf.methods }
