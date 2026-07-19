package classfile

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// compileFixture compiles tests/fixtures/<name>.java into t.TempDir() and
// returns the bytes of the resulting <name>.class. It fails the test (rather
// than skipping) if javac is missing — the whole project depends on a JDK.
func compileFixture(t *testing.T, name string) []byte {
	t.Helper()
	src := filepath.Join("..", "tests", "fixtures", name+".java")
	out := t.TempDir()
	cmd := exec.Command("javac", "-source", "8", "-target", "8", "-d", out, src)
	if err := cmd.Run(); err != nil {
		t.Fatalf("javac %s failed: %v", name, err)
	}
	data, err := os.ReadFile(filepath.Join(out, name+".class"))
	if err != nil {
		t.Fatalf("read class: %v", err)
	}
	return data
}

func TestParseHelloWorld(t *testing.T) {
	cf, err := Parse(compileFixture(t, "HelloWorld"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := cf.ClassName(); got != "HelloWorld" {
		t.Errorf("ClassName = %q, want HelloWorld", got)
	}
	if got := cf.SuperClassName(); got != "java/lang/Object" {
		t.Errorf("SuperClassName = %q, want java/lang/Object", got)
	}

	// HelloWorld declares no fields and exactly 2 methods: <init>() and main.
	if len(cf.Fields()) != 0 {
		t.Errorf("Fields = %d, want 0", len(cf.Fields()))
	}
	if len(cf.Methods()) != 2 {
		t.Fatalf("Methods = %d, want 2", len(cf.Methods()))
	}

	var main *MemberInfo
	for _, m := range cf.Methods() {
		if m.Name() == "main" && m.Descriptor() == "([Ljava/lang/String;)V" {
			main = m
		}
	}
	if main == nil {
		t.Fatalf("main method not found")
	}
	code := main.Code()
	if code == nil {
		t.Fatalf("main has no Code attribute")
	}
	if len(code.Code()) == 0 {
		t.Errorf("main Code is empty")
	}
	if code.MaxLocals() < 3 { // args + a + b
		t.Errorf("MaxLocals = %d, want >= 3", code.MaxLocals())
	}
}

// TestParseStackMapTable verifies the delta-frame parser against javap's output:
// Fibonacci.fib has a single SAME frame at pc 7 (the if_icmpge target), with
// locals [int] and an empty stack.
func TestParseStackMapTable(t *testing.T) {
	cf, err := Parse(compileFixture(t, "Fibonacci"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var fib *MemberInfo
	for _, m := range cf.Methods() {
		if m.Name() == "fib" && m.Descriptor() == "(I)I" {
			fib = m
		}
	}
	if fib == nil {
		t.Fatal("fib not found")
	}
	smt := fib.Code().StackMapTable()
	if smt == nil {
		t.Fatal("fib has no StackMapTable")
	}
	// fib is static with one int arg, so the initial locals = [Integer].
	frames := smt.Reconstruct([]VerifType{{Tag: ItemInteger}})
	if len(frames) != 1 {
		t.Fatalf("got %d frames, want 1", len(frames))
	}
	if frames[0].Offset != 7 {
		t.Errorf("frame offset = %d, want 7", frames[0].Offset)
	}
	if len(frames[0].Locals) != 1 || frames[0].Locals[0].Tag != ItemInteger {
		t.Errorf("frame locals = %+v, want [Integer]", frames[0].Locals)
	}
	if len(frames[0].Stack) != 0 {
		t.Errorf("frame stack = %+v, want empty", frames[0].Stack)
	}
}

// --- Phase 1: FormatError and truncation safety ---

func TestBadMagic(t *testing.T) {
	data := compileFixture(t, "HelloWorld")
	data[0] = 0xDE // corrupt magic
	cf, err := Parse(data)
	if err == nil {
		t.Fatal("expected FormatError for bad magic")
	}
	if cf != nil {
		t.Error("Parse must return nil *ClassFile on error")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "magic" {
		t.Errorf("FormatError.Op = %q, want magic", fe.Op)
	}
}

func TestTruncatedEmpty(t *testing.T) {
	_, err := Parse([]byte{})
	if err == nil {
		t.Fatal("expected FormatError for empty input")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

func TestTruncatedMagic(t *testing.T) {
	_, err := Parse([]byte{0xCA, 0xFE}) // only 2 of 4 magic bytes
	if err == nil {
		t.Fatal("expected FormatError for truncated magic")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

func TestTruncatedAfterMagic(t *testing.T) {
	// A valid classfile header is at least 8 bytes (magic + version). Truncate
	// before constant_pool_count so the reader hits bounds.
	data := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00} // 7 bytes — valid magic but truncated version
	_, err := Parse(data)
	if err == nil {
		t.Fatal("expected FormatError for truncated version")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

func TestTruncatedConstantPool(t *testing.T) {
	// Valid magic + version, constant_pool_count claims 99 entries but the data
	// ends immediately.
	buf := make([]byte, 10)
	binary.BigEndian.PutUint32(buf[0:4], 0xCAFEBABE)
	binary.BigEndian.PutUint16(buf[4:6], 0)   // minor_version
	binary.BigEndian.PutUint16(buf[6:8], 69)  // major_version (JDK 25)
	binary.BigEndian.PutUint16(buf[8:10], 99) // constant_pool_count = 99, but no entries follow
	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for truncated constant pool")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// --- Exact consumption and deterministic truncation tests ---

// TestBootstrapMethodsTrailingBytes verifies that trailing bytes after the
// BootstrapMethods attribute body are rejected.
func TestBootstrapMethodsTrailingBytes(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
	}
	// Valid body: 1 entry, no args = 6 bytes. Then append extra garbage.
	body := []byte{0x00, 0x01, 0x00, 0x07, 0x00, 0x00} // count=1, methodRef=7, argc=0
	body = append(body, 0xFF, 0xFF, 0xFF)              // trailing bytes
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 8) // name_index = "BootstrapMethods"
	attr = binary.BigEndian.AppendUint32(attr, uint32(len(body)))
	attr = append(attr, body...)
	attrs := [][]byte{attr}

	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for trailing bytes in BootstrapMethods")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "BootstrapMethods" {
		t.Errorf("FormatError.Op = %q, want BootstrapMethods", fe.Op)
	}
}

// TestCodeAttributeTrailingBytes verifies that Code rejects trailing bytes.
func TestCodeAttributeTrailingBytes(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildUTF8("Code"),                       // [5]
	}
	// Valid Code body: maxStack(2), maxLocals(2), codeLen(4), code(4), excLen(2), attrCnt(2)
	// Add trailing bytes.
	codeBody := []byte{
		0x00, 0x01, // max_stack
		0x00, 0x01, // max_locals
		0x00, 0x00, 0x00, 0x00, // code_length=0
		0x00, 0x00, // exception_table_length=0
		0x00, 0x00, // attributes_count=0
		0xAA, 0xBB, // trailing garbage
	}
	codeAttr := []byte{}
	codeAttr = binary.BigEndian.AppendUint16(codeAttr, 5) // name_index = "Code"
	codeAttr = binary.BigEndian.AppendUint32(codeAttr, uint32(len(codeBody)))
	codeAttr = append(codeAttr, codeBody...)

	methodAttrs := binary.BigEndian.AppendUint16(nil, 1)
	methodAttrs = append(methodAttrs, codeAttr...)
	method := buildMemberBytes(0x0009, 3, 4, methodAttrs)

	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001)
	buf = binary.BigEndian.AppendUint16(buf, 2)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 1)
	buf = append(buf, method...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)

	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for trailing bytes in Code attribute")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "Code" {
		t.Errorf("FormatError.Op = %q, want Code", fe.Op)
	}
}

// TestTruncatedBootstrapMethods verifies a BootstrapMethods attribute whose
// declared body length exceeds available data at a deterministic cut point.
func TestTruncatedBootstrapMethods(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("BootstrapMethods"),           // [3]
	}
	// Declared length 50 bytes, only 2 bytes provided.
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 3)  // name_index = "BootstrapMethods"
	attr = binary.BigEndian.AppendUint32(attr, 50) // declared length
	attr = append(attr, 0x00, 0x00)                // only 2 bytes of body
	attrs := [][]byte{attr}

	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if cf != nil {
		t.Error("Parse must return nil *ClassFile on error")
	}
	if err == nil {
		t.Fatal("expected FormatError for truncated BootstrapMethods")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestTruncatedStackMapTable verifies a StackMapTable attribute whose declared
// body length exceeds available data.
func TestTruncatedStackMapTable(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildUTF8("Code"),                       // [5]
		buildUTF8("StackMapTable"),              // [6]
	}
	// StackMapTable body with declared length 99, only 4 bytes provided.
	smtBody := []byte{0x00, 0x00, 0xAA, 0xBB}
	smtAttr := []byte{}
	smtAttr = binary.BigEndian.AppendUint16(smtAttr, 6) // "StackMapTable"
	smtAttr = binary.BigEndian.AppendUint32(smtAttr, 99)
	smtAttr = append(smtAttr, smtBody...)

	codeAttrs := binary.BigEndian.AppendUint16(nil, 1)
	codeAttrs = append(codeAttrs, smtAttr...)

	codeBody := []byte{
		0x00, 0x01, // max_stack
		0x00, 0x01, // max_locals
		0x00, 0x00, 0x00, 0x00, // code_length=0
		0x00, 0x00, // exception_table_length=0
	}
	codeBody = append(codeBody, codeAttrs...)

	codeAttr := []byte{}
	codeAttr = binary.BigEndian.AppendUint16(codeAttr, 5) // "Code"
	codeAttr = binary.BigEndian.AppendUint32(codeAttr, uint32(len(codeBody)))
	codeAttr = append(codeAttr, codeBody...)

	methodAttrs := binary.BigEndian.AppendUint16(nil, 1)
	methodAttrs = append(methodAttrs, codeAttr...)
	method := buildMemberBytes(0x0009, 3, 4, methodAttrs)

	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001)
	buf = binary.BigEndian.AppendUint16(buf, 2)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 1)
	buf = append(buf, method...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)

	cf, err := Parse(buf)
	if cf != nil {
		t.Error("Parse must return nil *ClassFile on error")
	}
	if err == nil {
		t.Fatal("expected FormatError for truncated StackMapTable")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

func TestUnknownConstantTag(t *testing.T) {
	// Valid magic + version, one entry in the pool with tag 99 (invalid).
	buf := make([]byte, 11)
	binary.BigEndian.PutUint32(buf[0:4], 0xCAFEBABE)
	binary.BigEndian.PutUint16(buf[4:6], 0)  // minor_version
	binary.BigEndian.PutUint16(buf[6:8], 69) // major_version
	binary.BigEndian.PutUint16(buf[8:10], 2) // constant_pool_count = 2 (index 1 occupied)
	buf[10] = 99                             // invalid tag
	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for unknown constant tag")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "constant pool" {
		t.Errorf("FormatError.Op = %q, want constant pool", fe.Op)
	}
}

// TestNoBoundsPanicOnTruncatedAttribute verifies that a truncated attribute
// body does not escape as a slice-bounds panic. The test builds a classfile
// with a precisely crafted Code attribute whose declared length exceeds the
// available data — no arbitrary percentage is used.
func TestNoBoundsPanicOnTruncatedAttribute(t *testing.T) {
	// Build a classfile where the last attribute (a Code attribute on the sole
	// method) declares a body length larger than the remaining bytes.
	// Structure:
	//   [1] CONSTANT_Utf8 "Code"
	//   [2] CONSTANT_Utf8 "C"
	//   [3] CONSTANT_Class -> 2
	//   [4] CONSTANT_Utf8 "m"
	//   [5] CONSTANT_Utf8 "()V"
	//   method: name=4, desc=5, attrs=[Code(name_idx=1, length=999)]
	// The attribute length declares 999 but we only provide 2 bytes — ReadSlice
	// must detect the bounds violation and panic with parsePanic.

	nameIdx := uint16(1)       // "Code"
	classNameIdx := uint16(3)  // CONSTANT_Class
	methodNameIdx := uint16(4) // "m"
	methodDescIdx := uint16(5) // "()V"

	cp := [][]byte{
		buildUTF8("Code"),                       // [1]
		buildUTF8("C"),                          // [2]
		buildCPEntry(ConstantClass, 0x00, 0x02), // [3]
		buildUTF8("m"),                          // [4]
		buildUTF8("()V"),                        // [5]
	}

	// Attribute: name_index=1, length=999 (far exceeding actual data), body=2 bytes
	codeAttr := []byte{}
	codeAttr = binary.BigEndian.AppendUint16(codeAttr, nameIdx)
	codeAttr = binary.BigEndian.AppendUint32(codeAttr, 999) // declared length — huge lie
	codeAttr = append(codeAttr, 0x00, 0x00)                 // only 2 bytes of body

	methodAttrs := binary.BigEndian.AppendUint16(nil, 1)
	methodAttrs = append(methodAttrs, codeAttr...)
	method := buildMemberBytes(0x0009, methodNameIdx, methodDescIdx, methodAttrs)

	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001) // access_flags
	buf = binary.BigEndian.AppendUint16(buf, classNameIdx)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super_class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces_count
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // fields_count
	buf = binary.BigEndian.AppendUint16(buf, 1)      // methods_count = 1
	buf = append(buf, method...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // attributes_count

	cf, err := Parse(buf)
	if cf != nil {
		t.Error("Parse must return nil *ClassFile on error")
	}
	if err == nil {
		t.Fatal("expected FormatError from truncated attribute")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestParseValidFixtures verifies existing positive coverage still works.
func TestParseValidFixtures(t *testing.T) {
	for _, name := range []string{"HelloWorld", "Fibonacci"} {
		t.Run(name, func(t *testing.T) {
			cf, err := Parse(compileFixture(t, name))
			if err != nil {
				t.Fatalf("Parse(%s): %v", name, err)
			}
			if cf.ClassName() == "" {
				t.Errorf("%s: empty class name", name)
			}
		})
	}
}

// --- Phase 3: Typed dynamic constant-pool tests ---

// readPool parses a raw constant_pool from bytes (without the surrounding
// classfile structure). Used for unit-testing pool accessors.
func readPool(t *testing.T, data []byte) *ConstantPool {
	t.Helper()
	r := NewClassReader(data)
	cp := readConstantPool(r)
	if r.Len() > 0 {
		t.Fatalf("unconsumed reader bytes: %d", r.Len())
	}
	return cp
}

// buildCPEntry builds raw bytes for one constant pool entry.
func buildCPEntry(tag uint8, payload ...byte) []byte {
	b := []byte{tag}
	return append(b, payload...)
}

// buildUTF8 builds a CONSTANT_Utf8 entry (tag 1, u2 length, bytes).
func buildUTF8(s string) []byte {
	b := []byte{ConstantUtf8}
	b = binary.BigEndian.AppendUint16(b, uint16(len(s)))
	return append(b, []byte(s)...)
}

// buildMHFixture builds a constant pool suitable for testing a MethodHandle
// entry at the given pool index. Returns (pool bytes, mhIndex).
//
// Layout:
//
//	[1] CONSTANT_Utf8 memberName
//	[2] CONSTANT_Utf8 memberDesc
//	[3] CONSTANT_Utf8 className
//	[4] CONSTANT_Class -> 3
//	[5] CONSTANT_NameAndType -> (1, 2)
//	[6] CONSTANT_MemberRef (memberTag) -> (4, 5)
//	[7] CONSTANT_MethodHandle (kind, ref=6)
//	... optional extras from appendEntries
func buildMHFixture(kind uint8, memberTag uint8, memberName, memberDesc, className string, appendEntries [][]byte) ([]byte, uint16) {
	entries := [][]byte{
		buildUTF8(memberName),                                     // [1]
		buildUTF8(memberDesc),                                     // [2]
		buildUTF8(className),                                      // [3]
		buildCPEntry(ConstantClass, 0x00, 0x03),                   // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02), // [5]
		buildCPEntry(memberTag, 0x00, 0x04, 0x00, 0x05),           // [6]
		buildCPEntry(ConstantMethodHandle, kind, 0x00, 0x06),      // [7]
	}
	entries = append(entries, appendEntries...)
	mhIndex := uint16(7)
	buf := []byte{}
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)+1)) // cp_count
	for _, e := range entries {
		buf = append(buf, e...)
	}
	return buf, mhIndex
}

// TestMHTableDrivenPositive tests all 9 MethodHandle reference kinds with
// valid targets, plus version-52 InterfaceMethodref acceptance.
func TestMHTableDrivenPositive(t *testing.T) {
	tests := []struct {
		name         string
		kind         uint8
		memberTag    uint8
		memberName   string
		memberDesc   string
		className    string
		majorVersion uint16
		wantName     string
		wantDesc     string
		wantClass    string
		wantRefTag   uint8
	}{
		// Fields (kinds 1–4)
		{"kind1_getField", RefGetField, ConstantFieldref, "f", "I", "MyClass", 69, "f", "I", "MyClass", ConstantFieldref},
		{"kind2_getStatic", RefGetStatic, ConstantFieldref, "f", "I", "MyClass", 69, "f", "I", "MyClass", ConstantFieldref},
		{"kind3_putField", RefPutField, ConstantFieldref, "f", "I", "MyClass", 69, "f", "I", "MyClass", ConstantFieldref},
		{"kind4_putStatic", RefPutStatic, ConstantFieldref, "f", "I", "MyClass", 69, "f", "I", "MyClass", ConstantFieldref},
		// Methods (kinds 5–9)
		{"kind5_invokeVirtual", RefInvokeVirtual, ConstantMethodref, "m", "()I", "MyClass", 69, "m", "()I", "MyClass", ConstantMethodref},
		{"kind6_invokeStatic_methodref", RefInvokeStatic, ConstantMethodref, "m", "()V", "MyClass", 69, "m", "()V", "MyClass", ConstantMethodref},
		{"kind7_invokeSpecial_methodref", RefInvokeSpecial, ConstantMethodref, "m", "(I)I", "MyClass", 69, "m", "(I)I", "MyClass", ConstantMethodref},
		{"kind8_newInvokeSpecial", RefNewInvokeSpecial, ConstantMethodref, "<init>", "()V", "MyClass", 69, "<init>", "()V", "MyClass", ConstantMethodref},
		{"kind9_invokeInterface", RefInvokeInterface, ConstantInterfaceMethodref, "m", "()I", "MyClass", 69, "m", "()I", "MyClass", ConstantInterfaceMethodref},
		// InterfaceMethodref at major >= 52 for kinds 6/7
		{"kind6_interface_v52", RefInvokeStatic, ConstantInterfaceMethodref, "m", "()V", "MyClass", 52, "m", "()V", "MyClass", ConstantInterfaceMethodref},
		{"kind7_interface_v52", RefInvokeSpecial, ConstantInterfaceMethodref, "m", "()V", "MyClass", 52, "m", "()V", "MyClass", ConstantInterfaceMethodref},
		// Kind 6/7 with Methodref at major 51 — always valid
		{"kind6_methodref_v51", RefInvokeStatic, ConstantMethodref, "m", "()V", "MyClass", 51, "m", "()V", "MyClass", ConstantMethodref},
		{"kind7_methodref_v51", RefInvokeSpecial, ConstantMethodref, "m", "()V", "MyClass", 51, "m", "()V", "MyClass", ConstantMethodref},
		// Kind 8 with parameterized constructor
		{"kind8_with_params", RefNewInvokeSpecial, ConstantMethodref, "<init>", "(Ljava/lang/String;)V", "MyClass", 69, "<init>", "(Ljava/lang/String;)V", "MyClass", ConstantMethodref},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, mhIdx := buildMHFixture(tt.kind, tt.memberTag, tt.memberName, tt.memberDesc, tt.className, nil)
			cp := readPool(t, buf)

			// Typed accessor path.
			rk, rt, cls, name, desc, err := cp.MethodHandleReference(mhIdx, tt.majorVersion)
			if err != nil {
				t.Fatalf("MethodHandleReference: %v", err)
			}
			if rk != tt.kind {
				t.Errorf("refKind = %d, want %d", rk, tt.kind)
			}
			if rt != tt.wantRefTag {
				t.Errorf("refTag = %d, want %d", rt, tt.wantRefTag)
			}
			if cls != tt.wantClass {
				t.Errorf("class = %q, want %q", cls, tt.wantClass)
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
			if desc != tt.wantDesc {
				t.Errorf("desc = %q, want %q", desc, tt.wantDesc)
			}

			// RefKind accessor must also work.
			kind, err := cp.MethodHandleRefKind(mhIdx)
			if err != nil {
				t.Fatalf("MethodHandleRefKind: %v", err)
			}
			if kind != tt.kind {
				t.Errorf("MethodHandleRefKind = %d, want %d", kind, tt.kind)
			}
		})
	}
}

