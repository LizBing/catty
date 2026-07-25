package rtda

import (
	"errors"
	"math"
	"testing"
)

func TestJavaValueConstructorsPreserveLogicalKinds(t *testing.T) {
	ref := &Object{}
	tests := []struct {
		name  string
		value JavaValue
		kind  JavaValueKind
	}{
		{"void", VoidValue(), JavaValueVoid},
		{"boolean", BooleanValue(true), JavaValueBoolean},
		{"byte", ByteValue(-12), JavaValueByte},
		{"char", CharValue(0xfffe), JavaValueChar},
		{"short", ShortValue(-1234), JavaValueShort},
		{"int", IntValue(-1234567), JavaValueInt},
		{"long", LongValue(math.MinInt64 + 9), JavaValueLong},
		{"float", FloatValue(math.Float32frombits(0x7fc01234)), JavaValueFloat},
		{"double", DoubleValue(math.Float64frombits(0x7ff8000000001234)), JavaValueDouble},
		{"reference", ReferenceValue(ref), JavaValueReference},
		{"null reference", ReferenceValue(nil), JavaValueReference},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.value.Kind(); got != test.kind {
				t.Fatalf("Kind() = %s, want %s", got, test.kind)
			}
		})
	}
	if got, ok := BooleanValue(true).Boolean(); !ok || !got {
		t.Fatal("boolean getter lost value")
	}
	if got, ok := ByteValue(-12).Byte(); !ok || got != -12 {
		t.Fatal("byte getter lost value")
	}
	if got, ok := CharValue(0xfffe).Char(); !ok || got != 0xfffe {
		t.Fatal("char getter lost value")
	}
	if got, ok := ShortValue(-1234).Short(); !ok || got != -1234 {
		t.Fatal("short getter lost value")
	}
	if got, ok := IntValue(-1234567).Int(); !ok || got != -1234567 {
		t.Fatal("int getter lost value")
	}
	if got, ok := LongValue(math.MinInt64 + 9).Long(); !ok || got != math.MinInt64+9 {
		t.Fatal("long getter lost value")
	}
	if got, ok := FloatValue(3.25).Float(); !ok || got != 3.25 {
		t.Fatal("float getter lost value")
	}
	if got, ok := DoubleValue(-3.25).Double(); !ok || got != -3.25 {
		t.Fatal("double getter lost value")
	}
	if got, ok := ReferenceValue(ref).Reference(); !ok || got != ref {
		t.Fatal("reference getter lost identity")
	}
	if _, ok := IntValue(1).Long(); ok {
		t.Fatal("logical value must not silently widen")
	}
}

func TestDynamicResultDistinguishesNormalThrowableAndFailure(t *testing.T) {
	throwable := &Object{}
	failure := errors.New("adapter defect")

	normal := NormalResult(ReferenceValue(nil))
	if !normal.IsNormal() || normal.IsThrowable() || normal.IsInternalFailure() {
		t.Fatal("normal null result has wrong state")
	}
	if value, ok := normal.Value(); !ok || value.Kind() != JavaValueReference {
		t.Fatal("normal null result lost its value category")
	}

	thrown := ThrowableResult(throwable)
	if !thrown.IsThrowable() || thrown.IsNormal() || thrown.IsInternalFailure() {
		t.Fatal("throwable result has wrong state")
	}
	if got, ok := thrown.Throwable(); !ok || got != throwable {
		t.Fatal("throwable identity was not preserved")
	}

	internal := InternalFailureResult(failure)
	if !internal.IsInternalFailure() || internal.IsNormal() || internal.IsThrowable() {
		t.Fatal("internal failure has wrong state")
	}
	if !errors.Is(internal.Failure(), failure) {
		t.Fatal("internal failure was not preserved")
	}
}

func TestJavaValueKindForDescriptor(t *testing.T) {
	tests := map[string]JavaValueKind{
		"V": JavaValueVoid, "Z": JavaValueBoolean, "B": JavaValueByte,
		"C": JavaValueChar, "S": JavaValueShort, "I": JavaValueInt,
		"J": JavaValueLong, "F": JavaValueFloat, "D": JavaValueDouble,
		"Ljava/lang/Object;": JavaValueReference, "[[I": JavaValueReference,
	}
	for descriptor, want := range tests {
		if got, ok := JavaValueKindForDescriptor(descriptor); !ok || got != want {
			t.Errorf("JavaValueKindForDescriptor(%q) = (%s, %v), want (%s, true)", descriptor, got, ok, want)
		}
	}
	for _, descriptor := range []string{"", "Q", "Lmissing", "L;", "[", "[[", "[[V", "Lx;extra", "(I)V"} {
		if _, ok := JavaValueKindForDescriptor(descriptor); ok {
			t.Errorf("JavaValueKindForDescriptor(%q) unexpectedly succeeded", descriptor)
		}
	}
}

