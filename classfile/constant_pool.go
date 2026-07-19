package classfile

import (
	"fmt"
	"math"
	"strings"
)

// ConstantInfo is the tagged-union root for constant pool entries (JVMS §4.4).
// All entries are parsed inline in readConstantInfo, so the interface only
// carries Tag(), used to detect Long/Double's two-slot occupancy.
type ConstantInfo interface {
	Tag() uint8
}

// ConstantPool is indexed from 1 (index 0 is unused per spec). Entries pointing
// at Long/Double consume the following slot too.
type ConstantPool struct {
	infos []ConstantInfo
}

func (cp *ConstantPool) Size() int { return len(cp.infos) }

// lookup returns the ConstantInfo at index, or an error. It rejects index 0,
// out-of-range, nil entries (unusable second slot of long/double), and is the
// single safe entry point for all constant-pool reads. Public typed accessors
// use lookup and propagate errors; parse-internal code wraps lookup errors into
// parsePanic via get/Tag/utf8Checked.
func (cp *ConstantPool) lookup(index uint16) (ConstantInfo, error) {
	if index == 0 {
		return nil, fmt.Errorf("index 0 is not a valid constant pool index")
	}
	if int(index) >= len(cp.infos) {
		return nil, fmt.Errorf("index %d out of bounds (pool size %d)", index, len(cp.infos)-1)
	}
	info := cp.infos[index]
	if info == nil {
		return nil, fmt.Errorf("index %d refers to an unusable slot (second half of long/double)", index)
	}
	return info, nil
}

// readConstantPool reads count-1 entries (count includes the unused slot 0).
func readConstantPool(reader *ClassReader) *ConstantPool {
	cpCount := reader.ReadUint16()
	infos := make([]ConstantInfo, cpCount)
	cp := &ConstantPool{infos: infos}

	for i := 1; i < int(cpCount); i++ { // index 0 unused, start at 1
		info := readConstantInfo(reader, cp)
		infos[i] = info
		if info.Tag() == ConstantLong || info.Tag() == ConstantDouble {
			i++ // long/double take two slots
		}
	}
	return cp
}

func readConstantInfo(reader *ClassReader, cp *ConstantPool) ConstantInfo {
	tag := reader.ReadUint8() // JVMS §4.4: the cp_info tag is a u1, not u2.
	switch tag {
	case ConstantInteger:
		return &ConstantIntegerInfo{val: int32(reader.ReadUint32())}
	case ConstantFloat:
		return &ConstantFloatInfo{val: math.Float32frombits(reader.ReadUint32())}
	case ConstantLong:
		return &ConstantLongInfo{val: int64(reader.ReadUint64())}
	case ConstantDouble:
		return &ConstantDoubleInfo{val: math.Float64frombits(reader.ReadUint64())}
	case ConstantUtf8:
		length := reader.ReadUint16()
		raw := make([]byte, length)
		rb := reader.ReadSlice(uint32(length))
		copy(raw, rb)
		return &ConstantUtf8Info{str: decodeMUTF8(raw), raw: raw}
	case ConstantString:
		return &ConstantStringInfo{cp: cp, strIndex: reader.ReadUint16()}
	case ConstantClass:
		return &ConstantClassInfo{cp: cp, nameIndex: reader.ReadUint16()}
	case ConstantFieldref, ConstantMethodref, ConstantInterfaceMethodref:
		return &ConstantMemberRefInfo{
			cp:         cp,
			tag:        tag,
			classIndex: reader.ReadUint16(),
			natIndex:   reader.ReadUint16(),
		}
	case ConstantNameAndType:
		return &ConstantNameAndTypeInfo{cp: cp,
			nameIndex: reader.ReadUint16(), descIndex: reader.ReadUint16()}
	case ConstantMethodType:
		return &ConstantMethodTypeInfo{descIndex: reader.ReadUint16()}
	case ConstantMethodHandle:
		return &ConstantMethodHandleInfo{
			refKind:  reader.ReadUint8(),
			refIndex: reader.ReadUint16(),
		}
	case ConstantDynamic:
		return &ConstantDynamicInfo{
			bootstrapMethodAttrIndex: reader.ReadUint16(),
			natIndex:                 reader.ReadUint16(),
		}
	case ConstantInvokeDynamic:
		return &ConstantInvokeDynamicInfo{
			bootstrapMethodAttrIndex: reader.ReadUint16(),
			natIndex:                 reader.ReadUint16(),
		}
	default:
		panic(parsePanic{err: &FormatError{
			Op:  "constant pool",
			Msg: fmt.Sprintf("unknown or unsupported constant tag %d", int(tag)),
		}})
	}
}

// Constant pool tags (JVMS §4.4, Table 4.4-B).
const (
	ConstantUtf8               uint8 = 1
	ConstantInteger            uint8 = 3
	ConstantFloat              uint8 = 4
	ConstantLong               uint8 = 5
	ConstantDouble             uint8 = 6
	ConstantClass              uint8 = 7
	ConstantString             uint8 = 8
	ConstantFieldref           uint8 = 9
	ConstantMethodref          uint8 = 10
	ConstantInterfaceMethodref uint8 = 11
	ConstantNameAndType        uint8 = 12
	ConstantMethodHandle       uint8 = 15
	ConstantMethodType         uint8 = 16
	ConstantDynamic            uint8 = 17
	ConstantInvokeDynamic      uint8 = 18
)