// TestMHFieldAcceptsInitClinit verifies that MethodHandle kinds 1–4 (field
// references) accept <init> and <clinit> as field names. JVMS field names are
// unqualified names with no extra restriction on angle brackets — JDK 25
// accepts classfiles with <init>-named fields referenced by REF_getField.
func TestMHFieldAcceptsInitClinit(t *testing.T) {
	tests := []struct {
		name      string
		kind      uint8
		fieldName string
	}{
		{"kind1_init", RefGetField, "<init>"},
		{"kind1_clinit", RefGetField, "<clinit>"},
		{"kind2_init", RefGetStatic, "<init>"},
		{"kind2_clinit", RefGetStatic, "<clinit>"},
		{"kind3_init", RefPutField, "<init>"},
		{"kind3_clinit", RefPutField, "<clinit>"},
		{"kind4_init", RefPutStatic, "<init>"},
		{"kind4_clinit", RefPutStatic, "<clinit>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolBytes, _ := buildMHFixture(tt.kind, ConstantFieldref, tt.fieldName, "I", "C", nil)
			buf := buildMHClassfile(poolBytes, 69, nil)
			cf, err := Parse(buf)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if cf == nil {
				t.Fatal("Parse returned nil")
			}
			// Also verify typed accessor.
			cp := cf.ConstantPool()
			rk, rt, _, name, desc, err := cp.MethodHandleReference(7, 69)
			if err != nil {
				t.Fatalf("MethodHandleReference: %v", err)
			}
			if rk != tt.kind {
				t.Errorf("refKind = %d, want %d", rk, tt.kind)
			}
			if rt != ConstantFieldref {
				t.Errorf("refTag = %d, want ConstantFieldref", rt)
			}
			if name != tt.fieldName {
				t.Errorf("name = %q, want %q", name, tt.fieldName)
			}
			if desc != "I" {
				t.Errorf("desc = %q, want I", desc)
			}
		})
	}
}

