package rtda

import "fmt"

// ReadFieldValue reads an already-resolved field through the typed invocation
// boundary. For instance fields receiver must be the target object; static
// fields ignore it. Class initialization is intentionally not triggered here:
// the direct invocation service that owns an ExecutionContext performs that
// JVMS §5.5 transition in a later slice.
func ReadFieldValue(field *Field, receiver *Object) (result DynamicResult) {
	defer containAdapterPanic(&result)
	kind, failure := fieldValueKind(field)
	if failure != nil {
		return InternalFailureResult(failure)
	}
	if field.IsStatic() {
		return readStaticFieldValue(field, kind)
	}
	if failure := validateFieldReceiver(field, receiver); failure != nil {
		return InternalFailureResult(failure)
	}
	return ReadHeapValue(receiver, int(field.SlotID()), kind)
}

// WriteFieldValue writes an already-resolved field through the typed invocation
// boundary. It permits no Java conversion, boxing, or widening: value must
// exactly match the field descriptor's logical JavaValueKind.
func WriteFieldValue(field *Field, receiver *Object, value JavaValue) (result DynamicResult) {
	defer containAdapterPanic(&result)
	kind, failure := fieldValueKind(field)
	if failure != nil {
		return InternalFailureResult(failure)
	}
	if value.Kind() != kind {
		return InternalFailureResult(fmt.Errorf("write field %s.%s%s: value kind %s, want %s",
			field.Owner().Name(), field.Name(), field.Descriptor(), value.Kind(), kind))
	}
	if field.IsStatic() {
		return writeStaticFieldValue(field, value)
	}
	if failure := validateFieldReceiver(field, receiver); failure != nil {
		return InternalFailureResult(failure)
	}
	return WriteHeapValue(receiver, int(field.SlotID()), value)
}

func fieldValueKind(field *Field) (JavaValueKind, error) {
	if field == nil {
		return 0, fmt.Errorf("typed field access: nil field")
	}
	if field.Owner() == nil {
		return 0, fmt.Errorf("typed field access %s%s: nil owner", field.Name(), field.Descriptor())
	}
	kind, ok := JavaValueKindForDescriptor(field.Descriptor())
	if !ok || kind == JavaValueVoid {
		return 0, fmt.Errorf("typed field access %s.%s%s: invalid field descriptor",
			field.Owner().Name(), field.Name(), field.Descriptor())
	}
	return kind, nil
}

func validateFieldReceiver(field *Field, receiver *Object) error {
	if receiver == nil {
		return fmt.Errorf("typed field access %s.%s%s: nil receiver",
			field.Owner().Name(), field.Name(), field.Descriptor())
	}
	if !receiver.IsInstanceOf(field.Owner()) {
		return fmt.Errorf("typed field access %s.%s%s: receiver is not assignable to owner",
			field.Owner().Name(), field.Name(), field.Descriptor())
	}
	return nil
}

func readStaticFieldValue(field *Field, kind JavaValueKind) DynamicResult {
	owner, slot := field.Owner(), field.SlotID()
	switch kind {
	case JavaValueBoolean:
		return NormalResult(BooleanValue(owner.GetStaticInt(slot) != 0))
	case JavaValueByte:
		return NormalResult(ByteValue(int8(owner.GetStaticInt(slot))))
	case JavaValueChar:
		return NormalResult(CharValue(uint16(owner.GetStaticInt(slot))))
	case JavaValueShort:
		return NormalResult(ShortValue(int16(owner.GetStaticInt(slot))))
	case JavaValueInt:
		return NormalResult(IntValue(owner.GetStaticInt(slot)))
	case JavaValueLong:
		return NormalResult(LongValue(owner.GetStaticLong(slot)))
	case JavaValueFloat:
		return NormalResult(FloatValue(owner.GetStaticFloat(slot)))
	case JavaValueDouble:
		return NormalResult(DoubleValue(owner.GetStaticDouble(slot)))
	case JavaValueReference:
		return NormalResult(ReferenceValue(owner.GetStaticRef(slot)))
	default:
		return InternalFailureResult(fmt.Errorf("read static field: invalid kind %s", kind))
	}
}

func writeStaticFieldValue(field *Field, value JavaValue) DynamicResult {
	owner, slot := field.Owner(), field.SlotID()
	switch value.Kind() {
	case JavaValueBoolean:
		v, _ := value.Boolean()
		if v {
			owner.SetStaticInt(slot, 1)
		} else {
			owner.SetStaticInt(slot, 0)
		}
	case JavaValueByte:
		v, _ := value.Byte()
		owner.SetStaticInt(slot, int32(v))
	case JavaValueChar:
		v, _ := value.Char()
		owner.SetStaticInt(slot, int32(v))
	case JavaValueShort:
		v, _ := value.Short()
		owner.SetStaticInt(slot, int32(v))
	case JavaValueInt:
		v, _ := value.Int()
		owner.SetStaticInt(slot, v)
	case JavaValueLong:
		v, _ := value.Long()
		owner.SetStaticLong(slot, v)
	case JavaValueFloat:
		v, _ := value.Float()
		owner.SetStaticFloat(slot, v)
	case JavaValueDouble:
		v, _ := value.Double()
		owner.SetStaticDouble(slot, v)
	case JavaValueReference:
		v, _ := value.Reference()
		owner.SetStaticRef(slot, v)
	default:
		return InternalFailureResult(fmt.Errorf("write static field: invalid kind %s", value.Kind()))
	}
	return NormalResult(VoidValue())
}
