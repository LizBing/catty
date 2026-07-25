package rtda

import (
	"fmt"
	"math"
)

func int32FromFloat32(v float32) int32 { return int32(math.Float32bits(v)) }
func float32FromInt32(v int32) float32 { return math.Float32frombits(uint32(v)) }
func int64FromFloat64(v float64) int64 { return int64(math.Float64bits(v)) }
func float64FromInt64(v int64) float64 { return math.Float64frombits(uint64(v)) }

// ReadFrameLocal adapts a Frame local at index into one logical JavaValue.
// It contains legacy frame panics so they cannot cross the typed boundary.
func ReadFrameLocal(frame *Frame, index int, kind JavaValueKind) (result DynamicResult) {
	defer containAdapterPanic(&result)
	if frame == nil {
		return InternalFailureResult(fmt.Errorf("read frame local: nil frame"))
	}
	switch kind {
	case JavaValueBoolean:
		return NormalResult(BooleanValue(frame.GetInt(index) != 0))
	case JavaValueByte:
		return NormalResult(ByteValue(int8(frame.GetInt(index))))
	case JavaValueChar:
		return NormalResult(CharValue(uint16(frame.GetInt(index))))
	case JavaValueShort:
		return NormalResult(ShortValue(int16(frame.GetInt(index))))
	case JavaValueInt:
		return NormalResult(IntValue(frame.GetInt(index)))
	case JavaValueLong:
		return NormalResult(LongValue(frame.GetLong(index)))
	case JavaValueFloat:
		return NormalResult(FloatValue(frame.GetFloat(index)))
	case JavaValueDouble:
		return NormalResult(DoubleValue(frame.GetDouble(index)))
	case JavaValueReference:
		return NormalResult(ReferenceValue(frame.GetRef(index)))
	default:
		return InternalFailureResult(fmt.Errorf("read frame local: invalid kind %s", kind))
	}
}

// WriteFrameLocal adapts a logical JavaValue into a Frame local. It accepts
// only an exact logical type: conversion, boxing, and widening belong to later
// profile-specific layers.
func WriteFrameLocal(frame *Frame, index int, value JavaValue) (result DynamicResult) {
	defer containAdapterPanic(&result)
	if frame == nil {
		return InternalFailureResult(fmt.Errorf("write frame local: nil frame"))
	}
	switch value.Kind() {
	case JavaValueBoolean:
		v, _ := value.Boolean()
		if v {
			frame.SetInt(index, 1)
		} else {
			frame.SetInt(index, 0)
		}
	case JavaValueByte:
		v, _ := value.Byte()
		frame.SetInt(index, int32(v))
	case JavaValueChar:
		v, _ := value.Char()
		frame.SetInt(index, int32(v))
	case JavaValueShort:
		v, _ := value.Short()
		frame.SetInt(index, int32(v))
	case JavaValueInt:
		v, _ := value.Int()
		frame.SetInt(index, v)
	case JavaValueLong:
		v, _ := value.Long()
		frame.SetLong(index, v)
	case JavaValueFloat:
		v, _ := value.Float()
		frame.SetFloat(index, v)
	case JavaValueDouble:
		v, _ := value.Double()
		frame.SetDouble(index, v)
	case JavaValueReference:
		v, _ := value.Reference()
		frame.SetRef(index, v)
	default:
		return InternalFailureResult(fmt.Errorf("write frame local: invalid kind %s", value.Kind()))
	}
	return NormalResult(VoidValue())
}

// ReadHeapValue adapts one typed heap cell. Heap storage remains private to
// Object; this function preserves the ADR-0030 typed-access boundary.
func ReadHeapValue(object *Object, index int, kind JavaValueKind) (result DynamicResult) {
	defer containAdapterPanic(&result)
	if object == nil {
		return InternalFailureResult(fmt.Errorf("read heap value: nil object"))
	}
	switch kind {
	case JavaValueBoolean:
		return NormalResult(BooleanValue(object.GetIntCell(index) != 0))
	case JavaValueByte:
		return NormalResult(ByteValue(int8(object.GetIntCell(index))))
	case JavaValueChar:
		return NormalResult(CharValue(uint16(object.GetIntCell(index))))
	case JavaValueShort:
		return NormalResult(ShortValue(int16(object.GetIntCell(index))))
	case JavaValueInt:
		return NormalResult(IntValue(object.GetIntCell(index)))
	case JavaValueLong:
		return NormalResult(LongValue(object.GetLongCell(index)))
	case JavaValueFloat:
		return NormalResult(FloatValue(object.GetFloatCell(index)))
	case JavaValueDouble:
		return NormalResult(DoubleValue(object.GetDoubleCell(index)))
	case JavaValueReference:
		return NormalResult(ReferenceValue(object.GetRefCell(index)))
	default:
		return InternalFailureResult(fmt.Errorf("read heap value: invalid kind %s", kind))
	}
}

// WriteHeapValue adapts one logical JavaValue into a typed heap cell.
func WriteHeapValue(object *Object, index int, value JavaValue) (result DynamicResult) {
	defer containAdapterPanic(&result)
	if object == nil {
		return InternalFailureResult(fmt.Errorf("write heap value: nil object"))
	}
	switch value.Kind() {
	case JavaValueBoolean:
		v, _ := value.Boolean()
		if v {
			object.SetIntCell(index, 1)
		} else {
			object.SetIntCell(index, 0)
		}
	case JavaValueByte:
		v, _ := value.Byte()
		object.SetIntCell(index, int32(v))
	case JavaValueChar:
		v, _ := value.Char()
		object.SetIntCell(index, int32(v))
	case JavaValueShort:
		v, _ := value.Short()
		object.SetIntCell(index, int32(v))
	case JavaValueInt:
		v, _ := value.Int()
		object.SetIntCell(index, v)
	case JavaValueLong:
		v, _ := value.Long()
		object.SetLongCell(index, v)
	case JavaValueFloat:
		v, _ := value.Float()
		object.SetFloatCell(index, v)
	case JavaValueDouble:
		v, _ := value.Double()
		object.SetDoubleCell(index, v)
	case JavaValueReference:
		v, _ := value.Reference()
		object.SetRefCell(index, v)
	default:
		return InternalFailureResult(fmt.Errorf("write heap value: invalid kind %s", value.Kind()))
	}
	return NormalResult(VoidValue())
}

func containAdapterPanic(result *DynamicResult) {
	if recovered := recover(); recovered != nil {
		*result = InternalFailureResult(fmt.Errorf("typed adapter panic: %v", recovered))
	}
}
