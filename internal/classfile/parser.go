package classfile

import (
	"fmt"
	"math"
)

// reader is a big-endian byte reader with panic-based short-read handling.
// Parse recovers the panic and converts it to an error.
type reader struct {
	data []byte
	pos  int
}

func (r *reader) u1() uint8 {
	if r.pos >= len(r.data) {
		panic(shortRead{need: 1})
	}
	v := r.data[r.pos]
	r.pos++
	return v
}

func (r *reader) u2() uint16 {
	if r.pos+2 > len(r.data) {
		panic(shortRead{need: 2})
	}
	v := uint16(r.data[r.pos])<<8 | uint16(r.data[r.pos+1])
	r.pos += 2
	return v
}

func (r *reader) u4() uint32 {
	if r.pos+4 > len(r.data) {
		panic(shortRead{need: 4})
	}
	v := uint32(r.data[r.pos])<<24 | uint32(r.data[r.pos+1])<<16 |
		uint32(r.data[r.pos+2])<<8 | uint32(r.data[r.pos+3])
	r.pos += 4
	return v
}

func (r *reader) u8() uint64 {
	hi, lo := r.u4(), r.u4()
	return uint64(hi)<<32 | uint64(lo)
}

func (r *reader) bytes(n int) []byte {
	if n < 0 || r.pos+n > len(r.data) {
		panic(shortRead{need: n})
	}
	b := r.data[r.pos : r.pos+n]
	r.pos += n
	return b
}

type shortRead struct{ need int }

func (s shortRead) Error() string {
	return fmt.Sprintf("short read: need %d more bytes", s.need)
}

// Parse decodes a class file image. M0 accepts major versions 45..52.
// Verification (JVMS §4.10) is out of scope for M0: input is trusted
// javac output; see docs/specifications/deviation-ledger.md.
func Parse(data []byte) (cf *ClassFile, err error) {
	defer func() {
		if p := recover(); p != nil {
			if sr, ok := p.(shortRead); ok {
				cf, err = nil, fmt.Errorf("truncated class file: %v", sr)
				return
			}
			panic(p)
		}
	}()

	r := &reader{data: data}
	if m := r.u4(); m != Magic {
		return nil, fmt.Errorf("bad magic %#08x (want %#08x)", m, Magic)
	}
	minor := r.u2()
	major := r.u2()
	if major < 45 || major > MaxSupportedMajor {
		return nil, fmt.Errorf("unsupported class file version %d.%d (parser supports %d..%d)",
			major, minor, 45, MaxSupportedMajor)
	}

	cf = &ClassFile{MinorVersion: minor, MajorVersion: major}

	// Constant pool (JVMS §4.4). Long/Double occupy two slots.
	cpCount := int(r.u2())
	cf.Constants = make([]ConstEntry, cpCount)
	for i := 1; i < cpCount; i++ {
		e := &cf.Constants[i]
		e.Tag = ConstTag(r.u1())
		switch e.Tag {
		case CUtf8:
			n := int(r.u2())
			e.Str = decodeMUTF8(r.bytes(n))
		case CInteger:
			e.IntVal = int32(r.u4())
		case CFloat:
			e.FloatVal = math.Float32frombits(r.u4())
		case CLong:
			e.LongVal = int64(r.u8())
			i++ // two slots
		case CDouble:
			e.DoubleVal = math.Float64frombits(r.u8())
			i++
		case CClass, CString, CMethodType, CModule, CPackage:
			e.Idx1 = r.u2()
		case CFieldref, CMethodref, CInterfaceMethodref, CNameAndType, CDynamic, CInvokeDynamic:
			e.Idx1 = r.u2()
			e.Idx2 = r.u2()
		case CMethodHandle:
			e.Idx1 = uint16(r.u1()) // reference_kind
			e.Idx2 = r.u2()
		default:
			return nil, fmt.Errorf("unknown constant pool tag %d at index %d", e.Tag, i)
		}
	}

	cf.AccessFlags = r.u2()

	var err2 error
	if cf.ThisClass, err2 = cf.ClassName(r.u2()); err2 != nil {
		return nil, err2
	}
	superIdx := r.u2()
	if superIdx == 0 {
		cf.SuperClass = "" // java/lang/Object itself
	} else if cf.SuperClass, err2 = cf.ClassName(superIdx); err2 != nil {
		return nil, err2
	}

	ifcCount := int(r.u2())
	cf.Interfaces = make([]string, ifcCount)
	for i := 0; i < ifcCount; i++ {
		if cf.Interfaces[i], err2 = cf.ClassName(r.u2()); err2 != nil {
			return nil, err2
		}
	}

	fieldCount := int(r.u2())
	for i := 0; i < fieldCount; i++ {
		fd, err := parseField(r, cf)
		if err != nil {
			return nil, fmt.Errorf("field #%d: %w", i, err)
		}
		cf.Fields = append(cf.Fields, *fd)
	}

	methodCount := int(r.u2())
	for i := 0; i < methodCount; i++ {
		mi, err := parseMethod(r, cf)
		if err != nil {
			return nil, fmt.Errorf("method #%d: %w", i, err)
		}
		cf.Methods = append(cf.Methods, *mi)
	}

	if cf.Attributes, err = parseAttributes(r, cf); err != nil {
		return nil, err
	}
	for _, a := range cf.Attributes {
		if a.Name == "RuntimeVisibleAnnotations" || a.Name == "RuntimeInvisibleAnnotations" {
			cf.Annotations = append(cf.Annotations, parseAnnotations(cf, a.Data)...)
		}
		if a.Name == "SourceFile" && len(a.Data) == 2 {
			idx := uint16(a.Data[0])<<8 | uint16(a.Data[1])
			cf.SourceFile, _ = cf.UTF8(idx)
		}
	}

	if r.pos != len(data) {
		return nil, fmt.Errorf("trailing bytes after class file: consumed %d of %d", r.pos, len(data))
	}
	return cf, nil
}

