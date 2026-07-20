package rtda

import (
	"sync"
	"testing"
)

func TestGetArrayClassNamePrimitive(t *testing.T) {
	InitVMTypes()

	tests := []struct {
		primitive *Class
		wantName  string
		wantKind  int
	}{
		{VMPrimitiveInt, "[I", kindInt},
		{VMPrimitiveLong, "[J", kindLong},
		{VMPrimitiveFloat, "[F", kindFloat},
		{VMPrimitiveDouble, "[D", kindDouble},
		{VMPrimitiveByte, "[B", kindByte},
		{VMPrimitiveChar, "[C", kindChar},
		{VMPrimitiveShort, "[S", kindShort},
		{VMPrimitiveBool, "[Z", kindBoolean},
	}

	for _, tt := range tests {
		arr := tt.primitive.GetArrayClass()
		if arr.Name() != tt.wantName {
			t.Errorf("%s.[]: name = %q, want %q", tt.primitive.Name(), arr.Name(), tt.wantName)
		}
		if !arr.IsArray() {
			t.Errorf("%s.[]: not an array class", tt.primitive.Name())
		}
		if arr.ComponentClass() != tt.primitive {
			t.Errorf("%s.[]: componentClass = %p, want %p", tt.primitive.Name(), arr.ComponentClass(), tt.primitive)
		}
		if arr.ComponentKind() != tt.wantKind {
			t.Errorf("%s.[]: componentKind = %d, want %d", tt.primitive.Name(), arr.ComponentKind(), tt.wantKind)
		}
		if arr.DefiningLoader() != VMIdentity {
			t.Errorf("%s.[]: definingLoader = %v, want VMIdentity", tt.primitive.Name(), arr.DefiningLoader())
		}
	}
}

func TestGetArrayClassNameReference(t *testing.T) {
	// Simulate a reference class: not primitive, not array.
	c := NewSyntheticClass("java/lang/String", nil)

	arr := c.GetArrayClass()
	wantName := "[Ljava/lang/String;"
	if arr.Name() != wantName {
		t.Errorf("String[]: name = %q, want %q", arr.Name(), wantName)
	}
	if arr.ComponentKind() != kindNone {
		t.Errorf("String[]: componentKind = %d, want kindNone (0)", arr.ComponentKind())
	}
	if arr.ComponentClass() != c {
		t.Errorf("String[]: componentClass = %p, want %p", arr.ComponentClass(), c)
	}
}

func TestGetArrayClassNameDeepNesting(t *testing.T) {
	InitVMTypes()

	// int[][][] — three levels of array nesting.
	i := VMPrimitiveInt
	i1 := i.GetArrayClass()  // [I    — kind=int
	i2 := i1.GetArrayClass() // [[I   — kind=none (component is array)
	i3 := i2.GetArrayClass() // [[[I  — kind=none

	if i1.Name() != "[I" {
		t.Errorf("int[]: name = %q, want [I", i1.Name())
	}
	if i1.ComponentKind() != kindInt {
		t.Errorf("int[]: componentKind = %d, want kindInt", i1.ComponentKind())
	}
	if i2.Name() != "[[I" {
		t.Errorf("int[][]: name = %q, want [[I", i2.Name())
	}
	if i2.ComponentKind() != kindNone {
		t.Errorf("int[][]: componentKind = %d, want kindNone (component is array)", i2.ComponentKind())
	}
	if i3.Name() != "[[[I" {
		t.Errorf("int[][][]: name = %q, want [[[I", i3.Name())
	}

	// java/lang/String[][]
	s := NewSyntheticClass("java/lang/String", nil)
	s1 := s.GetArrayClass()  // [Ljava/lang/String; — kind=none (reference)
	s2 := s1.GetArrayClass() // [[Ljava/lang/String; — kind=none

	if s1.Name() != "[Ljava/lang/String;" {
		t.Errorf("String[]: name = %q, want [Ljava/lang/String;", s1.Name())
	}
	if s2.Name() != "[[Ljava/lang/String;" {
		t.Errorf("String[][]: name = %q, want [[Ljava/lang/String;", s2.Name())
	}
}

