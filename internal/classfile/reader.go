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
	Name       string
	Super      string
	Interfaces []string
	Flags      uint16
	Fields     []*Field
	Methods    []*Method
	SourceFile string
	Major      uint16
	Minor      uint16
	Origin     string
}

// Parse reads one class file. The caller supplies the origin used in
// diagnostics.
func Parse(data []byte, origin string) (*Class, error) {
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
		m, err := parseMethod(r, cp, name)
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
		if attrName != "SourceFile" {
			return nil
		}
		if len(body) != 2 {
			return fmt.Errorf("malformed SourceFile attribute")
		}
		text, _, err := cp.utf8(binary.BigEndian.Uint16(body))
		if err != nil {
			return err
		}
		c.SourceFile = text
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
			if _, err := r.u4(); err != nil {
				return nil, err
			}
		case 5, 6:
			if _, err := r.bytes(8); err != nil {
				return nil, err
			}
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
	f := &Field{Owner: owner, Name: name, Descriptor: desc, Flags: flags}
	if err := parseAttributes(r, cp, nil); err != nil {
		return nil, err
	}
	return f, nil
}

func parseMethod(r *reader, cp *constantPool, owner string) (*Method, error) {
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
			return parseCodeAttribute(cp, m, body)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s.%s%s: %w", owner, name, desc, err)
	}
	return m, nil
}

func parseCodeAttribute(cp *constantPool, m *Method, body []byte) error {
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
	insns, refs, err := decodeInstructions(cp, code, m)
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