// assertParseFormatError is a helper that parses buf and asserts: cf==nil,
// err is *FormatError with the expected Op. Fails the test on any deviation.
func assertParseFormatError(t *testing.T, buf []byte, wantOp string) {
	t.Helper()
	cf, err := Parse(buf)
	if cf != nil {
		t.Error("Parse must return nil *ClassFile on error")
	}
	if err == nil {
		t.Fatal("expected *FormatError, got nil")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != wantOp {
		t.Errorf("FormatError.Op = %q, want %q", fe.Op, wantOp)
	}
}

// buildMHClassfile wraps a MethodHandle test pool into a minimal classfile
// suitable for Parse validation. The pool must have a CONSTANT_Class at index 2
// for this_class, or a separate Class entry if provided.
func buildMHClassfile(poolBytes []byte, majorVersion uint16, extraAttrs [][]byte) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, majorVersion)
	buf = append(buf, poolBytes...)
	// Access flags, this_class, super, interfaces, fields, methods
	buf = binary.BigEndian.AppendUint16(buf, 0x0001)
	buf = binary.BigEndian.AppendUint16(buf, 0x0004) // this_class -> index 4 (the CONSTANT_Class)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super_class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // fields
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // methods
	buf = append(buf, buildClassAttr(extraAttrs...)...)
	return buf
}

// buildCPBuf builds a constant pool byte slice from a list of entry byte
// slices. Used by nested-index-error tests that need custom pool layouts
// beyond what buildMHFixture provides.
func buildCPBuf(entries [][]byte) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)+1))
	for _, e := range entries {
		buf = append(buf, e...)
	}
	return buf
}

// buildCPBufWithCount is like buildCPBuf but accepts an explicit cp_count.
// Used when the pool contains Long/Double entries whose second slot must be
// counted in cp_count.
func buildCPBufWithCount(entries [][]byte, cpCount uint16) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint16(buf, cpCount)
	for _, e := range entries {
		buf = append(buf, e...)
	}
	return buf
}

// buildMHClassfileNested is like buildMHClassfile but accepts an explicit
// this_class index for pools where the Class entry is not at index 4.
func buildMHClassfileNested(poolBytes []byte, thisClassIndex uint16, majorVersion uint16) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, majorVersion)
	buf = append(buf, poolBytes...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0001)
	buf = binary.BigEndian.AppendUint16(buf, thisClassIndex)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // fields
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // methods
	buf = append(buf, buildClassAttr(nil...)...)
	return buf
}

// TestMHTableDrivenNegativeParse tests that Parse rejects malformed
// MethodHandle entries with *FormatError (cf==nil, no non-parsePanic panic).
func TestMHTableDrivenNegativeParse(t *testing.T) {
	tests := []struct {
		name         string
		kind         uint8
		memberTag    uint8
		memberName   string
		memberDesc   string
		className    string
		majorVersion uint16
		wantOp       string // expected FormatError.Op
	}{
		// --- Wrong target tag per kind ---
		{"kind1_wrongTag_methodref", RefGetField, ConstantMethodref, "f", "I", "C", 69, "constant pool"},
		{"kind2_wrongTag_methodref", RefGetStatic, ConstantMethodref, "f", "I", "C", 69, "constant pool"},
		{"kind3_wrongTag_imethodref", RefPutField, ConstantInterfaceMethodref, "f", "I", "C", 69, "constant pool"},
		{"kind4_wrongTag_imethodref", RefPutStatic, ConstantInterfaceMethodref, "f", "I", "C", 69, "constant pool"},
		{"kind5_wrongTag_fieldref", RefInvokeVirtual, ConstantFieldref, "m", "()V", "C", 69, "constant pool"},
		{"kind6_wrongTag_fieldref", RefInvokeStatic, ConstantFieldref, "m", "()V", "C", 69, "constant pool"},
		{"kind7_wrongTag_fieldref", RefInvokeSpecial, ConstantFieldref, "m", "()V", "C", 69, "constant pool"},
		{"kind8_wrongTag_fieldref", RefNewInvokeSpecial, ConstantFieldref, "<init>", "()V", "C", 69, "constant pool"},
		{"kind9_wrongTag_methodref", RefInvokeInterface, ConstantMethodref, "m", "()V", "C", 69, "constant pool"},

		// --- Kind 6/7 InterfaceMethodref rejected at major < 52 ---
		{"kind6_interface_v51", RefInvokeStatic, ConstantInterfaceMethodref, "m", "()V", "C", 51, "constant pool"},
		{"kind7_interface_v51", RefInvokeSpecial, ConstantInterfaceMethodref, "m", "()V", "C", 51, "constant pool"},

		// --- <init>/<clinit> rejected for kinds 5/6/7/9 ---
		{"kind5_init_rejected", RefInvokeVirtual, ConstantMethodref, "<init>", "()V", "C", 69, "constant pool"},
		{"kind5_clinit_rejected", RefInvokeVirtual, ConstantMethodref, "<clinit>", "()V", "C", 69, "constant pool"},
		{"kind6_init_rejected", RefInvokeStatic, ConstantMethodref, "<init>", "()V", "C", 69, "constant pool"},
		{"kind6_clinit_rejected", RefInvokeStatic, ConstantMethodref, "<clinit>", "()V", "C", 69, "constant pool"},
		{"kind7_init_rejected", RefInvokeSpecial, ConstantMethodref, "<init>", "()V", "C", 69, "constant pool"},
		{"kind7_clinit_rejected", RefInvokeSpecial, ConstantMethodref, "<clinit>", "()V", "C", 69, "constant pool"},
		{"kind9_init_rejected", RefInvokeInterface, ConstantInterfaceMethodref, "<init>", "()V", "C", 69, "constant pool"},
		{"kind9_clinit_rejected", RefInvokeInterface, ConstantInterfaceMethodref, "<clinit>", "()V", "C", 69, "constant pool"},

		// --- Kind 8: wrong name ---
		{"kind8_not_init", RefNewInvokeSpecial, ConstantMethodref, "m", "()V", "C", 69, "constant pool"},
		{"kind8_is_clinit", RefNewInvokeSpecial, ConstantMethodref, "<clinit>", "()V", "C", 69, "constant pool"},

		// --- Kind 8: field descriptor rejected ---
		{"kind8_field_desc_I", RefNewInvokeSpecial, ConstantMethodref, "<init>", "I", "C", 69, "constant pool"},

		// --- Kind 8: non-V return rejected ---
		{"kind8_return_I", RefNewInvokeSpecial, ConstantMethodref, "<init>", "()I", "C", 69, "constant pool"},
		{"kind8_return_object", RefNewInvokeSpecial, ConstantMethodref, "<init>", "(I)I", "C", 69, "constant pool"},

		// --- Method kind with field descriptor ---
		{"kind5_field_desc", RefInvokeVirtual, ConstantMethodref, "m", "I", "C", 69, "constant pool"},

		// --- Field kind with method descriptor ---
		{"kind1_method_desc", RefGetField, ConstantFieldref, "f", "()V", "C", 69, "constant pool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolBytes, _ := buildMHFixture(tt.kind, tt.memberTag, tt.memberName, tt.memberDesc, tt.className, nil)
			buf := buildMHClassfile(poolBytes, tt.majorVersion, nil)
			assertParseFormatError(t, buf, tt.wantOp)
		})
	}
}