func TestGetArrayClassIdentityStable(t *testing.T) {
	InitVMTypes()

	// Multiple calls to GetArrayClass must return the same pointer.
	i := VMPrimitiveInt
	a1 := i.GetArrayClass()
	a2 := i.GetArrayClass()
	if a1 != a2 {
		t.Error("GetArrayClass returned different pointers for same component")
	}

	// Deep nesting identity.
	s := NewSyntheticClass("java/lang/String", nil)
	s1 := s.GetArrayClass()
	s2 := s1.GetArrayClass()
	s3 := s2.ComponentClass().GetArrayClass()
	if s2 != s3 {
		t.Error("String[][]: GetArrayClass from componentClass gave different pointer")
	}
}

func TestGetArrayClassConcurrentSameIdentity(t *testing.T) {
	InitVMTypes()

	i := VMPrimitiveInt
	const N = 32
	var wg sync.WaitGroup
	results := make([]*Class, N)

	for n := 0; n < N; n++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = i.GetArrayClass()
		}(n)
	}
	wg.Wait()

	for idx := 1; idx < N; idx++ {
		if results[idx] != results[0] {
			t.Errorf("goroutine %d got %p, want %p", idx, results[idx], results[0])
		}
	}
}

func TestPrimitiveDescriptor(t *testing.T) {
	if d := PrimitiveDescriptor("int"); d != 'I' {
		t.Errorf("PrimitiveDescriptor(int) = %c, want I", d)
	}
	if d := PrimitiveDescriptor("boolean"); d != 'Z' {
		t.Errorf("PrimitiveDescriptor(boolean) = %c, want Z", d)
	}
	if d := PrimitiveDescriptor("void"); d != 'V' {
		t.Errorf("PrimitiveDescriptor(void) = %c, want V", d)
	}
	if d := PrimitiveDescriptor("String"); d != 0 {
		t.Errorf("PrimitiveDescriptor(String) = %c, want 0", d)
	}
}

func TestVMPrimitiveForNameLazyInit(t *testing.T) {
	// Call BEFORE InitVMTypes — must self-initialise.
	// Reset the once (test-only).
	vmTypesOnce = sync.Once{}
	for _, name := range []string{"boolean", "byte", "char", "short", "int", "long", "float", "double", "void"} {
		c := VMPrimitiveForName(name)
		if c == nil {
			t.Errorf("VMPrimitiveForName(%q) = nil", name)
		}
		if c.Name() != name {
			t.Errorf("VMPrimitiveForName(%q).Name() = %q", name, c.Name())
		}
		if c.DefiningLoader() != VMIdentity {
			t.Errorf("VMPrimitiveForName(%q).DefiningLoader() = %v, want VMIdentity", name, c.DefiningLoader())
		}
	}
}

func TestVMPrimitiveComponentKind(t *testing.T) {
	InitVMTypes()

	tests := map[string]int{
		"boolean": kindBoolean,
		"byte":    kindByte,
		"char":    kindChar,
		"short":   kindShort,
		"int":     kindInt,
		"long":    kindLong,
		"float":   kindFloat,
		"double":  kindDouble,
	}
	for name, wantKind := range tests {
		c := VMPrimitiveForName(name)
		if c == nil {
			t.Fatalf("VMPrimitiveForName(%q) = nil", name)
		}
		if c.ComponentKind() != wantKind {
			t.Errorf("ComponentKind(%s) = %d, want %d", name, c.ComponentKind(), wantKind)
		}
	}
	// void has kind 0.
	if VMVoid.ComponentKind() != 0 {
		t.Errorf("void.ComponentKind() = %d, want 0", VMVoid.ComponentKind())
	}
}
