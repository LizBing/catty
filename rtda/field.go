package rtda

// Field is the runtime representation of a class field (JVMS §4.5). Every field
// is assigned a slotID — the offset into an Object's fields slice (instance
// fields) or the class's staticVars slice (static fields). Slot IDs are
// contiguous per class, inheriting the superclass's count for instances.
type Field struct {
	owner       *Class
	name        string
	descriptor  string
	accessFlags uint16
	slotID      uint
	isStatic    bool
}

func NewField(owner *Class, name, descriptor string, access uint16, isStatic bool, slotID uint) *Field {
	return &Field{
		owner:       owner,
		name:        name,
		descriptor:  descriptor,
		accessFlags: access,
		slotID:      slotID,
		isStatic:    isStatic,
	}
}

func (f *Field) Owner() *Class    { return f.owner }
func (f *Field) Name() string     { return f.name }
func (f *Field) Descriptor() string { return f.descriptor }
func (f *Field) SlotID() uint     { return f.slotID }
func (f *Field) IsStatic() bool   { return f.isStatic }
func (f *Field) IsLongOrDouble() bool {
	return f.descriptor == "J" || f.descriptor == "D"
}
