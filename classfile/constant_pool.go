package classfile

import "math"

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
		raw := reader.data[:length]
		reader.data = reader.data[length:]
		return &ConstantUtf8Info{str: decodeMUTF8(raw)}
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
	case ConstantInvokeDynamic:
		return &ConstantInvokeDynamicInfo{
			bootstrapMethodAttrIndex: reader.ReadUint16(),
			natIndex:                 reader.ReadUint16(),
		}
	default:
		panic("catty: unknown constant pool tag: " + itoa(int(tag)))
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
	ConstantInvokeDynamic      uint8 = 18
)

// --- Concrete constant info types ---

type ConstantUtf8Info struct{ str string }

func (c *ConstantUtf8Info) Tag() uint8 { return ConstantUtf8 }
func (c *ConstantUtf8Info) Str() string { return c.str }

type ConstantIntegerInfo struct{ val int32 }

func (c *ConstantIntegerInfo) Tag() uint8 { return ConstantInteger }
func (c *ConstantIntegerInfo) Val() int32 { return c.val }

type ConstantFloatInfo struct{ val float32 }

func (c *ConstantFloatInfo) Tag() uint8   { return ConstantFloat }
func (c *ConstantFloatInfo) Val() float32 { return c.val }

type ConstantLongInfo struct{ val int64 }

func (c *ConstantLongInfo) Tag() uint8   { return ConstantLong }
func (c *ConstantLongInfo) Val() int64   { return c.val }

type ConstantDoubleInfo struct{ val float64 }

func (c *ConstantDoubleInfo) Tag() uint8    { return ConstantDouble }
func (c *ConstantDoubleInfo) Val() float64  { return c.val }

type ConstantStringInfo struct {
	cp       *ConstantPool
	strIndex uint16
}

func (c *ConstantStringInfo) Tag() uint8    { return ConstantString }
func (c *ConstantStringInfo) Str() string   { return c.cp.UTF8(c.strIndex) }

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

func (c *ConstantNameAndTypeInfo) Tag() uint8           { return ConstantNameAndType }
func (c *ConstantNameAndTypeInfo) Name() string         { return c.cp.UTF8(c.nameIndex) }
func (c *ConstantNameAndTypeInfo) Descriptor() string   { return c.cp.UTF8(c.descIndex) }

// ConstantMemberRefInfo backs Fieldref / Methodref / InterfaceMethodref; they
// share an identical encoding and differ only in tag, so one type serves all.
type ConstantMemberRefInfo struct {
	cp         *ConstantPool
	tag        uint8
	classIndex uint16
	natIndex   uint16
}

func (c *ConstantMemberRefInfo) Tag() uint8             { return c.tag }
func (c *ConstantMemberRefInfo) ClassName() string      { return c.cp.ClassName(c.classIndex) }
func (c *ConstantMemberRefInfo) Name() string           { return c.cp.NameAndTypeName(c.natIndex) }
func (c *ConstantMemberRefInfo) Descriptor() string     { return c.cp.NameAndTypeDescriptor(c.natIndex) }

type ConstantMethodTypeInfo struct{ descIndex uint16 }

func (c *ConstantMethodTypeInfo) Tag() uint8 { return ConstantMethodType }

type ConstantMethodHandleInfo struct {
	refKind  uint8
	refIndex uint16
}

func (c *ConstantMethodHandleInfo) Tag() uint8 { return ConstantMethodHandle }

type ConstantInvokeDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	natIndex                 uint16
}

func (c *ConstantInvokeDynamicInfo) Tag() uint8 { return ConstantInvokeDynamic }

// --- Typed accessors on the pool ---

func (cp *ConstantPool) get(index uint16) ConstantInfo { return cp.infos[index] }

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

// Tag returns the JVMS constant tag at index, used by ldc to dispatch on type.
func (cp *ConstantPool) Tag(index uint16) uint8 {
	return cp.infos[index].Tag()
}

// itoa avoids importing strconv just for a panic message.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
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
