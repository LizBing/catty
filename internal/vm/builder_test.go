package vm

import (
	"math"
	"testing"

	"catty/internal/classfile"
	"catty/internal/kernel"
)

// classBuilder assembles minimal-but-valid v52 class files so interpreter
// unit tests can run hand-written bytecode without javac.
type cpEntry struct {
	tag  classfile.ConstTag
	slot uint16 // resolved pool slot index (long/double consume 2)
	c1, c2 uint16
	str  string
	i32  int32
	i64v int64
	f32v float32
	f64v float64
}

type classBuilder struct {
	name     string
	cp       []cpEntry // 0 unused
	nextSlot uint16
	methods  []methodBlob
}

type methodBlob struct {
	flags      uint16
	name       string
	desc       string
	maxStack   uint16
	maxLocals  uint16
	code       []byte
	handlers   []excHandler
}

type excHandler struct {
	start, end, handler uint16
	catchClassName      string // "" = catch-all
}

func newClassBuilder(name string) *classBuilder {
	// Slot numbering starts at 1: pool index 0 is reserved (JVMS §4.4).
	return &classBuilder{name: name, cp: []cpEntry{{}}, nextSlot: 1}
}

// add registers an entry and returns its pool slot index. Category-2
// constants (long/double) occupy two slots per JVMS §4.4.5.
func (b *classBuilder) add(e cpEntry) uint16 {
	e.slot = b.nextSlot
	width := uint16(1)
	if e.tag == classfile.CLong || e.tag == classfile.CDouble {
		width = 2
	}
	b.nextSlot += width
	b.cp = append(b.cp, e)
	return e.slot
}

func (b *classBuilder) utf8(s string) uint16 { return b.add(cpEntry{tag: classfile.CUtf8, str: s}) }

func (b *classBuilder) class(name string) uint16 {
	return b.add(cpEntry{tag: classfile.CClass, c1: b.utf8(name)})
}

func (b *classBuilder) stringConst(s string) uint16 {
	return b.add(cpEntry{tag: classfile.CString, c1: b.utf8(s)})
}

func (b *classBuilder) nameAndType(name, desc string) uint16 {
	return b.add(cpEntry{tag: classfile.CNameAndType, c1: b.utf8(name), c2: b.utf8(desc)})
}

func (b *classBuilder) memberRef(tag classfile.ConstTag, cls, name, desc string) uint16 {
	return b.add(cpEntry{tag: tag, c1: b.class(cls), c2: b.nameAndType(name, desc)})
}

func (b *classBuilder) integer(v int32) uint16 { return b.add(cpEntry{tag: classfile.CInteger, i32: v}) }
func (b *classBuilder) long(v int64) uint16    { return b.add(cpEntry{tag: classfile.CLong, i64v: v}) }
func (b *classBuilder) float(v float32) uint16 { return b.add(cpEntry{tag: classfile.CFloat, f32v: v}) }
func (b *classBuilder) double(v float64) uint16 {
	return b.add(cpEntry{tag: classfile.CDouble, f64v: v})
}

func (b *classBuilder) addMethod(m methodBlob) { b.methods = append(b.methods, m) }

