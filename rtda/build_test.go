package rtda

import (
	"encoding/binary"
	"testing"

	"catty/classfile"
)

// failingLoader is a Loader that returns typed failures for specific names.
type failingLoader struct {
	classes     map[string]*Class
	failureName string
	failureKind FailureKind
}

func (l *failingLoader) LoadClass(name string) *Class {
	if name == l.failureName {
		return nil
	}
	return l.classes[name]
}

func (l *failingLoader) LoadClassResult(name string) ClassLoadResult {
	if name == l.failureName {
		return NewFailureResult(&ClassLoadFailure{Kind: l.failureKind, Name: name})
	}
	if c, ok := l.classes[name]; ok {
		return NewClassResult(c)
	}
	return NewFailureResult(&ClassLoadFailure{Kind: FailureNotFound, Name: name})
}

func (l *failingLoader) LoaderIdentity() *LoaderIdentity {
	return NewLoaderIdentity()
}

// newMinimalClassFile builds a minimal valid .class file (raw bytes), parses it,
// and returns the ClassFile. The generated class has the given name, superclass,
// and interfaces, with no fields or methods. This avoids requiring javac.
func newMinimalClassFile(t *testing.T, name, superName string, interfaces []string) *classfile.ClassFile {
	t.Helper()

	// We build constant-pool entries manually:
	//   cp[1] = CONSTANT_Class -> utf8(name)
	//   cp[2] = CONSTANT_Utf8  -> name
	//   cp[3] = CONSTANT_Class -> utf8(superName)
	//   cp[4] = CONSTANT_Utf8  -> superName
	//   (interfaces and their utf8s follow)
	b := newClassFileBuilder()
	clsIdx := b.addClass(name)
	superIdx := b.addClass(superName)
	ifaceIdxs := make([]uint16, len(interfaces))
	for i, n := range interfaces {
		ifaceIdxs[i] = b.addClass(n)
	}
	buf := b.finish(clsIdx, superIdx, ifaceIdxs, 0x0021) // ACC_PUBLIC | ACC_SUPER

	cf, err := classfile.Parse(buf)
	if err != nil {
		t.Fatalf("parse minimal classfile: %v", err)
	}
	return cf
}

// classFileBuilder incrementally builds a valid .class file binary.
type classFileBuilder struct {
	cpCount   uint16
	cpBytes   []byte
	cpOffsets map[string]uint16 // utf8 content -> cp index
}

func newClassFileBuilder() *classFileBuilder {
	// cp[0] is unused slot.
	return &classFileBuilder{
		cpCount:   1,
		cpOffsets: make(map[string]uint16),
	}
}

func (b *classFileBuilder) addUtf8(s string) uint16 {
	if idx, ok := b.cpOffsets[s]; ok {
		return idx
	}
	idx := b.cpCount
	b.cpOffsets[s] = idx
	b.cpCount++

	// CONSTANT_Utf8 tag (1) + length (2) + bytes
	b.cpBytes = append(b.cpBytes, 1) // tag
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(s)))
	b.cpBytes = append(b.cpBytes, lenBytes...)
	b.cpBytes = append(b.cpBytes, []byte(s)...)
	return idx
}

func (b *classFileBuilder) addClass(name string) uint16 {
	utfIdx := b.addUtf8(name)
	idx := b.cpCount
	b.cpCount++

	// CONSTANT_Class tag (7) + name_index (2)
	b.cpBytes = append(b.cpBytes, 7) // tag
	idxBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(idxBytes, utfIdx)
	b.cpBytes = append(b.cpBytes, idxBytes...)
	return idx
}