func parseField(r *reader, cf *ClassFile) (*FieldInfo, error) {
	fd := &FieldInfo{AccessFlags: r.u2()}
	var err error
	if fd.Name, err = cf.UTF8(r.u2()); err != nil {
		return nil, err
	}
	if fd.Desc, err = cf.UTF8(r.u2()); err != nil {
		return nil, err
	}
	if fd.Attributes, err = parseAttributes(r, cf); err != nil {
		return nil, err
	}
	for _, a := range fd.Attributes {
		if a.Name == "RuntimeVisibleAnnotations" || a.Name == "RuntimeInvisibleAnnotations" {
			fd.Annotations = append(fd.Annotations, parseAnnotations(cf, a.Data)...)
		}
	}
	return fd, nil
}

func parseMethod(r *reader, cf *ClassFile) (*MethodInfo, error) {
	mi := &MethodInfo{AccessFlags: r.u2()}
	var err error
	if mi.Name, err = cf.UTF8(r.u2()); err != nil {
		return nil, err
	}
	if mi.Desc, err = cf.UTF8(r.u2()); err != nil {
		return nil, err
	}
	attrs, err := parseAttributes(r, cf)
	if err != nil {
		return nil, err
	}
	mi.Attributes = attrs
	for _, a := range attrs {
		if a.Name == "Code" {
			c, err := parseCode(cf, a.Data)
			if err != nil {
				return nil, fmt.Errorf("Code attr of %s%s: %w", mi.Name, mi.Desc, err)
			}
			mi.Code = c
			break
		}
	}
		for _, a := range attrs {
		if a.Name == "RuntimeVisibleAnnotations" || a.Name == "RuntimeInvisibleAnnotations" {
			mi.Annotations = append(mi.Annotations, parseAnnotations(cf, a.Data)...)
		}
	}
	return mi, nil
}

// parseCode decodes the Code attribute (JVMS §4.7.3).
func parseCode(cf *ClassFile, data []byte) (*Code, error) {
	r := &reader{data: data}
	c := &Code{MaxStack: r.u2(), MaxLocals: r.u2()}
	codeLen := int(r.u4())
	c.Code = append([]byte(nil), r.bytes(codeLen)...)
	exCount := int(r.u2())
	if exCount > 0 {
		c.Handlers = make([]ExceptionHandler, exCount)
		for i := 0; i < exCount; i++ {
			c.Handlers[i] = ExceptionHandler{
				StartPc:   r.u2(),
				EndPc:     r.u2(),
				HandlerPc: r.u2(),
				CatchType: r.u2(),
			}
		}
	}
	// Nested attributes: LineNumberTable etc. are skipped by length;
	// StackMapTable is decoded for the verifier.
	nested, err := parseAttributes(r, cf)
	if err != nil {
		return nil, err
	}
	for _, a := range nested {
		if a.Name == "StackMapTable" {
			c.StackMaps, err = ParseStackMapTable(cf, a.Data)
			if err != nil {
				return nil, err
			}
			break
		}
		if a.Name == "LineNumberTable" && len(a.Data) >= 2 {
			n := int(uint16(a.Data[0])<<8 | uint16(a.Data[1]))
			lns := make([]LineNum, 0, n)
			for i := 0; i < n && 2+4*i+4 <= len(a.Data); i++ {
				o := 2 + 4*i
				lns = append(lns, LineNum{
					StartPc: uint16(a.Data[o])<<8 | uint16(a.Data[o+1]),
					Line:    uint16(a.Data[o+2])<<8 | uint16(a.Data[o+3]),
				})
			}
			c.LineNumbers = lns
		}
	}
	return c, nil
}