// TestMHMethodRejectsAngleBracketNames verifies that MethodHandle kinds 5/6/7/9
// reject names containing '<' or '>' (per JVMS §4.2.2 ordinary method name
// rules). This catches <foo>, <init>, <clinit>, and any other non-method-name.
func TestMHMethodRejectsAngleBracketNames(t *testing.T) {
	tests := []struct {
		name      string
		kind      uint8
		memberTag uint8
		badName   string
	}{
		{"kind5_foo", RefInvokeVirtual, ConstantMethodref, "<foo>"},
		{"kind5_bar", RefInvokeVirtual, ConstantMethodref, "<bar>"},
		{"kind5_empty", RefInvokeVirtual, ConstantMethodref, "<>"},
		{"kind6_foo", RefInvokeStatic, ConstantMethodref, "<foo>"},
		{"kind7_foo", RefInvokeSpecial, ConstantMethodref, "<foo>"},
		{"kind9_foo", RefInvokeInterface, ConstantInterfaceMethodref, "<foo>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolBytes, _ := buildMHFixture(tt.kind, tt.memberTag, tt.badName, "()V", "C", nil)
			buf := buildMHClassfile(poolBytes, 69, nil)
			assertParseFormatError(t, buf, "constant pool")
		})
	}
}

// TestMHNestedIndexErrors tests that every nested CP index in a MethodHandle
// chain is validated. Both the Parse path (via assertParseFormatError) and the
// typed accessor path (via cp.MethodHandleReference) are exercised.
func TestMHNestedIndexErrors(t *testing.T) {
	// --- Typed accessor path: reference_index errors ---

	t.Run("reference_index_zero", func(t *testing.T) {
		entries := [][]byte{
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x00), // [1]: kind=5, ref=0
		}
		buf := []byte{}
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)+1))
		for _, e := range entries {
			buf = append(buf, e...)
		}
		cp := readPool(t, buf)
		_, _, _, _, _, err := cp.MethodHandleReference(1, 69)
		if err == nil {
			t.Fatal("expected error for reference_index 0")
		}
	})

	t.Run("reference_index_OOB", func(t *testing.T) {
		entries := [][]byte{
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x63), // [1]: kind=5, ref=99 (OOB)
		}
		buf := []byte{}
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)+1))
		for _, e := range entries {
			buf = append(buf, e...)
		}
		cp := readPool(t, buf)
		_, _, _, _, _, err := cp.MethodHandleReference(1, 69)
		if err == nil {
			t.Fatal("expected error for reference_index OOB")
		}
	})

	t.Run("reference_index_second_slot", func(t *testing.T) {
		// Pool: [1] Long (takes slots 1-2), [3] MethodHandle kind=5 ref=2.
		// cp_count=4 accounts for the Long's second slot.
		buf := []byte{}
		buf = binary.BigEndian.AppendUint16(buf, 4)                                                      // cp_count=4 (indices 0–3)
		buf = append(buf, buildCPEntry(ConstantLong, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42)...) // [1]
		buf = append(buf, buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x02)...)           // [3]: ref=2 (unusable)
		cp := readPool(t, buf)
		_, _, _, _, _, err := cp.MethodHandleReference(3, 69)
		if err == nil {
			t.Fatal("expected error for reference_index pointing to second slot")
		}
	})

	// --- Parse path: nested index errors (each returns *FormatError, cf==nil) ---

	t.Run("parse_ref_index_OOB", func(t *testing.T) {
		// MH at index 7 with refIndex=99 (OOB for a 7-entry pool).
		pool := [][]byte{
			buildUTF8("m"),                          // [1]
			buildUTF8("()V"),                        // [2]
			buildUTF8("C"),                          // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4]
			buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02),        // [5]
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x05),          // [6]
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x63), // [7]: ref=99
		}
		buf := buildMHClassfile(buildCPBuf(pool), 69, nil)
		assertParseFormatError(t, buf, "constant pool")
	})

	t.Run("parse_ref_index_second_slot", func(t *testing.T) {
		// [1] Long (takes 1-2), [3] Utf8 "C", [4] Class→3, [5] Utf8 "m",
		// [6] Utf8 "()V", [7] NameAndType→(5,6), [8] Methodref→(4,7),
		// [9] MH kind=5 ref=2. cp_count=10 accounts for the Long double-slot.
		pool := [][]byte{
			buildCPEntry(ConstantLong, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42), // [1]
			buildUTF8("C"),                          // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4]
			buildUTF8("m"),                          // [5]
			buildUTF8("()V"),                        // [6]
			buildCPEntry(ConstantNameAndType, 0x00, 0x05, 0x00, 0x06),        // [7]
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x07),          // [8]
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x02), // [9]: ref=2
		}
		poolBytes := buildCPBufWithCount(pool, 10) // Long occupies two slots
		buf := buildMHClassfileNested(poolBytes, 4, 69)
		assertParseFormatError(t, buf, "constant pool")
	})

	t.Run("parse_class_index_wrong_tag", func(t *testing.T) {
		// MemberRef.classIndex → Utf8 instead of Class.
		// [1] Utf8 "m", [2] Utf8 "()V", [3] Utf8 "C" (WRONG — used as classIndex, should be Class)
		// [4] Class→3 (for this_class), [5] NameAndType→(1,2),
		// [6] Methodref→(3,5) ← classIndex=3 is Utf8!
		// [7] MH kind=5 ref=6
		pool := [][]byte{
			buildUTF8("m"),                          // [1]
			buildUTF8("()V"),                        // [2]
			buildUTF8("C"),                          // [3] Utf8 (not Class!)
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4] Class→3 (for this_class)
			buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02),        // [5]
			buildCPEntry(ConstantMethodref, 0x00, 0x03, 0x00, 0x05),          // [6] classIndex=3→Utf8
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x06), // [7]
		}
		buf := buildMHClassfile(buildCPBuf(pool), 69, nil)
		assertParseFormatError(t, buf, "constant pool")
	})

	t.Run("parse_nat_index_wrong_tag", func(t *testing.T) {
		// MemberRef.natIndex → Utf8 instead of NameAndType.
		// [1] Utf8 "m", [2] Utf8 "()V", [3] Utf8 "C",
		// [4] Class→3, [5] Utf8 "not-a-nat" (WRONG — used as natIndex),
		// [6] Methodref→(4,5) ← natIndex=5 is Utf8!
		// [7] MH kind=5 ref=6
		pool := [][]byte{
			buildUTF8("m"),                          // [1]
			buildUTF8("()V"),                        // [2]
			buildUTF8("C"),                          // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4]
			buildUTF8("not-a-nat"),                  // [5] Utf8 (not NameAndType!)
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x05),          // [6] natIndex=5→Utf8
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x06), // [7]
		}
		buf := buildMHClassfile(buildCPBuf(pool), 69, nil)
		assertParseFormatError(t, buf, "constant pool")
	})

	t.Run("parse_name_index_wrong_tag", func(t *testing.T) {
		// NameAndType.nameIndex → non-Utf8 (e.g., Integer).
		// [1] Utf8 "m", [2] Utf8 "()V", [3] Utf8 "C",
		// [4] Class→3, [5] Integer 42,
		// [6] NameAndType→(5,2) ← nameIndex=5 is Integer!
		// [7] Methodref→(4,6), [8] MH kind=5 ref=7
		pool := [][]byte{
			buildUTF8("m"),                          // [1]
			buildUTF8("()V"),                        // [2]
			buildUTF8("C"),                          // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4]
			buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x2A),            // [5] Integer (not Utf8!)
			buildCPEntry(ConstantNameAndType, 0x00, 0x05, 0x00, 0x02),        // [6] nameIndex=5→Integer
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x06),          // [7]
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x07), // [8]
		}
		buf := buildMHClassfileNested(buildCPBuf(pool), 4, 69)
		assertParseFormatError(t, buf, "constant pool")
	})

	t.Run("parse_descriptor_index_wrong_tag", func(t *testing.T) {
		// NameAndType.descIndex → non-Utf8 (e.g., Integer).
		// [1] Utf8 "m", [2] Integer 42, [3] Utf8 "C",
		// [4] Class→3, [5] NameAndType→(1,2) ← descIndex=2 is Integer!
		// [6] Methodref→(4,5), [7] MH kind=5 ref=6
		pool := [][]byte{
			buildUTF8("m"), // [1]
			buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x2A), // [2] Integer (not Utf8!)
			buildUTF8("C"),                                                   // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03),                          // [4]
			buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02),        // [5] descIndex=2→Integer
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x05),          // [6]
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x06), // [7]
		}
		buf := buildMHClassfile(buildCPBuf(pool), 69, nil)
		assertParseFormatError(t, buf, "constant pool")
	})

	// --- Typed accessor path: nested wrong-tag errors (no Parse wrapper) ---

	t.Run("typed_class_index_wrong_tag", func(t *testing.T) {
		// Same structure as parse_class_index_wrong_tag but via typed accessor.
		entries := [][]byte{
			buildUTF8("m"),   // [1]
			buildUTF8("()V"), // [2]
			buildUTF8("C"),   // [3] Utf8 (not Class!)
			buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02),        // [4]
			buildCPEntry(ConstantMethodref, 0x00, 0x03, 0x00, 0x04),          // [5] classIndex=3→Utf8
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x05), // [6]
		}
		cp := readPool(t, buildCPBuf(entries))
		_, _, _, _, _, err := cp.MethodHandleReference(6, 69)
		if err == nil {
			t.Fatal("expected error for class_index pointing to Utf8")
		}
	})

	t.Run("typed_nat_index_wrong_tag", func(t *testing.T) {
		// MemberRef.natIndex → Utf8.
		entries := [][]byte{
			buildUTF8("m"),                          // [1]
			buildUTF8("()V"),                        // [2]
			buildUTF8("C"),                          // [3]
			buildCPEntry(ConstantClass, 0x00, 0x03), // [4] Class→3
			buildUTF8("not-a-nat"),                  // [5] Utf8 (not NameAndType!)
			buildCPEntry(ConstantMethodref, 0x00, 0x04, 0x00, 0x05),          // [6] natIndex=5→Utf8
			buildCPEntry(ConstantMethodHandle, RefInvokeVirtual, 0x00, 0x06), // [7]
		}
		cp := readPool(t, buildCPBuf(entries))
		_, _, _, _, _, err := cp.MethodHandleReference(7, 69)
		if err == nil {
			t.Fatal("expected error for nat_index pointing to Utf8")
		}
	})
}