func (b *classFileBuilder) finish(thisClass, superClass uint16, interfaces []uint16, accessFlags uint16) []byte {
	var buf []byte

	// Magic.
	buf = append(buf, 0xCA, 0xFE, 0xBA, 0xBE)
	// Version (61.0 = Java 17).
	ver := make([]byte, 4)
	binary.BigEndian.PutUint16(ver[0:2], 0)
	binary.BigEndian.PutUint16(ver[2:4], 61)
	buf = append(buf, ver...)
	// Constant pool count + entries.
	cpCount := make([]byte, 2)
	binary.BigEndian.PutUint16(cpCount, b.cpCount)
	buf = append(buf, cpCount...)
	buf = append(buf, b.cpBytes...)
	// Access flags.
	af := make([]byte, 2)
	binary.BigEndian.PutUint16(af, accessFlags)
	buf = append(buf, af...)
	// This class.
	tc := make([]byte, 2)
	binary.BigEndian.PutUint16(tc, thisClass)
	buf = append(buf, tc...)
	// Super class.
	sc := make([]byte, 2)
	binary.BigEndian.PutUint16(sc, superClass)
	buf = append(buf, sc...)
	// Interfaces count + indices.
	ifCount := make([]byte, 2)
	binary.BigEndian.PutUint16(ifCount, uint16(len(interfaces)))
	buf = append(buf, ifCount...)
	for _, idx := range interfaces {
		ib := make([]byte, 2)
		binary.BigEndian.PutUint16(ib, idx)
		buf = append(buf, ib...)
	}
	// Fields count (0), methods count (0), attributes count (0).
	buf = append(buf, 0, 0, 0, 0, 0, 0)
	return buf
}

func TestBuildClassSuperclassFailure(t *testing.T) {
	loader := &failingLoader{
		classes:     make(map[string]*Class),
		failureName: "test/BadSuper",
		failureKind: FailureNotFound,
	}

	cf := newMinimalClassFile(t, "test/Child", "test/BadSuper", nil)

	result := BuildClass(cf, loader)
	if result.Class != nil {
		t.Error("BuildClass succeeded unexpectedly for class with failed superclass")
	}
	if result.Failure == nil {
		t.Fatal("BuildClass returned nil failure for class with failed superclass")
	}
	if result.Failure.Kind != FailureNotFound {
		t.Errorf("failure kind = %v, want FailureNotFound", result.Failure.Kind)
	}
	if result.Failure.Name != "test/Child" {
		t.Errorf("failure name = %q, want test/Child", result.Failure.Name)
	}
	if result.Failure.Cause == nil {
		t.Error("failure Cause is nil")
	}
}

func TestBuildClassInterfaceFailure(t *testing.T) {
	// Need to provide java/lang/Object so the superclass resolves.
	obj := NewSyntheticClass("java/lang/Object", nil)
	loader := &failingLoader{
		classes:     map[string]*Class{"java/lang/Object": obj},
		failureName: "test/BadInterface",
		failureKind: FailureLinkage,
	}

	cf := newMinimalClassFile(t, "test/Child", "java/lang/Object", []string{"test/BadInterface"})

	result := BuildClass(cf, loader)
	if result.Class != nil {
		t.Error("BuildClass succeeded unexpectedly for class with failed interface")
	}
	if result.Failure == nil {
		t.Fatal("BuildClass returned nil failure for class with failed interface")
	}
	if result.Failure.Kind != FailureLinkage {
		t.Errorf("failure kind = %v, want FailureLinkage", result.Failure.Kind)
	}
}

func TestBuildClassSuccess(t *testing.T) {
	super := NewSyntheticClass("test/Parent", nil)
	loader := &failingLoader{
		classes:     map[string]*Class{"test/Parent": super},
		failureName: "nonexistent",
		failureKind: FailureNotFound,
	}

	cf := newMinimalClassFile(t, "test/Child", "test/Parent", nil)

	result := BuildClass(cf, loader)
	if result.Failure != nil {
		t.Fatalf("BuildClass failed unexpectedly: %v", result.Failure)
	}
	if result.Class == nil {
		t.Fatal("BuildClass returned nil Class on success")
	}
	if result.Class.Name() != "test/Child" {
		t.Errorf("class name = %q, want test/Child", result.Class.Name())
	}
	if result.Class.SuperClass() != super {
		t.Errorf("superclass = %p, want %p", result.Class.SuperClass(), super)
	}
}

func TestNewClassReturnsNilOnFailure(t *testing.T) {
	loader := &failingLoader{
		classes:     make(map[string]*Class),
		failureName: "test/BadSuper",
		failureKind: FailureNotFound,
	}

	cf := newMinimalClassFile(t, "test/Child", "test/BadSuper", nil)

	c := NewClass(cf, loader)
	if c != nil {
		t.Error("NewClass should return nil on superclass resolution failure")
	}
}