// parseAttributes reads an attribute table, resolving names and capturing
// raw payloads. Unknown attributes are preserved but not interpreted.
// parseAnnotations decodes a Runtime{Visible,Invisible}Annotations payload
// (JVMS §4.7.16-17).
func parseAnnotations(cf *ClassFile, data []byte) []ParsedAnnotation {
	if len(data) < 2 {
		return nil
	}
	n := int(uint16(data[0])<<8 | uint16(data[1]))
	out := make([]ParsedAnnotation, 0, n)
	pos := 2
	for i := 0; i < n && pos+4 <= len(data); i++ {
		typeIdx := uint16(data[pos])<<8 | uint16(data[pos+1])
		pos += 2
		td, err := cf.UTF8(typeIdx)
		if err != nil {
			break
		}
		pa := ParsedAnnotation{TypeDesc: td, Elements: map[string]ElementValue{}}
		numPairs := int(uint16(data[pos])<<8 | uint16(data[pos+1]))
		pos += 2
		for j := 0; j < numPairs && pos+3 <= len(data); j++ {
			nameIdx := uint16(data[pos])<<8 | uint16(data[pos+1])
			pos += 2
			elemName, err := cf.UTF8(nameIdx)
			if err != nil {
				break
			}
			ev, consumed := decodeElementValue(cf, data, pos)
			if consumed == 0 {
				break
			}
			pos += consumed
			pa.Elements[elemName] = ev
		}
		out = append(out, pa)
	}
	return out
}

func decodeElementValue(cf *ClassFile, data []byte, pos int) (ElementValue, int) {
	if pos >= len(data) {
		return ElementValue{}, 0
	}
	tag := data[pos]
	pos++
	switch tag {
	case 'B', 'C', 'D', 'I', 'J', 'S', 'Z', 's':
		if pos+2 > len(data) {
			return ElementValue{}, 0
		}
		idx := uint16(data[pos])<<8 | uint16(data[pos+1])
		ev := ElementValue{Tag: tag, ConstIdx: idx}
		if tag == 's' {
			if sv, err := cf.UTF8(idx); err == nil {
				ev.EnumName = sv // reuse field for string value
			}
		}
		return ev, 3
	case 'e':
		if pos+4 > len(data) {
			return ElementValue{}, 0
		}
		typeIdx := uint16(data[pos])<<8 | uint16(data[pos+1])
		nameIdx := uint16(data[pos+2])<<8 | uint16(data[pos+3])
		td, _ := cf.UTF8(typeIdx)
		en, _ := cf.UTF8(nameIdx)
		return ElementValue{Tag: 'e', EnumType: td, EnumName: en}, 5
	case 'c':
		if pos+2 > len(data) {
			return ElementValue{}, 0
		}
		idx := uint16(data[pos])<<8 | uint16(data[pos+1])
		cd, _ := cf.UTF8(idx)
		return ElementValue{Tag: 'c', ClassDesc: cd}, 3
	case '@':
		// nested annotation — v1 skips (returns empty marker)
		return ElementValue{Tag: '@'}, 0 // 0 = abort outer parse
	case '[':
		if pos+2 > len(data) {
			return ElementValue{}, 0
		}
		count := int(uint16(data[pos])<<8 | uint16(data[pos+1]))
		startPos := pos
		pos += 2
		arr := ElementValue{Tag: '[', Array: make([]ElementValue, 0, count)}
		for i := 0; i < count; i++ {
			ev, used := decodeElementValue(cf, data, pos)
			if used == 0 {
				break
			}
			arr.Array = append(arr.Array, ev)
			pos += used
		}
		return arr, pos - startPos
	default:
		return ElementValue{}, 0
	}
}

func parseAttributes(r *reader, cf *ClassFile) ([]Attribute, error) {
	count := int(r.u2())
	if count == 0 {
		return nil, nil
	}
	attrs := make([]Attribute, count)
	for i := 0; i < count; i++ {
		name, err := cf.UTF8(r.u2())
		if err != nil {
			return nil, err
		}
		n := int(r.u4())
		attrs[i] = Attribute{Name: name, Data: append([]byte(nil), r.bytes(n)...)}
	}
	return attrs, nil
}