// TestMHNegativeTypedAccessor tests that the typed accessor returns error
// (never panics) for malformed MH entries.
func TestMHNegativeTypedAccessor(t *testing.T) {
	t.Run("kind0_invalid", func(t *testing.T) {
		entries := [][]byte{
			buildCPEntry(ConstantMethodHandle, 0x00, 0x00, 0x00), // kind=0, ref=0
		}
		buf := []byte{}
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)+1))
		for _, e := range entries {
			buf = append(buf, e...)
		}
		cp := readPool(t, buf)
		_, err := cp.MethodHandleRefKind(1)
		if err == nil {
			t.Fatal("expected error for invalid ref kind 0")
		}
		// Must not panic.
		_, _, _, _, _, err = cp.MethodHandleReference(1, 69)
		if err == nil {
			t.Fatal("expected error from MethodHandleReference for invalid kind 0")
		}
	})

	t.Run("index_OOB", func(t *testing.T) {
		buf := []byte{0x00, 0x02} // cp_count=2, only index 1 possible
		buf = append(buf, buildUTF8("x")...)
		cp := readPool(t, buf)
		_, err := cp.MethodHandleRefKind(99)
		if err == nil {
			t.Fatal("expected error for OOB index")
		}
	})

	t.Run("index_zero", func(t *testing.T) {
		buf := []byte{0x00, 0x02}
		buf = append(buf, buildUTF8("x")...)
		cp := readPool(t, buf)
		_, err := cp.MethodHandleRefKind(0)
		if err == nil {
			t.Fatal("expected error for index 0")
		}
	})

	t.Run("second_slot", func(t *testing.T) {
		// Build a pool where index 2 is the unusable second slot.
		entries := [][]byte{
			buildCPEntry(ConstantLong, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42), // [1] long (takes 1-2)
			buildUTF8("x"), // [3]
		}
		buf := []byte{0x00, 0x04}
		for _, e := range entries {
			buf = append(buf, e...)
		}
		cp := readPool(t, buf)
		_, err := cp.MethodHandleRefKind(2)
		if err == nil {
			t.Fatal("expected error for second slot")
		}
	})
}

// TestMHKind8DescriptorRejection tests kind 8 descriptor edge cases via Parse.
func TestMHKind8DescriptorRejection(t *testing.T) {
	tests := []struct {
		name string
		desc string
	}{
		{"init_return_I", "<init>:I"}, // bare field type, not a method descriptor
		{"init_return_nonV", "()I"},
		{"init_params_return_nonV", "(I)I"},
		{"init_void_param", "(V)V"},
		{"init_return_Q", "()Q"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build pool with kind 8 MH pointing at a Methodref with the bad desc.
			poolBytes, _ := buildMHFixture(RefNewInvokeSpecial, ConstantMethodref, "<init>", tt.desc, "C", nil)
			buf := buildMHClassfile(poolBytes, 69, nil)
			assertParseFormatError(t, buf, "constant pool")
		})
	}
}

// TestMHKind8DescriptorAcceptance tests kind 8 descriptor positive cases via Parse.
func TestMHKind8DescriptorAcceptance(t *testing.T) {
	tests := []struct {
		name string
		desc string
	}{
		{"init_no_params", "()V"},
		{"init_one_int", "(I)V"},
		{"init_one_ref", "(Ljava/lang/String;)V"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolBytes, _ := buildMHFixture(RefNewInvokeSpecial, ConstantMethodref, "<init>", tt.desc, "C", nil)
			buf := buildMHClassfile(poolBytes, 69, nil)
			cf, err := Parse(buf)
			if err != nil {
				t.Fatalf("Parse for valid kind 8 %s: %v", tt.name, err)
			}
			if cf == nil {
				t.Fatal("Parse returned nil")
			}
		})
	}
}

// TestMethodTypeDescriptor tests MethodTypeDescriptor with checked lookup.
func TestMethodTypeDescriptor(t *testing.T) {
	buf := []byte{0x00, 0x03}
	buf = append(buf, buildUTF8("(I)V")...)
	buf = append(buf, buildCPEntry(ConstantMethodType, 0x00, 0x01)...)
	cp := readPool(t, buf)
	desc, err := cp.MethodTypeDescriptor(2)
	if err != nil {
		t.Fatalf("MethodTypeDescriptor: %v", err)
	}
	if desc != "(I)V" {
		t.Errorf("descriptor = %q, want (I)V", desc)
	}
}

func TestMethodTypeDescriptorInvalid(t *testing.T) {
	buf := []byte{0x00, 0x03}
	buf = append(buf, buildUTF8("I")...)
	buf = append(buf, buildCPEntry(ConstantMethodType, 0x00, 0x01)...)
	cp := readPool(t, buf)
	_, err := cp.MethodTypeDescriptor(2)
	if err == nil {
		t.Fatal("expected error for field descriptor in MethodType")
	}
}

// TestInvokeDynamicInfo tests the typed accessor with checked validation.
func TestInvokeDynamicInfo(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("run")...)
	buf = append(buf, buildUTF8("()V")...)
	buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
	buf = append(buf, buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x03)...)
	cp := readPool(t, buf)
	bsIdx, name, desc, err := cp.InvokeDynamicInfo(4)
	if err != nil {
		t.Fatalf("InvokeDynamicInfo: %v", err)
	}
	if bsIdx != 0 {
		t.Errorf("bootstrap index = %d, want 0", bsIdx)
	}
	if name != "run" {
		t.Errorf("name = %q, want run", name)
	}
	if desc != "()V" {
		t.Errorf("descriptor = %q, want ()V", desc)
	}
}

func TestInvokeDynamicInfoWrongDescriptor(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("run")...)
	buf = append(buf, buildUTF8("I")...)
	buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
	buf = append(buf, buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x03)...)
	cp := readPool(t, buf)
	_, _, _, err := cp.InvokeDynamicInfo(4)
	if err == nil {
		t.Fatal("expected error for field descriptor in InvokeDynamic")
	}
}

// TestConstantDynamicFieldDescriptor validates ConstantDynamic with field descriptor.
func TestConstantDynamicFieldDescriptor(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("x")...)
	buf = append(buf, buildUTF8("I")...)
	buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
	buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...)
	cp := readPool(t, buf)
	_, _, _, err := cp.ConstantDynamicInfo(4)
	if err != nil {
		t.Fatalf("ConstantDynamicInfo with field descriptor: %v", err)
	}
}

func TestConstantDynamicRejectsMethodDescriptor(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("x")...)
	buf = append(buf, buildUTF8("()V")...)
	buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
	buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...)
	cp := readPool(t, buf)
	_, _, _, err := cp.ConstantDynamicInfo(4)
	if err == nil {
		t.Fatal("expected error for method descriptor in ConstantDynamic")
	}
}

// TestInvokeDynamicNatWrongTag verifies nat_index → non-NameAndType is rejected.
func TestInvokeDynamicNatWrongTag(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("x")...)
	buf = append(buf, buildUTF8("()V")...)
	buf = append(buf, buildUTF8("not-a-nat")...) // [3] Utf8 (NOT NameAndType!)
	buf = append(buf, buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x03)...)
	cp := readPool(t, buf)
	_, _, _, err := cp.InvokeDynamicInfo(4)
	if err == nil {
		t.Fatal("expected error for nat_index pointing at Utf8 instead of NameAndType")
	}
}

// TestConstantDynamicDecoding tests ConstantDynamicInfo with checked lookup.
func TestConstantDynamicDecoding(t *testing.T) {
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("x")...)
	buf = append(buf, buildUTF8("I")...)
	buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
	buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...)

	cp := readPool(t, buf)

	bsIdx, name, desc, err := cp.ConstantDynamicInfo(4)
	if err != nil {
		t.Fatalf("ConstantDynamicInfo: %v", err)
	}
	if bsIdx != 0 {
		t.Errorf("bootstrap index = %d, want 0", bsIdx)
	}
	if name != "x" {
		t.Errorf("name = %q, want x", name)
	}
	if desc != "I" {
		t.Errorf("descriptor = %q, want I", desc)
	}
}

func TestConstantDynamicWrongTag(t *testing.T) {
	buf := []byte{0x00, 0x02}
	buf = append(buf, buildUTF8("x")...)
	cp := readPool(t, buf)

	_, _, _, err := cp.ConstantDynamicInfo(1)
	if err == nil {
		t.Fatal("expected error for wrong tag")
	}
}