// MethodHandle reference kinds (JVMS §5.4.3.5).
const (
	RefGetField         uint8 = 1
	RefGetStatic        uint8 = 2
	RefPutField         uint8 = 3
	RefPutStatic        uint8 = 4
	RefInvokeVirtual    uint8 = 5
	RefInvokeStatic     uint8 = 6
	RefInvokeSpecial    uint8 = 7
	RefNewInvokeSpecial uint8 = 8
	RefInvokeInterface  uint8 = 9
)

// refKindTarget maps a MethodHandle reference kind to the constant-pool tag it
// must reference. Kinds 6 and 7 have dual targets — InterfaceMethodref is
// permitted only for classfile version ≥ 52.
var refKindTarget = [10]struct{ primary, alt uint8 }{
	RefGetField:         {ConstantFieldref, 0},
	RefGetStatic:        {ConstantFieldref, 0},
	RefPutField:         {ConstantFieldref, 0},
	RefPutStatic:        {ConstantFieldref, 0},
	RefInvokeVirtual:    {ConstantMethodref, 0},
	RefInvokeStatic:     {ConstantMethodref, ConstantInterfaceMethodref},
	RefInvokeSpecial:    {ConstantMethodref, ConstantInterfaceMethodref},
	RefNewInvokeSpecial: {ConstantMethodref, 0},
	RefInvokeInterface:  {ConstantInterfaceMethodref, 0},
}

// --- Concrete constant info types ---

type ConstantUtf8Info struct {
	str string // Go string decoded via decodeMUTF8 (for names/descriptors)
	raw []byte // raw MUTF-8 bytes (for lossless UTF-16 String constant decoding)
}

func (c *ConstantUtf8Info) Tag() uint8  { return ConstantUtf8 }
func (c *ConstantUtf8Info) Str() string { return c.str }

// UTF16 returns the lossless UTF-16 code units decoded from the raw MUTF-8
// bytes. Used only for CONSTANT_String entries.
func (c *ConstantUtf8Info) UTF16() []uint16 { return decodeMUTF8ToUTF16(c.raw) }

type ConstantIntegerInfo struct{ val int32 }

func (c *ConstantIntegerInfo) Tag() uint8 { return ConstantInteger }
func (c *ConstantIntegerInfo) Val() int32 { return c.val }

type ConstantFloatInfo struct{ val float32 }

func (c *ConstantFloatInfo) Tag() uint8   { return ConstantFloat }
func (c *ConstantFloatInfo) Val() float32 { return c.val }

type ConstantLongInfo struct{ val int64 }

func (c *ConstantLongInfo) Tag() uint8 { return ConstantLong }
func (c *ConstantLongInfo) Val() int64 { return c.val }

type ConstantDoubleInfo struct{ val float64 }

func (c *ConstantDoubleInfo) Tag() uint8   { return ConstantDouble }
func (c *ConstantDoubleInfo) Val() float64 { return c.val }

type ConstantStringInfo struct {
	cp       *ConstantPool
	strIndex uint16
}

func (c *ConstantStringInfo) Tag() uint8  { return ConstantString }
func (c *ConstantStringInfo) Str() string { return c.cp.UTF8(c.strIndex) }

// UTF16 returns the CONSTANT_String literal as lossless UTF-16 code units.
func (c *ConstantStringInfo) UTF16() []uint16 {
	if info, ok := c.cp.get(c.strIndex).(*ConstantUtf8Info); ok {
		return info.UTF16()
	}
	return nil
}

type ConstantClassInfo struct {
	cp        *ConstantPool
	nameIndex uint16
}

func (c *ConstantClassInfo) Tag() uint8   { return ConstantClass }
func (c *ConstantClassInfo) Name() string { return c.cp.UTF8(c.nameIndex) }

type ConstantNameAndTypeInfo struct {
	cp        *ConstantPool
	nameIndex uint16
	descIndex uint16
}

func (c *ConstantNameAndTypeInfo) Tag() uint8         { return ConstantNameAndType }
func (c *ConstantNameAndTypeInfo) Name() string       { return c.cp.UTF8(c.nameIndex) }
func (c *ConstantNameAndTypeInfo) Descriptor() string { return c.cp.UTF8(c.descIndex) }

// ConstantMemberRefInfo backs Fieldref / Methodref / InterfaceMethodref; they
// share an identical encoding and differ only in tag, so one type serves all.
type ConstantMemberRefInfo struct {
	cp         *ConstantPool
	tag        uint8
	classIndex uint16
	natIndex   uint16
}

func (c *ConstantMemberRefInfo) Tag() uint8 { return c.tag }

// These legacy accessors serve the interpreter/classloader through the pool's
// MemberRef() helper. They use get() (which panics on bad indices) and return
// zero values on type mismatches.
func (c *ConstantMemberRefInfo) ClassName() string  { return c.cp.ClassName(c.classIndex) }
func (c *ConstantMemberRefInfo) Name() string       { return c.cp.NameAndTypeName(c.natIndex) }
func (c *ConstantMemberRefInfo) Descriptor() string { return c.cp.NameAndTypeDescriptor(c.natIndex) }