func TestTypedAdaptersRoundTripAllNonVoidKinds(t *testing.T) {
	ref := &Object{}
	tests := []struct {
		name  string
		value JavaValue
		index int
	}{
		{"boolean", BooleanValue(true), 0},
		{"byte", ByteValue(-12), 1},
		{"char", CharValue(0xfffe), 2},
		{"short", ShortValue(-1234), 3},
		{"int", IntValue(-1234567), 4},
		{"long", LongValue(math.MinInt64 + 9), 5},
		{"float", FloatValue(math.Float32frombits(0x7fc01234)), 7},
		{"double", DoubleValue(math.Float64frombits(0x7ff8000000001234)), 8},
		{"reference", ReferenceValue(ref), 10},
		{"null reference", ReferenceValue(nil), 11},
	}
	frame := NewFrame(nil, &Method{maxLocals: 12})
	object := &Object{heapCells: NewHeapCells(12)}
	for _, test := range tests {
		t.Run(test.name+" frame", func(t *testing.T) {
			if result := WriteFrameLocal(frame, test.index, test.value); !result.IsNormal() {
				t.Fatalf("WriteFrameLocal failed: %v", result.Failure())
			}
			assertNormalValue(t, ReadFrameLocal(frame, test.index, test.value.Kind()), test.value)
		})
		t.Run(test.name+" heap", func(t *testing.T) {
			if result := WriteHeapValue(object, test.index, test.value); !result.IsNormal() {
				t.Fatalf("WriteHeapValue failed: %v", result.Failure())
			}
			assertNormalValue(t, ReadHeapValue(object, test.index, test.value.Kind()), test.value)
		})
	}
}

func TestTypedAdaptersContainLegacyPanics(t *testing.T) {
	frame := NewFrame(nil, &Method{maxLocals: 1})
	if result := ReadFrameLocal(frame, 1, JavaValueInt); !result.IsInternalFailure() || result.Failure() == nil {
		t.Fatal("out-of-bounds frame read must become an internal failure")
	}
	if result := WriteHeapValue(nil, 0, IntValue(1)); !result.IsInternalFailure() || result.Failure() == nil {
		t.Fatal("nil heap object must become an internal failure")
	}
	if result := ReadHeapValue(&Object{heapCells: NewHeapCells(1)}, 0, JavaValueVoid); !result.IsInternalFailure() {
		t.Fatal("void heap read must become an internal failure")
	}
}

func TestTypedFieldAccessRoundTripsAllKinds(t *testing.T) {
	ref := &Object{}
	tests := []struct {
		name  string
		value JavaValue
		index uint
	}{
		{"boolean", BooleanValue(true), 0},
		{"byte", ByteValue(-12), 1},
		{"char", CharValue(0xfffe), 2},
		{"short", ShortValue(-1234), 3},
		{"int", IntValue(-1234567), 4},
		{"long", LongValue(math.MinInt64 + 9), 5},
		{"float", FloatValue(math.Float32frombits(0x7fc01234)), 6},
		{"double", DoubleValue(math.Float64frombits(0x7ff8000000001234)), 7},
		{"reference", ReferenceValue(ref), 8},
		{"null reference", ReferenceValue(nil), 9},
	}
	owner := &Class{name: "test/Owner", instCellCount: 10, staticCells: NewHeapCells(10)}
	instance := NewObject(owner)
	for _, test := range tests {
		t.Run(test.name+" instance", func(t *testing.T) {
			field := NewField(owner, test.name, descriptorForValue(test.value), 0, false, test.index)
			if result := WriteFieldValue(field, instance, test.value); !result.IsNormal() {
				t.Fatalf("WriteFieldValue failed: %v", result.Failure())
			}
			assertNormalValue(t, ReadFieldValue(field, instance), test.value)
		})
		t.Run(test.name+" static", func(t *testing.T) {
			field := NewField(owner, test.name, descriptorForValue(test.value), 0, true, test.index)
			if result := WriteFieldValue(field, nil, test.value); !result.IsNormal() {
				t.Fatalf("WriteFieldValue failed: %v", result.Failure())
			}
			assertNormalValue(t, ReadFieldValue(field, nil), test.value)
		})
	}
}

func TestTypedFieldAccessRejectsInvalidInputs(t *testing.T) {
	owner := &Class{name: "test/Owner", instCellCount: 1}
	field := NewField(owner, "i", "I", 0, false, 0)
	if result := ReadFieldValue(field, nil); !result.IsInternalFailure() {
		t.Fatal("nil instance receiver must be rejected")
	}
	if result := WriteFieldValue(field, NewObject(owner), LongValue(1)); !result.IsInternalFailure() {
		t.Fatal("mismatched logical value kind must be rejected")
	}
	voidField := NewField(owner, "bad", "V", 0, false, 0)
	if result := ReadFieldValue(voidField, NewObject(owner)); !result.IsInternalFailure() {
		t.Fatal("void field descriptor must be rejected")
	}
}

func descriptorForValue(value JavaValue) string {
	switch value.Kind() {
	case JavaValueBoolean:
		return "Z"
	case JavaValueByte:
		return "B"
	case JavaValueChar:
		return "C"
	case JavaValueShort:
		return "S"
	case JavaValueInt:
		return "I"
	case JavaValueLong:
		return "J"
	case JavaValueFloat:
		return "F"
	case JavaValueDouble:
		return "D"
	case JavaValueReference:
		return "Ljava/lang/Object;"
	default:
		panic("no field descriptor for void")
	}
}

func assertNormalValue(t *testing.T, result DynamicResult, want JavaValue) {
	t.Helper()
	got, ok := result.Value()
	if !ok {
		t.Fatalf("adapter did not return normal value: kind=%d failure=%v", result.Kind(), result.Failure())
	}
	if got != want {
		t.Fatalf("adapter round trip = %#v, want %#v", got, want)
	}
}
