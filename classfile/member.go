package classfile

// MemberInfo describes a field or method (JVMS §4.5). The encoding is identical
// for both; only the access flags and presence of a Code attribute distinguish
// them in practice.
type MemberInfo struct {
	cp              *ConstantPool
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []AttributeInfo
}

// readMembers reads a u2 count followed by that many field_info or method_info
// structures. loc distinguishes field from method attributes.
func readMembers(reader *ClassReader, cp *ConstantPool, loc attrLocation) []*MemberInfo {
	count := reader.ReadUint16()
	members := make([]*MemberInfo, count)
	for i := range members {
		members[i] = readMember(reader, cp, loc)
	}
	return members
}

func readMember(reader *ClassReader, cp *ConstantPool, loc attrLocation) *MemberInfo {
	return &MemberInfo{
		cp:              cp,
		accessFlags:     reader.ReadUint16(),
		nameIndex:       reader.ReadUint16(),
		descriptorIndex: reader.ReadUint16(),
		attributes:      readAttributes(reader, cp, loc),
	}
}

func (m *MemberInfo) AccessFlags() uint16         { return m.accessFlags }
func (m *MemberInfo) Name() string                { return m.cp.UTF8(m.nameIndex) }
func (m *MemberInfo) Descriptor() string          { return m.cp.UTF8(m.descriptorIndex) }
func (m *MemberInfo) Attributes() []AttributeInfo { return m.attributes }

// Code returns the parsed Code attribute (methods) or nil (fields/abstract).
func (m *MemberInfo) Code() *CodeAttribute { return CodeAttributeOf(m.attributes) }