func u2(w []byte, v uint16) []byte { return append(w, byte(v>>8), byte(v)) }
func u4(w []byte, v uint32) []byte {
	return append(w, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// build serializes and parses the class through the real pipeline, then
// loads it into a fresh kernel.
func (b *classBuilder) build(t *testing.T) (*kernel.Kernel, *kernel.Class) {
	t.Helper()

	// Reserve standard entries up front so ordering is stable.
	objUtf8 := b.utf8("java/lang/Object")
	objClsIdx := b.add(cpEntry{tag: classfile.CClass, c1: objUtf8})
	thisIdx := b.class(b.name)
	codeName := b.utf8("Code")

	// Serialize method attribute bodies first: handler class refs may grow
	// the constant pool, and the pool size must be final before assembly.
	type serMethod struct {
		head []byte // flags/name/desc + attr_count placeholder filler
		attr []byte // Code attribute payload
	}
	sers := make([]serMethod, len(b.methods))
	for i, m := range b.methods {
		attr := []byte{}
		attr = u2(attr, m.maxStack)
		attr = u2(attr, m.maxLocals)
		attr = u4(attr, uint32(len(m.code)))
		attr = append(attr, m.code...)
		attr = u2(attr, uint16(len(m.handlers)))
		for _, h := range m.handlers {
			ct := uint16(0)
			if h.catchClassName != "" {
				ct = b.class(h.catchClassName)
			}
			attr = u2(attr, h.start)
			attr = u2(attr, h.end)
			attr = u2(attr, h.handler)
			attr = u2(attr, ct)
		}
		attr = u2(attr, 0) // nested attributes

		head := []byte{}
		head = u2(head, m.flags)
		head = u2(head, b.utf8(m.name))
		head = u2(head, b.utf8(m.desc))
		sers[i] = serMethod{head: head, attr: attr}
	}

	// Pool is final. Valid indices are 1..nextSlot-1, so the count field
	// (which includes the reserved slot 0) equals nextSlot itself.
	slots := b.nextSlot

	w := []byte{}
	w = u4(w, classfile.Magic)
	w = u2(w, 0)  // minor
	w = u2(w, 52) // major
	w = u2(w, slots)
	for _, e := range b.cp[1:] {
		w = append(w, byte(e.tag))
		switch e.tag {
		case classfile.CUtf8:
			w = u2(w, uint16(len(e.str)))
			w = append(w, e.str...)
		case classfile.CInteger:
			w = u4(w, uint32(e.i32))
		case classfile.CFloat:
			w = u4(w, f32bits(e.f32v))
		case classfile.CLong:
			w = u4(w, uint32(uint64(e.i64v)>>32))
			w = u4(w, uint32(uint64(e.i64v)))
		case classfile.CDouble:
			w = u4(w, f64hi(e.f64v))
			w = u4(w, f64lo(e.f64v))
		case classfile.CClass, classfile.CString:
			w = u2(w, e.c1)
		case classfile.CFieldref, classfile.CMethodref, classfile.CInterfaceMethodref, classfile.CNameAndType:
			w = u2(w, e.c1)
			w = u2(w, e.c2)
		default:
			t.Fatalf("builder: unsupported tag %v", e.tag)
		}
	}

	w = u2(w, 0x0021) // public super
	w = u2(w, thisIdx)
	w = u2(w, objClsIdx)
	w = u2(w, 0) // interfaces
	w = u2(w, 0) // fields
	w = u2(w, uint16(len(sers)))
	for _, sm := range sers {
		w = append(w, sm.head...)
		w = u2(w, 1) // one attribute: Code
		w = u2(w, codeName)
		w = u4(w, uint32(len(sm.attr)))
		w = append(w, sm.attr...)
	}
	w = u2(w, 0) // class attributes

	// Interpreter unit tests target the engine, not the verifier: the
	// builder does not synthesize StackMapTable frames yet (DEBT-0009).
	k := kernel.New(kernel.Options{SkipVerify: true})
	cls, err := k.LoadClassBytes(w)
	if err != nil {
		t.Fatalf("build/load %s: %v", b.name, err)
	}
	return k, cls
}

func f32bits(f float32) uint32 { return math.Float32bits(f) }
func f64hi(f float64) uint32   { return uint32(math.Float64bits(f) >> 32) }
func f64lo(f float64) uint32   { return uint32(math.Float64bits(f)) }

// callStatic runs a static method of the built class by name+desc.
func callStatic(t *testing.T, k *kernel.Kernel, cls *kernel.Class, name, desc string, args []kernel.Value) (kernel.Value, error) {
	t.Helper()
	th := New(k)
	m, err := k.ResolveMethod(cls, name, desc)
	if err != nil {
		t.Fatalf("resolve %s: %v", name, err)
	}
	return th.Call(m, nil, args)
}