type ConstantMethodTypeInfo struct{ descIndex uint16 }

func (c *ConstantMethodTypeInfo) Tag() uint8 { return ConstantMethodType }

type ConstantMethodHandleInfo struct {
	refKind  uint8
	refIndex uint16
}

func (c *ConstantMethodHandleInfo) Tag() uint8 { return ConstantMethodHandle }

type ConstantDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	natIndex                 uint16
}

func (c *ConstantDynamicInfo) Tag() uint8 { return ConstantDynamic }

type ConstantInvokeDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	natIndex                 uint16
}

func (c *ConstantInvokeDynamicInfo) Tag() uint8 { return ConstantInvokeDynamic }

// --- Parse-internal accessors (panic on malformed index) ---
//
// These are the only functions that convert lookup errors into parsePanic. They
// are called by legacy zero-value accessors, parse validation, and internal
// readers — all of which run under Parse's recovery boundary.

// get looks up a constant pool entry and panics with parsePanic on any error.
// Legacy accessors use it for their type-assertion pattern; Parse converts the
// panic into *FormatError.
func (cp *ConstantPool) get(index uint16) ConstantInfo {
	info, err := cp.lookup(index)
	if err != nil {
		panicf("constant pool", "%v", err)
	}
	return info
}

// Tag returns the JVMS constant tag at index. It panics with parsePanic on
// malformed index so parse validation code can use it without error threading.
func (cp *ConstantPool) Tag(index uint16) uint8 {
	return cp.get(index).Tag()
}

// utf8Checked returns the UTF-8 string at index or panics with parsePanic.
// Used by attribute name resolution: a malformed name_index is a format error,
// not a silently-ignored unknown attribute.
func (cp *ConstantPool) utf8Checked(index uint16) string {
	info := cp.get(index)
	u, ok := info.(*ConstantUtf8Info)
	if !ok {
		panicf("constant pool", "index %d: expected CONSTANT_Utf8 for attribute name, got tag %d", index, info.Tag())
	}
	return u.str
}

// --- Legacy zero-value accessors on the pool ---
//
// These accessors use get(), which panics on malformed indices. They return zero
// values for type mismatches and serve existing callers in the interpreter,
// classloader, and AOT pipeline.

// UTF8 returns the string at a CONSTANT_Utf8 index.
func (cp *ConstantPool) UTF8(index uint16) string {
	if info, ok := cp.get(index).(*ConstantUtf8Info); ok {
		return info.str
	}
	return ""
}

// ClassName returns the internal class name at a CONSTANT_Class index.
func (cp *ConstantPool) ClassName(index uint16) string {
	if info, ok := cp.get(index).(*ConstantClassInfo); ok {
		return info.Name()
	}
	return ""
}

func (cp *ConstantPool) NameAndTypeName(index uint16) string {
	if info, ok := cp.get(index).(*ConstantNameAndTypeInfo); ok {
		return info.Name()
	}
	return ""
}

func (cp *ConstantPool) NameAndTypeDescriptor(index uint16) string {
	if info, ok := cp.get(index).(*ConstantNameAndTypeInfo); ok {
		return info.Descriptor()
	}
	return ""
}

// MemberRef returns the (class, name, descriptor) triple at a Field/Method ref.
func (cp *ConstantPool) MemberRef(index uint16) (className, name, descriptor string) {
	if info, ok := cp.get(index).(*ConstantMemberRefInfo); ok {
		return info.ClassName(), info.Name(), info.Descriptor()
	}
	return "", "", ""
}

// String returns the literal at a CONSTANT_String index (for ldc).
func (cp *ConstantPool) String(index uint16) string {
	if info, ok := cp.get(index).(*ConstantStringInfo); ok {
		return info.Str()
	}
	return ""
}

// UTF16 returns the lossless UTF-16 code units for the CONSTANT_String at index.
func (cp *ConstantPool) UTF16(index uint16) []uint16 {
	if info, ok := cp.get(index).(*ConstantStringInfo); ok {
		return info.UTF16()
	}
	return nil
}

// Integer returns the int constant at a CONSTANT_Integer index (for ldc).
func (cp *ConstantPool) Integer(index uint16) int32 {
	if info, ok := cp.get(index).(*ConstantIntegerInfo); ok {
		return info.Val()
	}
	return 0
}

// Float returns the float constant at a CONSTANT_Float index (for ldc).
func (cp *ConstantPool) Float(index uint16) float32 {
	if info, ok := cp.get(index).(*ConstantFloatInfo); ok {
		return info.Val()
	}
	return 0
}

// Long returns the long constant at a CONSTANT_Long index (for ldc2_w).
func (cp *ConstantPool) Long(index uint16) int64 {
	if info, ok := cp.get(index).(*ConstantLongInfo); ok {
		return info.Val()
	}
	return 0
}

// Double returns the double constant at a CONSTANT_Double index (for ldc2_w).
func (cp *ConstantPool) Double(index uint16) float64 {
	if info, ok := cp.get(index).(*ConstantDoubleInfo); ok {
		return info.Val()
	}
	return 0
}