// TestConstantDynamicRejectsSpecialName verifies dynamic entries reject <init>/<clinit>.
func TestConstantDynamicRejectsSpecialName(t *testing.T) {
	for _, name := range []string{"<init>", "<clinit>"} {
		t.Run(name, func(t *testing.T) {
			buf := []byte{0x00, 0x05}
			buf = append(buf, buildUTF8(name)...)
			buf = append(buf, buildUTF8("I")...)
			buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
			buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...)
			cp := readPool(t, buf)
			_, _, _, err := cp.ConstantDynamicInfo(4)
			if err == nil {
				t.Fatalf("expected error for ConstantDynamic naming %q", name)
			}
		})
	}
}

// TestConstantDynamicAcceptsAngleBracketFieldName verifies that
// ConstantDynamic accepts names like <foo> — the NameAndType represents a
// field descriptor (JVMS §4.4.10), so field name rules apply. Only <init>
// and <clinit> are explicitly forbidden by JVMS §5.4.3.6; other names
// containing '<' or '>' are valid field identifiers.
func TestConstantDynamicAcceptsAngleBracketFieldName(t *testing.T) {
	for _, name := range []string{"<foo>", "<bar>"} {
		t.Run(name, func(t *testing.T) {
			buf := []byte{0x00, 0x05}
			buf = append(buf, buildUTF8(name)...)
			buf = append(buf, buildUTF8("I")...)
			buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
			buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...)
			cp := readPool(t, buf)
			_, _, _, err := cp.ConstantDynamicInfo(4)
			if err != nil {
				t.Fatalf("ConstantDynamicInfo for valid field name %q: %v", name, err)
			}
		})
	}
}

// TestConstantDynamicNatWrongTag verifies nat_index → non-NameAndType is rejected.
func TestConstantDynamicNatWrongTag(t *testing.T) {
	// [1] Utf8 "x", [2] Utf8 "I", [3] Utf8 "x" (NOT NameAndType!),
	// [4] ConstantDynamic(bootstrap=0, nat=3) — nat points at Utf8, not NameAndType.
	buf := []byte{0x00, 0x05}
	buf = append(buf, buildUTF8("x")...)                                        // [1] (unused)
	buf = append(buf, buildUTF8("I")...)                                        // [2] (unused)
	buf = append(buf, buildUTF8("not-a-nat")...)                                // [3] Utf8 (NOT NameAndType!)
	buf = append(buf, buildCPEntry(ConstantDynamic, 0x00, 0x00, 0x00, 0x03)...) // [4] nat=3→Utf8
	cp := readPool(t, buf)
	_, _, _, err := cp.ConstantDynamicInfo(4)
	if err == nil {
		t.Fatal("expected error for nat_index pointing at Utf8 instead of NameAndType")
	}
}

// TestInvokeDynamicRejectsSpecialName verifies InvokeDynamic rejects <init>/<clinit>.
func TestInvokeDynamicRejectsSpecialName(t *testing.T) {
	for _, name := range []string{"<init>", "<clinit>"} {
		t.Run(name, func(t *testing.T) {
			buf := []byte{0x00, 0x05}
			buf = append(buf, buildUTF8(name)...)
			buf = append(buf, buildUTF8("()V")...)
			buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
			buf = append(buf, buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x03)...)
			cp := readPool(t, buf)
			_, _, _, err := cp.InvokeDynamicInfo(4)
			if err == nil {
				t.Fatalf("expected error for InvokeDynamic naming %q", name)
			}
		})
	}
}

// TestInvokeDynamicRejectsAngleBracketNames verifies that InvokeDynamic
// rejects names containing '<' or '>' (per JVMS §4.2.2 ordinary method
// name rules). The NameAndType represents a method descriptor (JVMS §4.4.10).
func TestInvokeDynamicRejectsAngleBracketNames(t *testing.T) {
	for _, name := range []string{"<foo>", "<bar>", "<>"} {
		t.Run(name, func(t *testing.T) {
			buf := []byte{0x00, 0x05}
			buf = append(buf, buildUTF8(name)...)
			buf = append(buf, buildUTF8("()V")...)
			buf = append(buf, buildCPEntry(ConstantNameAndType, 0x00, 0x01, 0x00, 0x02)...)
			buf = append(buf, buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x03)...)
			cp := readPool(t, buf)
			_, _, _, err := cp.InvokeDynamicInfo(4)
			if err == nil {
				t.Fatalf("expected error for InvokeDynamic naming %q", name)
			}
		})
	}
}

// --- Phase 4: BootstrapMethods tests ---

// buildBootstrapMethodsAttr builds the raw bytes of a BootstrapMethods
// attribute given the attribute name index and a list of (methodRef, args)
// pairs.
func buildBootstrapMethodsAttr(nameIndex uint16, entries []struct {
	methodRef uint16
	args      []uint16
}) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(entries)))
	for _, e := range entries {
		buf = binary.BigEndian.AppendUint16(buf, e.methodRef)
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(e.args)))
		for _, a := range e.args {
			buf = binary.BigEndian.AppendUint16(buf, a)
		}
	}
	// Wrap in attribute format: name_index + length + body.
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, nameIndex)
	attr = binary.BigEndian.AppendUint32(attr, uint32(len(buf)))
	attr = append(attr, buf...)
	return attr
}

// buildClassAttr wraps raw attribute bytes in the class attribute format
// (u2 count followed by attributes).
func buildClassAttr(attrs ...[]byte) []byte {
	buf := binary.BigEndian.AppendUint16(nil, uint16(len(attrs)))
	for _, a := range attrs {
		buf = append(buf, a...)
	}
	return buf
}

// buildMinimalClassfile constructs a minimal parseable classfile with the given
// constant-pool entries and class-level attributes.
func buildMinimalClassfile(cpEntries [][]byte, classAttrs [][]byte, majorVersion uint16) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, majorVersion)
	// cp_count = entries + 1 (slot 0)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cpEntries)+1))
	for _, e := range cpEntries {
		buf = append(buf, e...)
	}
	// this_class always index 1
	buf = binary.BigEndian.AppendUint16(buf, 0x0001) // access_flags
	buf = binary.BigEndian.AppendUint16(buf, 0x0002) // this_class -> assume index 2 is Class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super_class (Object)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces_count
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // fields_count
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // methods_count
	buf = append(buf, buildClassAttr(classAttrs...)...)
	return buf
}

func TestBootstrapMethodsValid(t *testing.T) {
	// Build a classfile with:
	// [1] CONSTANT_Utf8 "C"
	// [2] CONSTANT_Class -> 1
	// [3] CONSTANT_Utf8 "m"
	// [4] CONSTANT_Utf8 "()V"
	// [5] CONSTANT_NameAndType -> (3, 4)
	// [6] CONSTANT_Methodref -> (2, 5)
	// [7] CONSTANT_MethodHandle (kind=6, ref=6)
	// [8] CONSTANT_Utf8 "BootstrapMethods"
	// BootstrapMethods attribute: 1 entry (methodRef=7, args=[])

	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
	}
	attrs := [][]byte{
		buildBootstrapMethodsAttr(8, []struct {
			methodRef uint16
			args      []uint16
		}{{methodRef: 7}}),
	}
	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err != nil {
		t.Fatalf("Parse with valid BootstrapMethods: %v", err)
	}
	bm := cf.BootstrapMethods()
	if bm == nil {
		t.Fatal("BootstrapMethods is nil")
	}
	if bm.NumEntries() != 1 {
		t.Fatalf("NumEntries = %d, want 1", bm.NumEntries())
	}
	e := bm.Entry(0)
	if e.MethodRef != 7 {
		t.Errorf("Entry.MethodRef = %d, want 7", e.MethodRef)
	}
	if len(e.Arguments) != 0 {
		t.Errorf("len(Arguments) = %d, want 0", len(e.Arguments))
	}
}

func TestBootstrapMethodsMultipleArgs(t *testing.T) {
	// BootstrapMethods with 3 arguments preserves order.
	// [1]-[8] same as above
	// [9] CONSTANT_Integer (value 42)
	// [10] CONSTANT_Integer (value 99)
	// [11] CONSTANT_String -> "s" UTF8
	// [12] CONSTANT_Utf8 "s"
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
		buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x2A),           // [9] val=42
		buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x63),           // [10] val=99
		buildUTF8("s"),                           // [11]
		buildCPEntry(ConstantString, 0x00, 0x0B), // [12] -> 11
	}
	attrs := [][]byte{
		buildBootstrapMethodsAttr(8, []struct {
			methodRef uint16
			args      []uint16
		}{{methodRef: 7, args: []uint16{9, 10, 12}}}),
	}
	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	e := cf.BootstrapMethods().Entry(0)
	if len(e.Arguments) != 3 {
		t.Fatalf("len(Arguments) = %d, want 3", len(e.Arguments))
	}
	if e.Arguments[0] != 9 || e.Arguments[1] != 10 || e.Arguments[2] != 12 {
		t.Errorf("Arguments = %v, want [9, 10, 12]", e.Arguments)
	}
}

func TestBootstrapMethodsArgumentImmutability(t *testing.T) {
	// Mutating the returned argument slice must not change retained metadata.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
		buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x2A),           // [9]
	}
	attrs := [][]byte{
		buildBootstrapMethodsAttr(8, []struct {
			methodRef uint16
			args      []uint16
		}{{methodRef: 7, args: []uint16{9}}}),
	}
	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	e1 := cf.BootstrapMethods().Entry(0)
	e1.Arguments[0] = 999 // mutate the copy
	e2 := cf.BootstrapMethods().Entry(0)
	if e2.Arguments[0] != 9 {
		t.Errorf("Argument mutated: got %d, want 9", e2.Arguments[0])
	}
}

