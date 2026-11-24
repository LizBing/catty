package classfile

import (
	"encoding/binary"
	"fmt"

	"catty/rt"
)

// Field is one declared class field.
type Field struct {
	Owner      string
	Name       string
	Descriptor string
	Flags      uint16
	// Constant is the validated ConstantValue payload, or nil.
	Constant *rt.Constant
}

// LineNumber is one source line entry.
type LineNumber struct {
	PC   int32
	Line int
}

// Method is one declared class method with decoded, verified-shape bytecode.
type Method struct {
	Owner      string
	Name       string
	Descriptor string
	Flags      uint16

	MaxStack    int
	MaxLocals   int
	CodeLength  int
	Insns       []rt.Insn
	Refs        []rt.Ref
	Handlers    []rt.Handler
	Frames      []Frame
	HasStackMap bool
	LineNumbers []LineNumber
}

// Class is a parsed class definition.
type Class struct {
	Name        string
	Super       string
	Interfaces  []string
	Flags       uint16
	NestHost    string
	NestMembers []string
	Fields      []*Field
	Methods     []*Method
	SourceFile  string
	Major       uint16
	Minor       uint16
	Origin      string
}

// Parse reads one class file. The caller supplies the origin used in
// diagnostics. It applies the default (Gate 2) opcode inventory.
func Parse(data []byte, origin string) (*Class, error) {
	return ParseFor(data, origin, "")
}