// ClassRefNames returns the internal names of all CONSTANT_Class entries in the
// pool — used by the AOT build's reachability traversal to find every class
// referenced by a class's constant pool.
func (cp *ConstantPool) ClassRefNames() []string {
	var names []string
	for i := 1; i < len(cp.infos); i++ {
		info := cp.infos[i]
		if cls, ok := info.(*ConstantClassInfo); ok {
			names = append(names, cls.Name())
		}
		if info != nil && (info.Tag() == ConstantLong || info.Tag() == ConstantDouble) {
			i++ // skip the 2nd slot of long/double
		}
	}
	return names
}

// --- Typed dynamic accessors ---
//
// These accessors use lookup() (not get()), so they return ordinary errors on
// malformed indices instead of panicking. They do not expose the constant-pool
// backing slice.

// MethodTypeDescriptor returns the method descriptor at a CONSTANT_MethodType
// index. Returns an error for wrong tag, invalid index, or non-method descriptor.
func (cp *ConstantPool) MethodTypeDescriptor(index uint16) (string, error) {
	info, err := cp.lookup(index)
	if err != nil {
		return "", fmt.Errorf("constant pool index %d: %w", index, err)
	}
	mt, ok := info.(*ConstantMethodTypeInfo)
	if !ok {
		return "", fmt.Errorf("constant pool index %d: expected CONSTANT_MethodType, got tag %d", index, info.Tag())
	}
	// Verify descriptor_index is Utf8 via checked lookup.
	descInfo, err := cp.lookup(mt.descIndex)
	if err != nil {
		return "", fmt.Errorf("constant pool index %d: MethodType descriptor_index %d: %w", index, mt.descIndex, err)
	}
	descUtf8, ok := descInfo.(*ConstantUtf8Info)
	if !ok {
		return "", fmt.Errorf("constant pool index %d: MethodType descriptor_index %d is not CONSTANT_Utf8 (tag %d)", index, mt.descIndex, descInfo.Tag())
	}
	desc := descUtf8.str
	if err := validateMethodDescriptor(desc); err != nil {
		return "", fmt.Errorf("constant pool index %d: invalid method descriptor %q: %w", index, desc, err)
	}
	return desc, nil
}

// MethodHandleRefKind returns the reference kind (1–9) at a
// CONSTANT_MethodHandle index.
func (cp *ConstantPool) MethodHandleRefKind(index uint16) (uint8, error) {
	info, err := cp.lookup(index)
	if err != nil {
		return 0, fmt.Errorf("constant pool index %d: %w", index, err)
	}
	mh, ok := info.(*ConstantMethodHandleInfo)
	if !ok {
		return 0, fmt.Errorf("constant pool index %d: expected CONSTANT_MethodHandle, got tag %d", index, info.Tag())
	}
	if mh.refKind < 1 || mh.refKind > 9 {
		return 0, fmt.Errorf("constant pool index %d: invalid MethodHandle reference kind %d", index, mh.refKind)
	}
	return mh.refKind, nil
}

// lookupUtf8 is a checked helper: looks up a CP index and asserts it is
// CONSTANT_Utf8, returning the string value.
func (cp *ConstantPool) lookupUtf8(index uint16, label string) (string, error) {
	info, err := cp.lookup(index)
	if err != nil {
		return "", fmt.Errorf("%s %d: %w", label, index, err)
	}
	u, ok := info.(*ConstantUtf8Info)
	if !ok {
		return "", fmt.Errorf("%s %d is not CONSTANT_Utf8 (tag %d)", label, index, info.Tag())
	}
	return u.str, nil
}

// lookupClass is a checked helper: looks up a CP index and asserts it is
// CONSTANT_Class, returning the ClassInfo.
func (cp *ConstantPool) lookupClass(index uint16, label string) (*ConstantClassInfo, error) {
	info, err := cp.lookup(index)
	if err != nil {
		return nil, fmt.Errorf("%s %d: %w", label, index, err)
	}
	cls, ok := info.(*ConstantClassInfo)
	if !ok {
		return nil, fmt.Errorf("%s %d is not CONSTANT_Class (tag %d)", label, index, info.Tag())
	}
	return cls, nil
}

// lookupNameAndType is a checked helper: looks up a CP index and asserts it is
// CONSTANT_NameAndType.
func (cp *ConstantPool) lookupNameAndType(index uint16, label string) (*ConstantNameAndTypeInfo, error) {
	info, err := cp.lookup(index)
	if err != nil {
		return nil, fmt.Errorf("%s %d: %w", label, index, err)
	}
	nat, ok := info.(*ConstantNameAndTypeInfo)
	if !ok {
		return nil, fmt.Errorf("%s %d is not CONSTANT_NameAndType (tag %d)", label, index, info.Tag())
	}
	return nat, nil
}

// validateMethodDescriptorReturnV checks that desc is a valid method descriptor
// whose return descriptor is exactly V.
func validateMethodDescriptorReturnV(desc string) error {
	if err := validateMethodDescriptor(desc); err != nil {
		return err
	}
	// Find the return descriptor: everything after ')'.
	closeIdx := strings.IndexByte(desc, ')')
	if closeIdx == -1 {
		return fmt.Errorf("method descriptor missing ')'")
	}
	returnDesc := desc[closeIdx+1:]
	if returnDesc != "V" {
		return fmt.Errorf("method descriptor return must be V for constructor, got %q", returnDesc)
	}
	return nil
}

