package rtda

import (
	"fmt"
)

// JavaValueKind identifies one logical Java value category. It deliberately
// describes JVMS values, rather than an interpreter slot, IR register, or heap
// representation. Long and double are each one JavaValue even though a Frame
// uses two local/operand slots for them.
type JavaValueKind uint8

const (
	JavaValueVoid JavaValueKind = iota
	JavaValueBoolean
	JavaValueByte
	JavaValueChar
	JavaValueShort
	JavaValueInt
	JavaValueLong
	JavaValueFloat
	JavaValueDouble
	JavaValueReference
)

func (k JavaValueKind) String() string {
	switch k {
	case JavaValueVoid:
		return "void"
	case JavaValueBoolean:
		return "boolean"
	case JavaValueByte:
		return "byte"
	case JavaValueChar:
		return "char"
	case JavaValueShort:
		return "short"
	case JavaValueInt:
		return "int"
	case JavaValueLong:
		return "long"
	case JavaValueFloat:
		return "float"
	case JavaValueDouble:
		return "double"
	case JavaValueReference:
		return "reference"
	default:
		return fmt.Sprintf("JavaValueKind(%d)", k)
	}
}

// JavaValue is an engine-neutral logical Java value. Its payload is private so
// callers cannot accidentally substitute a Slot, IR register, or heap bits for
// the dynamic-invocation boundary.
type JavaValue struct {
	kind JavaValueKind
	bits int64
	ref  *Object
}

func VoidValue() JavaValue { return JavaValue{kind: JavaValueVoid} }
func BooleanValue(v bool) JavaValue {
	if v {
		return JavaValue{kind: JavaValueBoolean, bits: 1}
	}
	return JavaValue{kind: JavaValueBoolean}
}
func ByteValue(v int8) JavaValue   { return JavaValue{kind: JavaValueByte, bits: int64(v)} }
func CharValue(v uint16) JavaValue { return JavaValue{kind: JavaValueChar, bits: int64(v)} }
func ShortValue(v int16) JavaValue { return JavaValue{kind: JavaValueShort, bits: int64(v)} }
func IntValue(v int32) JavaValue   { return JavaValue{kind: JavaValueInt, bits: int64(v)} }
func LongValue(v int64) JavaValue  { return JavaValue{kind: JavaValueLong, bits: v} }
func FloatValue(v float32) JavaValue {
	return JavaValue{kind: JavaValueFloat, bits: int64(int32FromFloat32(v))}
}
func DoubleValue(v float64) JavaValue {
	return JavaValue{kind: JavaValueDouble, bits: int64FromFloat64(v)}
}
func ReferenceValue(v *Object) JavaValue { return JavaValue{kind: JavaValueReference, ref: v} }

// Kind returns the logical Java category of v.
func (v JavaValue) Kind() JavaValueKind { return v.kind }

func (v JavaValue) Boolean() (bool, bool) { return v.bits != 0, v.kind == JavaValueBoolean }
func (v JavaValue) Byte() (int8, bool)    { return int8(v.bits), v.kind == JavaValueByte }
func (v JavaValue) Char() (uint16, bool)  { return uint16(v.bits), v.kind == JavaValueChar }
func (v JavaValue) Short() (int16, bool)  { return int16(v.bits), v.kind == JavaValueShort }
func (v JavaValue) Int() (int32, bool)    { return int32(v.bits), v.kind == JavaValueInt }
func (v JavaValue) Long() (int64, bool)   { return v.bits, v.kind == JavaValueLong }
func (v JavaValue) Float() (float32, bool) {
	return float32FromInt32(int32(v.bits)), v.kind == JavaValueFloat
}
func (v JavaValue) Double() (float64, bool) {
	return float64FromInt64(v.bits), v.kind == JavaValueDouble
}
func (v JavaValue) Reference() (*Object, bool) { return v.ref, v.kind == JavaValueReference }
func (v JavaValue) IsVoid() bool               { return v.kind == JavaValueVoid }

// DynamicResultKind distinguishes a normal Java result (including void and
// null), a Java throwable, and a Catty-internal adapter failure.
type DynamicResultKind uint8

const (
	DynamicResultNormal DynamicResultKind = iota
	DynamicResultThrowable
	DynamicResultInternalFailure
)

// DynamicResult transports an invocation outcome without using Go panic as a
// Java-visible failure mechanism. Throwable identity is preserved exactly.
type DynamicResult struct {
	kind      DynamicResultKind
	value     JavaValue
	throwable *Object
	failure   error
}

func NormalResult(value JavaValue) DynamicResult {
	return DynamicResult{kind: DynamicResultNormal, value: value}
}
func ThrowableResult(throwable *Object) DynamicResult {
	return DynamicResult{kind: DynamicResultThrowable, throwable: throwable}
}
func InternalFailureResult(err error) DynamicResult {
	return DynamicResult{kind: DynamicResultInternalFailure, failure: err}
}

func (r DynamicResult) Kind() DynamicResultKind    { return r.kind }
func (r DynamicResult) IsNormal() bool             { return r.kind == DynamicResultNormal }
func (r DynamicResult) IsThrowable() bool          { return r.kind == DynamicResultThrowable }
func (r DynamicResult) IsInternalFailure() bool    { return r.kind == DynamicResultInternalFailure }
func (r DynamicResult) Value() (JavaValue, bool)   { return r.value, r.IsNormal() }
func (r DynamicResult) Throwable() (*Object, bool) { return r.throwable, r.IsThrowable() }
func (r DynamicResult) Failure() error             { return r.failure }

// JavaValueKindForDescriptor maps one field/return descriptor to its logical
// value kind. Arrays and object descriptors are references. Void is accepted
// only for method returns; consumers decide where it is valid.
func JavaValueKindForDescriptor(descriptor string) (JavaValueKind, bool) {
	switch descriptor {
	case "V":
		return JavaValueVoid, true
	case "Z":
		return JavaValueBoolean, true
	case "B":
		return JavaValueByte, true
	case "C":
		return JavaValueChar, true
	case "S":
		return JavaValueShort, true
	case "I":
		return JavaValueInt, true
	case "J":
		return JavaValueLong, true
	case "F":
		return JavaValueFloat, true
	case "D":
		return JavaValueDouble, true
	}
	if isReferenceDescriptor(descriptor) {
		return JavaValueReference, true
	}
	return 0, false
}

func isReferenceDescriptor(descriptor string) bool {
	i := 0
	for i < len(descriptor) && descriptor[i] == '[' {
		i++
	}
	if i == len(descriptor) {
		return false
	}
	switch descriptor[i] {
	case 'Z', 'B', 'C', 'S', 'I', 'J', 'F', 'D':
		return i+1 == len(descriptor)
	case 'L':
		if i+2 >= len(descriptor) || descriptor[len(descriptor)-1] != ';' {
			return false
		}
		for _, ch := range descriptor[i+1 : len(descriptor)-1] {
			if ch == '.' || ch == ';' || ch == '[' {
				return false
			}
		}
		return true
	default:
		return false
	}
}