// ParseFor reads one class file under one profile's opcode inventory. Gate 4
// decodes monitorenter and monitorexit as single-byte instructions; earlier
// profiles keep rejecting them at the same offset.
func ParseFor(data []byte, origin, profile string) (*Class, error) {
	r := &reader{data: data}
	magic, err := r.u4()
	if err != nil {
		return nil, err
	}
	if magic != 0xCAFEBABE {
		return nil, fmt.Errorf("bad class file magic 0x%08x", magic)
	}
	minor, err := r.u2()
	if err != nil {
		return nil, err
	}
	major, err := r.u2()
	if err != nil {
		return nil, err
	}
	if major != 65 || minor != 0 {
		return nil, fmt.Errorf("unsupported class file version %d.%d (expected 65.0)", major, minor)
	}
	cp, err := parseConstantPool(r)
	if err != nil {
		return nil, err
	}
	flags, err := r.u2()
	if err != nil {
		return nil, err
	}
	thisIndex, err := r.u2()
	if err != nil {
		return nil, err
	}
	name, err := cp.className(thisIndex)
	if err != nil {
		return nil, err
	}
	superIndex, err := r.u2()
	if err != nil {
		return nil, err
	}
	c := &Class{Name: name, Flags: flags, Origin: origin, Major: major, Minor: minor}
	if superIndex != 0 {
		c.Super, err = cp.className(superIndex)
		if err != nil {
			return nil, err
		}
	}
	interfaceCount, err := r.u2()
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(interfaceCount); i++ {
		idx, err := r.u2()
		if err != nil {
			return nil, err
		}
		in, err := cp.className(idx)
		if err != nil {
			return nil, err
		}
		c.Interfaces = append(c.Interfaces, in)
	}
	fieldCount, err := r.u2()
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(fieldCount); i++ {
		f, err := parseField(r, cp, name)
		if err != nil {
			return nil, fmt.Errorf("field %d: %w", i, err)
		}
		c.Fields = append(c.Fields, f)
	}
	methodCount, err := r.u2()
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(methodCount); i++ {
		m, err := parseMethod(r, cp, name, profile)
		if err != nil {
			return nil, fmt.Errorf("method %d: %w", i, err)
		}
		for _, prev := range c.Methods {
			if prev.Name == m.Name && prev.Descriptor == m.Descriptor {
				return nil, fmt.Errorf("duplicate method %s%s", m.Name, m.Descriptor)
			}
		}
		c.Methods = append(c.Methods, m)
	}
	err = parseAttributes(r, cp, func(attrName string, body []byte) error {
		switch attrName {
		case "SourceFile":
			if len(body) != 2 {
				return fmt.Errorf("malformed SourceFile attribute")
			}
			text, _, err := cp.utf8(binary.BigEndian.Uint16(body))
			if err != nil {
				return err
			}
			c.SourceFile = text
		case "NestHost":
			if len(body) != 2 {
				return fmt.Errorf("malformed NestHost attribute")
			}
			host, err := cp.className(binary.BigEndian.Uint16(body))
			if err != nil {
				return err
			}
			c.NestHost = host
		case "NestMembers":
			if len(body) < 2 || len(body) != 2+2*int(binary.BigEndian.Uint16(body)) {
				return fmt.Errorf("malformed NestMembers attribute")
			}
			count := int(binary.BigEndian.Uint16(body))
			for i := 0; i < count; i++ {
				member, err := cp.className(binary.BigEndian.Uint16(body[2+i*2:]))
				if err != nil {
					return err
				}
				c.NestMembers = append(c.NestMembers, member)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func parseConstantPool(r *reader) (*constantPool, error) {
	count, err := r.u2()
	if err != nil {
		return nil, err
	}
	cp := &constantPool{entries: make([]cpEntry, int(count))}
	for i := 1; i < int(count); i++ {
		tag, err := r.u1()
		if err != nil {
			return nil, err
		}
		e := cpEntry{tag: tag}
		switch tag {
		case 1:
			n, err := r.u2()
			if err != nil {
				return nil, err
			}
			raw, err := r.bytes(int(n))
			if err != nil {
				return nil, err
			}
			units, err := decodeMUTF8(raw)
			if err != nil {
				return nil, fmt.Errorf("constant pool %d: %w", i, err)
			}
			e.unit = units
			e.text = unitsText(units)
		case 3:
			v, err := r.u4()
			if err != nil {
				return nil, err
			}
			e.intVal = int32(v)
		case 4:
			v, err := r.u4()
			if err != nil {
				return nil, err
			}
			e.bits = uint64(v)
		case 5, 6:
			raw, err := r.bytes(8)
			if err != nil {
				return nil, err
			}
			e.bits = binary.BigEndian.Uint64(raw)
			cp.entries[i] = e
			i++
			if i >= int(count) {
				return nil, fmt.Errorf("constant pool double-slot entry overruns pool")
			}
			cp.entries[i] = cpEntry{}
			continue
		case 7, 8:
			v, err := r.u2()
			if err != nil {
				return nil, err
			}
			e.classIndex = v
			e.stringIndex = v
		case 9, 10, 11, 12, 17, 18:
			a, err := r.u2()
			if err != nil {
				return nil, err
			}
			b, err := r.u2()
			if err != nil {
				return nil, err
			}
			e.first = a
			e.second = b
			e.nameTypeIndex = b
			if tag == 12 {
				e.nameIndex = a
				e.descriptorIndex = b
			}
			if tag == 17 || tag == 18 {
				e.bootstrapIndex = a
			}
		case 15:
			k, err := r.u1()
			if err != nil {
				return nil, err
			}
			v, err := r.u2()
			if err != nil {
				return nil, err
			}
			e.refKind = k
			e.refIndex = v
		case 16:
			v, err := r.u2()
			if err != nil {
				return nil, err
			}
			e.descriptorIndex = v
		case 19, 20:
			if _, err := r.u2(); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported constant pool tag %d at index %d", tag, i)
		}
		cp.entries[i] = e
	}
	return cp, nil
}

func parseField(r *reader, cp *constantPool, owner string) (*Field, error) {
	flags, err := r.u2()
	if err != nil {
		return nil, err
	}
	nameIdx, err := r.u2()
	if err != nil {
		return nil, err
	}
	descIdx, err := r.u2()
	if err != nil {
		return nil, err
	}
	name, _, err := cp.utf8(nameIdx)
	if err != nil {
		return nil, err
	}
	desc, _, err := cp.utf8(descIdx)
	if err != nil {
		return nil, err
	}
	ft, err := rt.ParseFieldType(desc)
	if err != nil {
		return nil, err
	}
	f := &Field{Owner: owner, Name: name, Descriptor: desc, Flags: flags}
	seenConstant := false
	err = parseAttributes(r, cp, func(attrName string, body []byte) error {
		if attrName != "ConstantValue" {
			return nil
		}
		if seenConstant {
			return fmt.Errorf("duplicate ConstantValue attribute on %s", name)
		}
		seenConstant = true
		if len(body) != 2 {
			return fmt.Errorf("malformed ConstantValue attribute on %s", name)
		}
		if flags&rt.AccStatic == 0 {
			return fmt.Errorf("ConstantValue on non-static field %s", name)
		}
		idx := binary.BigEndian.Uint16(body)
		e, err := cp.entry(idx)
		if err != nil {
			return err
		}
		value, err := constantValueFor(ft, e, cp, idx)
		if err != nil {
			return fmt.Errorf("field %s: %w", name, err)
		}
		f.Constant = value
		return nil
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

// constantValueFor validates that a ConstantValue constant matches the field
// descriptor and converts it to the image representation.
func constantValueFor(ft rt.Type, e *cpEntry, cp *constantPool, idx uint16) (*rt.Constant, error) {
	switch ft.Kind {
	case 'I':
		if e.tag != 3 {
			return nil, fmt.Errorf("ConstantValue tag %d does not match descriptor %s", e.tag, ft.Descriptor)
		}
		return &rt.Constant{Kind: "int", Bits: uint64(uint32(e.intVal))}, nil
	case 'J':
		if e.tag != 5 {
			return nil, fmt.Errorf("ConstantValue tag %d does not match descriptor %s", e.tag, ft.Descriptor)
		}
		return &rt.Constant{Kind: "long", Bits: e.bits}, nil
	case 'F':
		if e.tag != 4 {
			return nil, fmt.Errorf("ConstantValue tag %d does not match descriptor %s", e.tag, ft.Descriptor)
		}
		return &rt.Constant{Kind: "float", Bits: e.bits}, nil
	case 'D':
		if e.tag != 6 {
			return nil, fmt.Errorf("ConstantValue tag %d does not match descriptor %s", e.tag, ft.Descriptor)
		}
		return &rt.Constant{Kind: "double", Bits: e.bits}, nil
	case 'L':
		if ft.ClassName != "java/lang/String" {
			return nil, fmt.Errorf("ConstantValue is not supported for descriptor %s", ft.Descriptor)
		}
		if e.tag != 8 {
			return nil, fmt.Errorf("ConstantValue tag %d does not match descriptor %s", e.tag, ft.Descriptor)
		}
		units, err := cp.ldcString(idx)
		if err != nil {
			return nil, err
		}
		return &rt.Constant{Kind: "string", Units: append([]uint16(nil), units...)}, nil
	case 'A':
		return nil, fmt.Errorf("ConstantValue is not supported for array descriptor %s", ft.Descriptor)
	default:
		return nil, fmt.Errorf("unsupported ConstantValue descriptor %s", ft.Descriptor)
	}
}

func parseMethod(r *reader, cp *constantPool, owner, profile string) (*Method, error) {
	flags, err := r.u2()
	if err != nil {
		return nil, err
	}
	nameIdx, err := r.u2()
	if err != nil {
		return nil, err
	}
	descIdx, err := r.u2()
	if err != nil {
		return nil, err
	}
	name, _, err := cp.utf8(nameIdx)
	if err != nil {
		return nil, err
	}
	desc, _, err := cp.utf8(descIdx)
	if err != nil {
		return nil, err
	}
	m := &Method{Owner: owner, Name: name, Descriptor: desc, Flags: flags}
	seenCode := false
	err = parseAttributes(r, cp, func(attrName string, body []byte) error {
		switch attrName {
		case "Code":
			if seenCode {
				return fmt.Errorf("duplicate Code attribute on %s%s", name, desc)
			}
			seenCode = true
			return parseCodeAttribute(cp, m, body, profile)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s.%s%s: %w", owner, name, desc, err)
	}
	return m, nil
}

func parseCodeAttribute(cp *constantPool, m *Method, body []byte, profile string) error {
	r := &reader{data: body}
	maxStack, err := r.u2()
	if err != nil {
		return err
	}
	maxLocals, err := r.u2()
	if err != nil {
		return err
	}
	codeLen, err := r.u4()
	if err != nil {
		return err
	}
	if codeLen == 0 || codeLen > 1<<20 {
		return fmt.Errorf("unsupported Code length %d", codeLen)
	}
	code, err := r.bytes(int(codeLen))
	if err != nil {
		return err
	}
	handlerCount, err := r.u2()
	if err != nil {
		return err
	}
	for i := 0; i < int(handlerCount); i++ {
		start, err := r.u2()
		if err != nil {
			return err
		}
		end, err := r.u2()
		if err != nil {
			return err
		}
		handler, err := r.u2()
		if err != nil {
			return err
		}
		catchIndex, err := r.u2()
		if err != nil {
			return err
		}
		h := rt.Handler{StartPC: int32(start), EndPC: int32(end), HandlerPC: int32(handler)}
		if catchIndex != 0 {
			h.CatchType, err = cp.className(catchIndex)
			if err != nil {
				return err
			}
		}
		m.Handlers = append(m.Handlers, h)
	}
	m.MaxStack = int(maxStack)
	m.MaxLocals = int(maxLocals)
	m.CodeLength = int(codeLen)
	insns, refs, err := decodeInstructions(cp, code, m, profile)
	if err != nil {
		return err
	}
	m.Insns = insns
	m.Refs = refs
	return parseAttributes(r, cp, func(attrName string, attr []byte) error {
		switch attrName {
		case "StackMapTable":
			frames, err := parseStackMapTable(cp, attr, m)
			if err != nil {
				return err
			}
			m.Frames = frames
			m.HasStackMap = true
		case "LineNumberTable":
			lines, err := parseLineNumberTable(attr)
			if err != nil {
				return err
			}
			m.LineNumbers = lines
		}
		return nil
	})
}

func parseLineNumberTable(body []byte) ([]LineNumber, error) {
	if len(body) < 2 {
		return nil, fmt.Errorf("malformed LineNumberTable")
	}
	count := int(binary.BigEndian.Uint16(body))
	if len(body) != 2+count*4 {
		return nil, fmt.Errorf("malformed LineNumberTable length")
	}
	out := make([]LineNumber, 0, count)
	for i := 0; i < count; i++ {
		pc := binary.BigEndian.Uint16(body[2+i*4:])
		line := binary.BigEndian.Uint16(body[4+i*4:])
		out = append(out, LineNumber{PC: int32(pc), Line: int(line)})
	}
	return out, nil
}

func parseAttributes(r *reader, cp *constantPool, handle func(name string, body []byte) error) error {
	count, err := r.u2()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		nameIdx, err := r.u2()
		if err != nil {
			return err
		}
		length, err := r.u4()
		if err != nil {
			return err
		}
		body, err := r.bytes(int(length))
		if err != nil {
			return err
		}
		name, _, err := cp.utf8(nameIdx)
		if err != nil {
			return err
		}
		if handle != nil {
			if err := handle(name, body); err != nil {
				return fmt.Errorf("attribute %s: %w", name, err)
			}
		}
	}
	return nil
}