// validateMethodHandle performs complete structural validation of a
// CONSTANT_MethodHandle entry per JVMS §5.4.3.5. Both the public typed accessor
// MethodHandleReference and the parse-time validation path call this function.
//
// The validation chain is:
//
//	MethodHandle → reference_index (MemberRef) → class_index (Class)
//	     ↓                                         → name_index (Utf8)
//	     → name_and_type_index (NameAndType) → name_index (Utf8)
//	                                          → descriptor_index (Utf8)
//
// Every link is verified with checked lookup + type assertion before reading
// values. Returns the fully validated (refKind, refTag, className, name, descriptor).
func (cp *ConstantPool) validateMethodHandle(index uint16, majorVersion uint16) (refKind uint8, refTag uint8, className, name, descriptor string, err error) {
	// 1. Safe lookup of the MethodHandle entry.
	info, err := cp.lookup(index)
	if err != nil {
		return 0, 0, "", "", "", fmt.Errorf("constant pool index %d: %w", index, err)
	}
	mh, ok := info.(*ConstantMethodHandleInfo)
	if !ok {
		return 0, 0, "", "", "", fmt.Errorf("constant pool index %d: expected CONSTANT_MethodHandle, got tag %d", index, info.Tag())
	}
	refKind = mh.refKind

	// 2. Validate reference kind 1–9.
	if refKind < 1 || refKind > 9 {
		return refKind, 0, "", "", "", fmt.Errorf("constant pool index %d: invalid MethodHandle reference kind %d", index, refKind)
	}

	// 3. Safe lookup of the reference_index.
	refInfo, err := cp.lookup(mh.refIndex)
	if err != nil {
		return refKind, 0, "", "", "", fmt.Errorf("constant pool index %d: MethodHandle reference_index %d: %w", index, mh.refIndex, err)
	}
	refTag = refInfo.Tag()

	// 4. Validate reference tag against refKindTarget.
	rt := refKindTarget[refKind]
	valid := refTag == rt.primary
	if !valid && rt.alt != 0 && refTag == rt.alt {
		if majorVersion < 52 {
			return refKind, refTag, "", "", "", fmt.Errorf(
				"constant pool index %d: MethodHandle kind %d targets InterfaceMethodref but classfile version %d < 52",
				index, refKind, majorVersion)
		}
		valid = true
	}
	if !valid {
		return refKind, refTag, "", "", "", fmt.Errorf(
			"constant pool index %d: MethodHandle kind %d references tag %d, want %d",
			index, refKind, refTag, rt.primary)
	}

	// 5. reference_index must be a ConstantMemberRefInfo.
	mr, ok := refInfo.(*ConstantMemberRefInfo)
	if !ok {
		return refKind, refTag, "", "", "", fmt.Errorf(
			"constant pool index %d: MethodHandle reference_index %d is not a member ref (tag %d)",
			index, mh.refIndex, refTag)
	}

	// 6. class_index → CONSTANT_Class → name_index → CONSTANT_Utf8.
	cls, err := cp.lookupClass(mr.classIndex, fmt.Sprintf("constant pool index %d: class_index", index))
	if err != nil {
		return refKind, refTag, "", "", "", err
	}
	className, err = cp.lookupUtf8(cls.nameIndex, fmt.Sprintf("constant pool index %d: class name_index", index))
	if err != nil {
		return refKind, refTag, "", "", "", err
	}
	if err := validateInternalName(className); err != nil {
		return refKind, refTag, "", "", "", fmt.Errorf("constant pool index %d: invalid class name %q: %w", index, className, err)
	}

	// 7. name_and_type_index → CONSTANT_NameAndType → name_index/desc_index → CONSTANT_Utf8.
	nat, err := cp.lookupNameAndType(mr.natIndex, fmt.Sprintf("constant pool index %d: name_and_type_index", index))
	if err != nil {
		return refKind, refTag, "", "", "", err
	}
	name, err = cp.lookupUtf8(nat.nameIndex, fmt.Sprintf("constant pool index %d: name_index", index))
	if err != nil {
		return refKind, refTag, "", "", "", err
	}
	descriptor, err = cp.lookupUtf8(nat.descIndex, fmt.Sprintf("constant pool index %d: descriptor_index", index))
	if err != nil {
		return refKind, refTag, "", "", "", err
	}

	// 8. Per-kind name and descriptor rules (JVMS §5.4.3.5).
	switch refKind {
	case RefGetField, RefGetStatic, RefPutField, RefPutStatic: // kinds 1–4: field
		// Field names are unqualified names; <init> and <clinit> are valid.
		if err := validateUnqualifiedName(name); err != nil {
			return refKind, refTag, "", "", "", fmt.Errorf(
				"constant pool index %d: invalid field name %q: %w", index, name, err)
		}
		if !validFieldDescriptor(descriptor) {
			return refKind, refTag, "", "", "", fmt.Errorf(
				"constant pool index %d: field MethodHandle kind %d has non-field descriptor %q", index, refKind, descriptor)
		}

	case RefInvokeVirtual, RefInvokeStatic, RefInvokeSpecial,
		RefNewInvokeSpecial, RefInvokeInterface: // kinds 5–9: method

		// All method kinds require a valid method descriptor.
		if refKind == RefNewInvokeSpecial {
			// Kind 8: must name <init> and return descriptor must be V.
			if name != "<init>" {
				return refKind, refTag, "", "", "", fmt.Errorf(
					"constant pool index %d: MethodHandle kind 8 (newInvokeSpecial) must name <init>, got %q", index, name)
			}
			if err := validateMethodDescriptorReturnV(descriptor); err != nil {
				return refKind, refTag, "", "", "", fmt.Errorf(
					"constant pool index %d: MethodHandle kind 8 constructor descriptor %q: %w", index, descriptor, err)
			}
		} else {
			// Kinds 5, 6, 7, 9: must be a valid ordinary method name
			// (no '<' or '>', per JVMS §4.2.2). This rejects <init>,
			// <clinit>, and names like <foo>.
			if err := validateOrdinaryMethodName(name); err != nil {
				return refKind, refTag, "", "", "", fmt.Errorf(
					"constant pool index %d: invalid method name %q for MethodHandle kind %d: %w", index, name, refKind, err)
			}
			if err := validateMethodDescriptor(descriptor); err != nil {
				return refKind, refTag, "", "", "", fmt.Errorf(
					"constant pool index %d: method MethodHandle kind %d has invalid method descriptor %q: %w", index, refKind, descriptor, err)
			}
		}
	}

	return refKind, refTag, className, name, descriptor, nil
}