func TestBootstrapMethodsMissingForInvokeDynamic(t *testing.T) {
	// InvokeDynamic entry but no BootstrapMethods attribute → fail.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("run"),                        // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),   // [5]
		buildCPEntry(ConstantInvokeDynamic, 0x00, 0x00, 0x00, 0x05), // [6]
	}
	_, err := Parse(buildMinimalClassfile(cp, nil, 69))
	if err == nil {
		t.Fatal("expected FormatError for InvokeDynamic without BootstrapMethods")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "constant pool" {
		t.Errorf("FormatError.Op = %q, want constant pool", fe.Op)
	}
}

func TestBootstrapMethodsDuplicate(t *testing.T) {
	// Two BootstrapMethods attributes → fail.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
	}
	attr := buildBootstrapMethodsAttr(8, []struct {
		methodRef uint16
		args      []uint16
	}{{methodRef: 7}})
	attrs := [][]byte{attr, attr} // duplicate
	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for duplicate BootstrapMethods")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "BootstrapMethods" {
		t.Errorf("FormatError.Op = %q, want BootstrapMethods", fe.Op)
	}
}

func TestBootstrapMethodsWrongMethodRefTag(t *testing.T) {
	// BootstrapMethods entry references a Class entry instead of MethodHandle.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("BootstrapMethods"),           // [3]
	}
	attrs := [][]byte{
		buildBootstrapMethodsAttr(3, []struct {
			methodRef uint16
			args      []uint16
		}{{methodRef: 2}}), // refs Class tag 7, not MethodHandle tag 15
	}
	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for wrong method_ref tag")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "BootstrapMethods" {
		t.Errorf("FormatError.Op = %q, want BootstrapMethods", fe.Op)
	}
}

func TestBootstrapMethodsInvalidArgument(t *testing.T) {
	// BootstrapMethods argument references a UTF8 entry (not loadable).
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("m"),                          // [3]
		buildUTF8("()V"),                        // [4]
		buildCPEntry(ConstantNameAndType, 0x00, 0x03, 0x00, 0x04),       // [5]
		buildCPEntry(ConstantMethodref, 0x00, 0x02, 0x00, 0x05),         // [6]
		buildCPEntry(ConstantMethodHandle, RefInvokeStatic, 0x00, 0x06), // [7]
		buildUTF8("BootstrapMethods"),                                   // [8]
		buildUTF8("not-loadable"),                                       // [9]
	}
	attrs := [][]byte{
		buildBootstrapMethodsAttr(8, []struct {
			methodRef uint16
			args      []uint16
		}{{methodRef: 7, args: []uint16{9}}}), // 9 = UTF8, not loadable
	}
	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for invalid argument")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "BootstrapMethods" {
		t.Errorf("FormatError.Op = %q, want BootstrapMethods", fe.Op)
	}
}

func TestBootstrapMethodsEmpty(t *testing.T) {
	// Empty BootstrapMethods table is valid (only rejected when referenced
	// by dynamic entries).
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("BootstrapMethods"),           // [3]
	}
	// Build an empty BootstrapMethods attribute manually.
	body := []byte{0x00, 0x00} // num_bootstrap_methods = 0
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 3) // name_index = 3
	attr = binary.BigEndian.AppendUint32(attr, uint32(len(body)))
	attr = append(attr, body...)
	attrs := [][]byte{attr}
	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err != nil {
		t.Fatalf("Parse with empty BootstrapMethods: %v", err)
	}
	bm := cf.BootstrapMethods()
	if bm == nil {
		t.Fatal("BootstrapMethods is nil for empty table")
	}
	if bm.NumEntries() != 0 {
		t.Fatalf("NumEntries = %d, want 0", bm.NumEntries())
	}
}

// --- Attribute location tests (Code, StackMapTable, BootstrapMethods) ---

// buildMemberBytes builds raw bytes for a single field_info or method_info.
func buildMemberBytes(accessFlags, nameIndex, descIndex uint16, attrs []byte) []byte {
	buf := []byte{}
	buf = binary.BigEndian.AppendUint16(buf, accessFlags)
	buf = binary.BigEndian.AppendUint16(buf, nameIndex)
	buf = binary.BigEndian.AppendUint16(buf, descIndex)
	buf = append(buf, attrs...)
	return buf
}

func TestCodeAttributeOnField(t *testing.T) {
	// Build a classfile with a Code attribute on a field → should fail.
	// [1] CONSTANT_Utf8 "C"
	// [2] CONSTANT_Class -> 1
	// [3] CONSTANT_Utf8 "f"
	// [4] CONSTANT_Utf8 "I"
	// [5] CONSTANT_Utf8 "Code"
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("f"),                          // [3]
		buildUTF8("I"),                          // [4]
		buildUTF8("Code"),                       // [5]
	}
	// Code attribute body: max_stack=0, max_locals=0, code_length=0,
	// exception_table_length=0, attributes_count=0
	codeBody := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	codeAttr := []byte{}
	codeAttr = binary.BigEndian.AppendUint16(codeAttr, 5) // name_index = 5 ("Code")
	codeAttr = binary.BigEndian.AppendUint32(codeAttr, uint32(len(codeBody)))
	codeAttr = append(codeAttr, codeBody...)

	// Wrap in attributes (u2 count + attribute).
	fieldAttrs := binary.BigEndian.AppendUint16(nil, 1)
	fieldAttrs = append(fieldAttrs, codeAttr...)

	field := buildMemberBytes(0, 3, 4, fieldAttrs) // access_flags, name=3("f"), desc=4("I")

	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001) // access_flags
	buf = binary.BigEndian.AppendUint16(buf, 0x0002) // this_class -> 2
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super_class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces_count
	buf = binary.BigEndian.AppendUint16(buf, 1)      // fields_count = 1
	buf = append(buf, field...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // methods_count
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // attributes_count

	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for Code attribute on field")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "Code attribute" {
		t.Errorf("FormatError.Op = %q, want Code attribute", fe.Op)
	}
}

func TestStackMapTableOnClass(t *testing.T) {
	// Build a classfile with a StackMapTable attribute at the class level → fail.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("StackMapTable"),              // [3]
	}
	// StackMapTable body: number_of_entries = 0
	smtBody := []byte{0x00, 0x00}
	smtAttr := []byte{}
	smtAttr = binary.BigEndian.AppendUint16(smtAttr, 3) // name_index = 3
	smtAttr = binary.BigEndian.AppendUint32(smtAttr, uint32(len(smtBody)))
	smtAttr = append(smtAttr, smtBody...)
	attrs := [][]byte{smtAttr}

	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for StackMapTable on class")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "StackMapTable" {
		t.Errorf("FormatError.Op = %q, want StackMapTable", fe.Op)
	}
}

func TestBootstrapMethodsOnMethod(t *testing.T) {
	// Build a classfile with a BootstrapMethods attribute on a method → fail.
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("BootstrapMethods"),           // [3]
		buildUTF8("m"),                          // [4]
		buildUTF8("()V"),                        // [5]
	}
	// BootstrapMethods body: num_bootstrap_methods = 1 (with a dummy entry)
	bsBody := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
	bsAttr := []byte{}
	bsAttr = binary.BigEndian.AppendUint16(bsAttr, 3) // name_index = 3 ("BootstrapMethods")
	bsAttr = binary.BigEndian.AppendUint32(bsAttr, uint32(len(bsBody)))
	bsAttr = append(bsAttr, bsBody...)

	methodAttrs := binary.BigEndian.AppendUint16(nil, 1)
	methodAttrs = append(methodAttrs, bsAttr...)

	method := buildMemberBytes(0x0009, 4, 5, methodAttrs) // public static, name=4("m"), desc=5("()V")

	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001) // access_flags
	buf = binary.BigEndian.AppendUint16(buf, 0x0002) // this_class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // super_class
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // interfaces_count
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // fields_count
	buf = binary.BigEndian.AppendUint16(buf, 1)      // methods_count = 1
	buf = append(buf, method...)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000) // attributes_count

	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for BootstrapMethods on method")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
	if fe.Op != "BootstrapMethods" {
		t.Errorf("FormatError.Op = %q, want BootstrapMethods", fe.Op)
	}
}

// --- Attribute name_index validation tests ---
// Five categories of malformed attribute_name_index, each producing *FormatError
// (not silently ignored as unknown).

// TestAttrNameIndexZero verifies that an attribute with name_index==0 is
// rejected as a *FormatError.
func TestAttrNameIndexZero(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
	}
	// Attribute with name_index = 0 (invalid).
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 0) // name_index = 0
	attr = binary.BigEndian.AppendUint32(attr, 0) // length = 0
	attrs := [][]byte{attr}

	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for attribute name_index 0")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestAttrNameIndexOOB verifies that an attribute with out-of-bounds
// name_index is rejected as a *FormatError.
func TestAttrNameIndexOOB(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
	}
	// Only 2 entries (index 0 + 1), index 99 is OOB.
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 99) // name_index OOB
	attr = binary.BigEndian.AppendUint32(attr, 0)  // length = 0
	attrs := [][]byte{attr}

	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for OOB attribute name_index")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestAttrNameIndexSecondSlot verifies that an attribute name_index pointing