// MethodHandleReference returns the validated symbolic reference target of a
// CONSTANT_MethodHandle. It delegates to validateMethodHandle.
func (cp *ConstantPool) MethodHandleReference(index uint16, majorVersion uint16) (refKind uint8, refTag uint8, className, name, descriptor string, err error) {
	return cp.validateMethodHandle(index, majorVersion)
}

// validateDynamicNat resolves a NameAndType index from a dynamic entry
// (InvokeDynamic or ConstantDynamic) using checked lookup, and validates the
// name and descriptor. The descriptorClass must be "method" or "field".
// Name rules differ by descriptor class:
//   - "method" (InvokeDynamic): must be a valid ordinary method name
//     (no '<' or '>', per JVMS §4.2.2).
//   - "field" (ConstantDynamic): must be a valid unqualified name, with
//     <init> and <clinit> explicitly rejected per JVMS §5.4.3.6.
func (cp *ConstantPool) validateDynamicNat(natIndex uint16, label string, descriptorClass string) (name, descriptor string, err error) {
	nat, err := cp.lookupNameAndType(natIndex, label+" name_and_type_index")
	if err != nil {
		return "", "", err
	}
	name, err = cp.lookupUtf8(nat.nameIndex, label+" name_index")
	if err != nil {
		return "", "", err
	}
	descriptor, err = cp.lookupUtf8(nat.descIndex, label+" descriptor_index")
	if err != nil {
		return "", "", err
	}
	switch descriptorClass {
	case "method":
		// InvokeDynamic: the NameAndType represents a method descriptor
		// (JVMS §4.4.10). The name must be a valid ordinary method name,
		// which inherently rejects <init>, <clinit>, and any name
		// containing '<' or '>'.
		if err := validateOrdinaryMethodName(name); err != nil {
			return "", "", fmt.Errorf("%s: invalid method name %q: %w", label, name, err)
		}
		if err := validateMethodDescriptor(descriptor); err != nil {
			return "", "", fmt.Errorf("%s: invalid method descriptor %q: %w", label, descriptor, err)
		}
	case "field":
		// ConstantDynamic: the NameAndType represents a field descriptor
		// (JVMS §4.4.10). Field names are unqualified names, but
		// <init>/<clinit> are explicitly rejected per JVMS §5.4.3.6.
		if err := validateUnqualifiedName(name); err != nil {
			return "", "", fmt.Errorf("%s: invalid field name %q: %w", label, name, err)
		}
		if name == "<init>" || name == "<clinit>" {
			return "", "", fmt.Errorf("%s: dynamic entry must not name %q", label, name)
		}
		if !validFieldDescriptor(descriptor) {
			return "", "", fmt.Errorf("%s: invalid field descriptor %q", label, descriptor)
		}
	}
	return name, descriptor, nil
}

// InvokeDynamicInfo returns the bootstrap table index and (name, method
// descriptor) for a CONSTANT_InvokeDynamic entry. Uses checked lookup
// throughout — never returns silently empty values for tag mismatches.
func (cp *ConstantPool) InvokeDynamicInfo(index uint16) (bootstrapIndex uint16, name, descriptor string, err error) {
	info, err := cp.lookup(index)
	if err != nil {
		return 0, "", "", fmt.Errorf("constant pool index %d: %w", index, err)
	}
	indy, ok := info.(*ConstantInvokeDynamicInfo)
	if !ok {
		return 0, "", "", fmt.Errorf("constant pool index %d: expected CONSTANT_InvokeDynamic, got tag %d", index, info.Tag())
	}
	name, descriptor, err = cp.validateDynamicNat(indy.natIndex,
		fmt.Sprintf("constant pool index %d", index), "method")
	if err != nil {
		return 0, "", "", err
	}
	return indy.bootstrapMethodAttrIndex, name, descriptor, nil
}

// ConstantDynamicInfo returns the bootstrap table index and (name, field
// descriptor) for a CONSTANT_Dynamic entry. Uses checked lookup throughout.
func (cp *ConstantPool) ConstantDynamicInfo(index uint16) (bootstrapIndex uint16, name, descriptor string, err error) {
	info, err := cp.lookup(index)
	if err != nil {
		return 0, "", "", fmt.Errorf("constant pool index %d: %w", index, err)
	}
	cdy, ok := info.(*ConstantDynamicInfo)
	if !ok {
		return 0, "", "", fmt.Errorf("constant pool index %d: expected CONSTANT_Dynamic, got tag %d", index, info.Tag())
	}
	name, descriptor, err = cp.validateDynamicNat(cdy.natIndex,
		fmt.Sprintf("constant pool index %d", index), "field")
	if err != nil {
		return 0, "", "", err
	}
	return cdy.bootstrapMethodAttrIndex, name, descriptor, nil
}

// --- Parse-time validation ---

// loadableConstant checks whether a constant-pool tag represents a loadable
// constant suitable as a bootstrap argument (JVMS §5.4.3.6).
func loadableConstant(tag uint8) bool {
	switch tag {
	case ConstantInteger, ConstantFloat, ConstantLong, ConstantDouble,
		ConstantClass, ConstantString, ConstantMethodHandle,
		ConstantMethodType, ConstantDynamic:
		return true
	default:
		return false
	}
}

// validateDynamicPool performs phase-2 cross-reference validation after the
// constant pool is fully parsed. It checks MethodType descriptors, MethodHandle
// targets (via validateMethodHandle), dynamic entry cross-references, and
// BootstrapMethods. Validation failures are reported as parsePanic so Parse
// converts them to *FormatError.
func validateDynamicPool(cp *ConstantPool, majorVersion uint16, bootstrapMethods *BootstrapMethodsAttr) {
	dynamicSeen := false

	for i := 1; i < len(cp.infos); i++ {
		info := cp.infos[i]
		if info == nil {
			continue
		}
		switch info.Tag() {
		case ConstantMethodType:
			mt := info.(*ConstantMethodTypeInfo)
			// Use checked lookup for the descriptor.
			_, err := cp.lookupUtf8(mt.descIndex, fmt.Sprintf("constant pool index %d: MethodType descriptor_index", i))
			if err != nil {
				panicf("constant pool", "%v", err)
			}
			desc := cp.UTF8(mt.descIndex) // safe: we just verified it's Utf8
			if err := validateMethodDescriptor(desc); err != nil {
				panicf("constant pool", "index %d: CONSTANT_MethodType descriptor %q is not a method descriptor: %s", i, desc, err)
			}

		case ConstantMethodHandle:
			_, _, _, _, _, err := cp.validateMethodHandle(uint16(i), majorVersion)
			if err != nil {
				panicf("constant pool", "%v", err)
			}

		case ConstantInvokeDynamic:
			dynamicSeen = true
			indy := info.(*ConstantInvokeDynamicInfo)
			// Use checked validation for the nat entry.
			_, _, err := cp.validateDynamicNat(indy.natIndex,
				fmt.Sprintf("constant pool index %d", i), "method")
			if err != nil {
				panicf("constant pool", "%v", err)
			}
			if bootstrapMethods == nil {
				panicf("constant pool", "index %d: InvokeDynamic references bootstrap method %d but no BootstrapMethods attribute present", i, indy.bootstrapMethodAttrIndex)
			}
			if int(indy.bootstrapMethodAttrIndex) >= len(bootstrapMethods.entries) {
				panicf("constant pool", "index %d: InvokeDynamic bootstrap method index %d out of range (table has %d entries)", i, indy.bootstrapMethodAttrIndex, len(bootstrapMethods.entries))
			}

		case ConstantDynamic:
			dynamicSeen = true
			cdy := info.(*ConstantDynamicInfo)
			_, _, err := cp.validateDynamicNat(cdy.natIndex,
				fmt.Sprintf("constant pool index %d", i), "field")
			if err != nil {
				panicf("constant pool", "%v", err)
			}
			if bootstrapMethods == nil {
				panicf("constant pool", "index %d: ConstantDynamic references bootstrap method %d but no BootstrapMethods attribute present", i, cdy.bootstrapMethodAttrIndex)
			}
			if int(cdy.bootstrapMethodAttrIndex) >= len(bootstrapMethods.entries) {
				panicf("constant pool", "index %d: ConstantDynamic bootstrap method index %d out of range (table has %d entries)", i, cdy.bootstrapMethodAttrIndex, len(bootstrapMethods.entries))
			}
		}
		// Skip the second slot of Long/Double entries.
		if info.Tag() == ConstantLong || info.Tag() == ConstantDouble {
			i++
		}
	}

	// Validate BootstrapMethods entries against the constant pool.
	if bootstrapMethods != nil {
		bsSize := uint16(len(cp.infos))
		for i, entry := range bootstrapMethods.entries {
			if entry.MethodRef == 0 || entry.MethodRef >= bsSize {
				panicf("BootstrapMethods", "entry %d: bootstrap_method_ref index %d out of constant pool bounds", i, entry.MethodRef)
			}
			mrInfo := cp.get(entry.MethodRef)
			if mrInfo.Tag() != ConstantMethodHandle {
				panicf("BootstrapMethods", "entry %d: bootstrap_method_ref index %d is not CONSTANT_MethodHandle (tag %d)", i, entry.MethodRef, mrInfo.Tag())
			}
			for j, argIdx := range entry.Arguments {
				if argIdx == 0 || argIdx >= bsSize {
					panicf("BootstrapMethods", "entry %d argument %d: index %d out of constant pool bounds", i, j, argIdx)
				}
				argInfo := cp.get(argIdx)
				if !loadableConstant(argInfo.Tag()) {
					panicf("BootstrapMethods", "entry %d argument %d: index %d has tag %d, not a loadable constant", i, j, argIdx, argInfo.Tag())
				}
			}
		}
	} else if dynamicSeen {
		panicf("BootstrapMethods", "dynamic entries present but no BootstrapMethods attribute found")
	}
}