// at a long/double second slot is rejected.
func TestAttrNameIndexSecondSlot(t *testing.T) {
	// Build a pool where index 2 is the unusable second slot of a long.
	// [1] CONSTANT_Long (takes slots 1 and 2)
	// [2] nil (second slot)
	// The attribute name_index will be 2 (the unusable slot).
	cp := [][]byte{
		buildCPEntry(ConstantLong, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42), // [1]
		buildUTF8("C"),                          // [3] — note: index 2 is skipped by long
		buildCPEntry(ConstantClass, 0x00, 0x03), // [4]
	}
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 2) // name_index = 2 (unusable second slot)
	attr = binary.BigEndian.AppendUint32(attr, 0) // length = 0
	attrs := [][]byte{attr}

	// Override this_class to use index 4 (Class) to avoid further errors.
	buf := []byte{}
	buf = binary.BigEndian.AppendUint32(buf, 0xCAFEBABE)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, 69)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(cp)+1))
	for _, e := range cp {
		buf = append(buf, e...)
	}
	buf = binary.BigEndian.AppendUint16(buf, 0x0001)
	buf = binary.BigEndian.AppendUint16(buf, 4) // this_class -> index 4 (Class)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = binary.BigEndian.AppendUint16(buf, 0x0000)
	buf = append(buf, buildClassAttr(attrs...)...)

	_, err := Parse(buf)
	if err == nil {
		t.Fatal("expected FormatError for attribute name_index pointing at second slot")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestAttrNameIndexNotUtf8 verifies that an attribute name_index pointing
// at a non-Utf8 entry (e.g. Integer) is rejected.
func TestAttrNameIndexNotUtf8(t *testing.T) {
	cp := [][]byte{
		buildCPEntry(ConstantInteger, 0x00, 0x00, 0x00, 0x2A), // [1]
		buildUTF8("C"),                          // [2]
		buildCPEntry(ConstantClass, 0x00, 0x02), // [3]
	}
	// name_index = 1 points at ConstantInteger, not Utf8.
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 1) // name_index -> Integer
	attr = binary.BigEndian.AppendUint32(attr, 0) // length = 0
	attrs := [][]byte{attr}

	_, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err == nil {
		t.Fatal("expected FormatError for non-UTF8 attribute name_index")
	}
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("expected *FormatError, got %T: %v", err, err)
	}
}

// TestAttrNameIndexValidButUnknown verifies that a syntactically valid
// attribute name that is unrecognized is NOT rejected — unknown attributes
// are silently consumed, not an error.
func TestAttrNameIndexValidButUnknown(t *testing.T) {
	cp := [][]byte{
		buildUTF8("C"),                          // [1]
		buildCPEntry(ConstantClass, 0x00, 0x01), // [2]
		buildUTF8("UnknownCustomAttribute"),     // [3]
	}
	attr := []byte{}
	attr = binary.BigEndian.AppendUint16(attr, 3) // name_index = "UnknownCustomAttribute"
	attr = binary.BigEndian.AppendUint32(attr, 0) // length = 0
	attrs := [][]byte{attr}

	cf, err := Parse(buildMinimalClassfile(cp, attrs, 69))
	if err != nil {
		t.Fatalf("Parse with unknown attribute name: %v", err)
	}
	if cf == nil {
		t.Fatal("Parse returned nil ClassFile for valid-but-unknown attribute")
	}
}

// --- JDK 25 integration fixture tests ---

// compileModernFixture compiles tests/fixtures/<name>.java with the system JDK
// (no source/target restriction) into t.TempDir() and returns the class bytes.
func compileModernFixture(t *testing.T, name string) []byte {
	t.Helper()
	src := filepath.Join("..", "tests", "fixtures", name+".java")
	out := t.TempDir()
	cmd := exec.Command("javac", "-d", out, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("javac %s failed: %v\n%s", name, err, out)
	}
	data, err := os.ReadFile(filepath.Join(out, name+".class"))
	if err != nil {
		t.Fatalf("read class: %v", err)
	}
	return data
}

func TestParseJDK25StringConcat(t *testing.T) {
	data := compileModernFixture(t, "DynStringConcat")
	cf, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse DynStringConcat: %v", err)
	}
	if cf.ClassName() != "DynStringConcat" {
		t.Errorf("ClassName = %q, want DynStringConcat", cf.ClassName())
	}
	if cf.MajorVersion() < 65 {
		t.Errorf("MajorVersion = %d, want >= 65", cf.MajorVersion())
	}

	// Must have the greet method.
	foundGreet := false
	for _, m := range cf.Methods() {
		if m.Name() == "greet" && m.Descriptor() == "(Ljava/lang/String;)Ljava/lang/String;" {
			foundGreet = true
		}
	}
	if !foundGreet {
		t.Error("greet method not found")
	}

	// BootstrapMethods must be present and non-empty.
	bm := cf.BootstrapMethods()
	if bm == nil {
		t.Fatal("BootstrapMethods is nil for string concat class")
	}
	if bm.NumEntries() == 0 {
		t.Fatal("BootstrapMethods is empty for string concat class")
	}

	pool := cf.ConstantPool()
	major := cf.MajorVersion()

	// Must find at least one InvokeDynamic.
	indyFound := false
	for i := uint16(1); i < uint16(pool.Size()); i++ {
		if pool.Tag(i) == ConstantInvokeDynamic {
			indyFound = true
			bsIdx, name, desc, infoErr := pool.InvokeDynamicInfo(i)
			if infoErr != nil {
				t.Errorf("InvokeDynamicInfo(%d): %v", i, infoErr)
				continue
			}
			// Name and descriptor must be non-empty.
			if name == "" {
				t.Error("InvokeDynamic name is empty")
			}
			if desc == "" {
				t.Error("InvokeDynamic descriptor is empty")
			}
			// Bootstrap index must be in range.
			if int(bsIdx) >= bm.NumEntries() {
				t.Errorf("InvokeDynamic bootstrap index %d out of range (table has %d entries)", bsIdx, bm.NumEntries())
				continue
			}
			// The corresponding BootstrapMethods method ref must be a valid MethodHandle.
			e := bm.Entry(int(bsIdx))
			_, _, _, _, _, mhErr := pool.MethodHandleReference(e.MethodRef, major)
			if mhErr != nil {
				t.Errorf("BootstrapMethods entry %d method ref %d: %v", bsIdx, e.MethodRef, mhErr)
			}
			// Check bootstrap arguments order and tags.
			for j, argIdx := range e.Arguments {
				argTag := pool.Tag(argIdx)
				t.Logf("BootstrapMethods entry %d arg[%d]: index=%d tag=%d", bsIdx, j, argIdx, argTag)
				if !loadableConstant(argTag) {
					t.Errorf("BootstrapMethods entry %d arg[%d]: index %d tag %d is not loadable", bsIdx, j, argIdx, argTag)
				}
			}
		}
		if pool.Tag(i) == ConstantLong || pool.Tag(i) == ConstantDouble {
			i++
		}
	}
	if !indyFound {
		t.Fatal("no InvokeDynamic entry found in string concat class")
	}
}

func TestParseJDK25Lambda(t *testing.T) {
	data := compileModernFixture(t, "DynLambda")
	cf, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse DynLambda: %v", err)
	}
	if cf.ClassName() != "DynLambda" {
		t.Errorf("ClassName = %q, want DynLambda", cf.ClassName())
	}

	bm := cf.BootstrapMethods()
	if bm == nil {
		t.Fatal("BootstrapMethods is nil for lambda class")
	}
	if bm.NumEntries() == 0 {
		t.Fatal("BootstrapMethods is empty for lambda class")
	}

	pool := cf.ConstantPool()
	major := cf.MajorVersion()

	// Must find InvokeDynamic.
	indyFound := false
	for i := uint16(1); i < uint16(pool.Size()); i++ {
		if pool.Tag(i) == ConstantInvokeDynamic {
			indyFound = true
			bsIdx, name, desc, infoErr := pool.InvokeDynamicInfo(i)
			if infoErr != nil {
				t.Errorf("InvokeDynamicInfo(%d): %v", i, infoErr)
				continue
			}
			if name == "" {
				t.Error("InvokeDynamic name is empty")
			}
			if desc == "" {
				t.Error("InvokeDynamic descriptor is empty")
			}
			if int(bsIdx) >= bm.NumEntries() {
				t.Errorf("InvokeDynamic bootstrap index %d out of range", bsIdx)
				continue
			}

			// Read the bootstrap method entry.
			e := bm.Entry(int(bsIdx))

			// Verify the method ref is a valid MethodHandle.
			_, _, _, _, _, mhErr := pool.MethodHandleReference(e.MethodRef, major)
			if mhErr != nil {
				t.Errorf("BootstrapMethods entry %d method ref %d: %v", bsIdx, e.MethodRef, mhErr)
			}

			// Verify MethodType and MethodHandle arguments — both MUST be found.
			foundMethodType := false
			foundMethodHandle := false
			for j, argIdx := range e.Arguments {
				argTag := pool.Tag(argIdx)
				switch argTag {
				case ConstantMethodType:
					desc, mtErr := pool.MethodTypeDescriptor(argIdx)
					if mtErr != nil {
						t.Errorf("Bootstrap arg %d (MethodType %d): %v", j, argIdx, mtErr)
					} else if desc == "" {
						t.Errorf("Bootstrap arg %d (MethodType %d): empty descriptor", j, argIdx)
					}
					foundMethodType = true
				case ConstantMethodHandle:
					_, _, _, _, _, mhErr := pool.MethodHandleReference(argIdx, major)
					if mhErr != nil {
						t.Errorf("Bootstrap arg %d (MethodHandle %d): %v", j, argIdx, mhErr)
					}
					foundMethodHandle = true
				}
			}
			// Mandatory assertions: a lambda bootstrap MUST have both.
			if !foundMethodType {
				t.Error("no MethodType found in BootstrapMethods arguments for lambda — expected MethodType is required for lambda metafactory")
			}
			if !foundMethodHandle {
				t.Error("no MethodHandle found in BootstrapMethods arguments for lambda — expected MethodHandle is required for lambda metafactory")
			}
		}
		if pool.Tag(i) == ConstantLong || pool.Tag(i) == ConstantDouble {
			i++
		}
	}
	if !indyFound {
		t.Fatal("no InvokeDynamic entry found in lambda class")
	}
}