// --- Name validation (JVMS §4.2.2) ---

// validateUnqualifiedName checks that name is a valid JVM unqualified name
// (JVMS §4.2.2): non-empty, no '.', ';', '[', or '/'. Used for field names,
// which may include '<' and '>' — <init> and <clinit> are valid field
// identifiers.
func validateUnqualifiedName(name string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}
	if strings.ContainsAny(name, ".;[/") {
		return fmt.Errorf("invalid character in name %q", name)
	}
	return nil
}

// validateOrdinaryMethodName checks that name is a valid JVM ordinary method
// name (JVMS §4.2.2): non-empty, no '.', ';', '[', '/', '<', or '>'.
// This rejects <init>, <clinit>, and any other name containing angle brackets.
func validateOrdinaryMethodName(name string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}
	if strings.ContainsAny(name, ".;[/<>") {
		return fmt.Errorf("invalid character in name %q", name)
	}
	return nil
}

// --- Descriptor grammar validation (JVMS §4.3.2, §4.3.3) ---

// validateInternalName checks that a JVM internal class name contains no illegal
// characters and has non-empty parts.
func validateInternalName(name string) error {
	if name == "" {
		return fmt.Errorf("empty internal name")
	}
	parts := strings.Split(name, "/")
	for _, p := range parts {
		if p == "" {
			return fmt.Errorf("empty part in internal name %q", name)
		}
	}
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '.', ';', '[':
			return fmt.Errorf("invalid character %c in internal name %q", name[i], name)
		}
	}
	return nil
}

// validateFieldDescriptor parses a field descriptor (JVMS §4.3.2). It returns
// the number of bytes consumed.
func validateFieldDescriptor(desc string) (int, error) {
	if desc == "" {
		return 0, fmt.Errorf("empty field descriptor")
	}
	switch desc[0] {
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z':
		return 1, nil
	case 'L':
		idx := strings.IndexByte(desc, ';')
		if idx == -1 {
			return 0, fmt.Errorf("object type missing ';' in %q", desc)
		}
		if idx == 1 {
			return 0, fmt.Errorf("empty class name in object type")
		}
		if err := validateInternalName(desc[1:idx]); err != nil {
			return 0, err
		}
		return idx + 1, nil
	case '[':
		n, err := validateFieldDescriptor(desc[1:])
		if err != nil {
			return 0, fmt.Errorf("array component: %w", err)
		}
		return 1 + n, nil
	default:
		return 0, fmt.Errorf("invalid base type %c in descriptor %q", desc[0], desc)
	}
}

// validateMethodDescriptor validates a method descriptor (JVMS §4.3.3).
func validateMethodDescriptor(desc string) error {
	if desc == "" || desc[0] != '(' {
		return fmt.Errorf("method descriptor must start with '('")
	}
	i := 1
	for i < len(desc) && desc[i] != ')' {
		n, err := validateFieldDescriptor(desc[i:])
		if err != nil {
			return fmt.Errorf("parameter descriptor at offset %d: %w", i, err)
		}
		i += n
	}
	if i >= len(desc) {
		return fmt.Errorf("method descriptor missing ')'")
	}
	i++ // skip ')'
	if i >= len(desc) {
		return fmt.Errorf("method descriptor missing return descriptor")
	}
	if desc[i] == 'V' {
		i++
	} else {
		n, err := validateFieldDescriptor(desc[i:])
		if err != nil {
			return fmt.Errorf("return descriptor: %w", err)
		}
		i += n
	}
	if i != len(desc) {
		return fmt.Errorf("trailing bytes in method descriptor %q at offset %d", desc, i)
	}
	return nil
}

// validFieldDescriptor reports whether desc is a valid field descriptor.
func validFieldDescriptor(desc string) bool {
	n, err := validateFieldDescriptor(desc)
	return err == nil && n == len(desc)
}
