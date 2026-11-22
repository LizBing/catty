package rt

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"
)

// Runtime owns loader state, class metadata, static storage, mirrors, strings
// and the single logical Java thread used by the launcher.
type Runtime struct {
	Image     *Image
	Classes   map[string]*Class
	Bootstrap map[string]*Class
	Loader    *Loader

	primitives   map[byte]*Class
	arrayClasses map[string]*Class
	strings      map[string]*Object
	boxedInts    map[int32]*Object
	boxedLongs   map[int64]*Object
	boxedShorts  map[int16]*Object
	boxedBytes   map[int8]*Object
	boxedChars   map[uint16]*Object
	boxedBools   map[bool]*Object
	classOrder   []*Class
	primitiveOrder []*Class

	resolutions map[string]*Resolution

	nextIdentity uint64
	nextThread   int64
	installed    bool
	trace        *trace
	mainThread   *Thread

	// AllocHook is an optional testing hook invoked after each instance
	// allocation. It must not retain the object; the garbage-collection
	// acceptance test uses it to install a finalization probe and record a
	// non-rooting identity witness.
	AllocHook func(o *Object)
}

// Resolution is one cached symbolic-resolution outcome. A non-nil Thrown is a
// persistent failure outcome rethrown by later uses of the same symbol.
type Resolution struct {
	Method *Method
	Field  *Field
	Thrown *Object
}

// NewRuntime creates an uninstalled runtime for img.
func NewRuntime(img *Image) *Runtime {
	r := &Runtime{
		Image:        img,
		Classes:      map[string]*Class{},
		Bootstrap:    map[string]*Class{},
		primitives:   map[byte]*Class{},
		arrayClasses: map[string]*Class{},
		strings:      map[string]*Object{},
		boxedInts:    map[int32]*Object{},
		boxedLongs:   map[int64]*Object{},
		boxedShorts:  map[int16]*Object{},
		boxedBytes:   map[int8]*Object{},
		boxedChars:   map[uint16]*Object{},
		boxedBools:   map[bool]*Object{},
		resolutions:  map[string]*Resolution{},
	}
	r.Loader = &Loader{Name: "app", Runtime: r}
	return r
}

// Installed reports whether Install completed.
func (r *Runtime) Installed() bool { return r.installed }

func supportedOpcodeSet() map[string]bool {
	out := map[string]bool{}
	for _, name := range SupportedOpcodeNames() {
		out[name] = true
	}
	return out
}

// Install validates the image ABI, profile and opcode inventory and
// materializes immutable class definitions, static storage, mirrors and method
// metadata. No Java initializer runs here.
func (r *Runtime) Install() error {
	if r.installed {
		return fmt.Errorf("catty: runtime already installed")
	}
	if r.Image == nil {
		return fmt.Errorf("catty: nil image")
	}
	if r.Image.ABI != ABIRevision {
		return fmt.Errorf("catty: unsupported image ABI %d (runtime ABI %d)", r.Image.ABI, ABIRevision)
	}
	if r.Image.Profile != ProfileID {
		return fmt.Errorf("catty: unsupported image profile %q (runtime profile %q)", r.Image.Profile, ProfileID)
	}
	if len(r.Image.Opcodes) > 0 {
		supported := supportedOpcodeSet()
		for _, name := range r.Image.Opcodes {
			if !supported[name] {
				return fmt.Errorf("catty: image opcode inventory contains unsupported opcode %q", name)
			}
		}
	}

	for _, bc := range bootstrapInventory {
		c := &Class{
			Name:        bc.name,
			Flags:       bc.flags | AccPublic,
			MethodIndex: map[string]*Method{},
			State:       classInitialized,
			Loader:      r.Loader,
		}
		r.Bootstrap[bc.name] = c
		r.classOrder = append(r.classOrder, c)
	}
	for _, pc := range PrimitiveClasses {
		c := &Class{
			Name:        pc.Name,
			Descriptor:  pc.Descriptor,
			IsPrimitive: true,
			Flags:       AccPublic | AccFinal,
			MethodIndex: map[string]*Method{},
			State:       classInitialized,
			Loader:      r.Loader,
		}
		r.primitives[pc.Descriptor[0]] = c
		r.classOrder = append(r.classOrder, c)
		r.primitiveOrder = append(r.primitiveOrder, c)
	}

	seen := map[string]bool{}
	for _, def := range r.Image.Classes {
		if def == nil || def.Name == "" {
			return fmt.Errorf("catty: image contains an unnamed class definition")
		}
		if strings.HasPrefix(def.Name, "[") {
			return fmt.Errorf("catty: image class %q has an array name", def.Name)
		}
		if _, ok := r.Bootstrap[def.Name]; ok || strings.HasPrefix(def.Name, "java/") {
			return fmt.Errorf("catty: image class %q collides with the reserved bootstrap namespace", def.Name)
		}
		if seen[def.Name] {
			return fmt.Errorf("catty: duplicate class definition %q in image", def.Name)
		}
		seen[def.Name] = true
		c := &Class{
			Name:        def.Name,
			Flags:       def.Flags,
			MethodIndex: map[string]*Method{},
			NestHost:    def.NestHost,
			NestMembers: append([]string(nil), def.NestMembers...),
			Loader:      r.Loader,
		}
		r.Classes[def.Name] = c
	}

	// Link bootstrap hierarchies first.
	for _, bc := range bootstrapInventory {
		c := r.Bootstrap[bc.name]
		if bc.super != "" {
			s, ok := r.Bootstrap[bc.super]
			if !ok {
				return fmt.Errorf("catty: bootstrap class %s has unknown superclass %s", bc.name, bc.super)
			}
			c.Super = s
		}
		for _, in := range bc.interfaces {
			i, ok := r.Bootstrap[in]
			if !ok {
				return fmt.Errorf("catty: bootstrap class %s has unknown interface %s", bc.name, in)
			}
			c.Interfaces = append(c.Interfaces, i)
		}
	}

	// Link application hierarchies.
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		if def.Super != "" {
			s, ok := r.lookupClassShell(def.Super)
			if !ok {
				return fmt.Errorf("catty: class %s has unknown superclass %s", def.Name, def.Super)
			}
			c.Super = s
		}
		for _, in := range def.Interfaces {
			i, ok := r.lookupClassShell(in)
			if !ok {
				return fmt.Errorf("catty: class %s has unknown interface %s", def.Name, in)
			}
			c.Interfaces = append(c.Interfaces, i)
		}
	}
	if err := r.checkHierarchies(); err != nil {
		return err
	}

	// Declared fields and static storage. Layout is assigned superclass-first.
	layoutState := map[*Class]uint8{}
	var layout func(c *Class) error
	layout = func(c *Class) error {
		switch layoutState[c] {
		case 1:
			return fmt.Errorf("catty: class hierarchy cycle at %s", c.Name)
		case 2:
			return nil
		}
		layoutState[c] = 1
		if c.Super != nil {
			if err := layout(c.Super); err != nil {
				return err
			}
		}
		size := 0
		if c.Super != nil {
			size = c.Super.InstanceSize
		}
		c.InstanceField = nil
		if c.Super != nil {
			c.InstanceKinds = append([]string(nil), c.Super.InstanceKinds...)
		} else {
			c.InstanceKinds = nil
		}
		c.StaticField = nil
		if bc, ok := bootstrapByName[c.Name]; ok {
			for _, f := range bc.fields {
				ft, err := ParseFieldType(f.descriptor)
				if err != nil {
					return fmt.Errorf("catty: bootstrap field %s.%s: %w", c.Name, f.name, err)
				}
				c.StaticField = append(c.StaticField, &Field{
					Owner: c, Name: f.name, Descriptor: f.descriptor, Flags: f.flags,
					Static: true, Index: len(c.StaticField), Type: ft,
				})
			}
			c.InstanceSize = 0
			c.Statics = make([]Value, len(c.StaticField))
			layoutState[c] = 2
			return nil
		}
		if c.IsPrimitive {
			c.Statics = nil
			c.InstanceSize = 0
			layoutState[c] = 2
			return nil
		}
		for _, def := range r.Image.Classes {
			if def.Name != c.Name {
				continue
			}
			for _, f := range def.InstanceFields {
				ft, err := ParseFieldType(f.Descriptor)
				if err != nil {
					return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
				}
				c.InstanceField = append(c.InstanceField, &Field{
					Owner: c, Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags,
					Offset: size, Type: ft,
				})
				c.InstanceKinds = append(c.InstanceKinds, f.Descriptor)
				size++
			}
			c.InstanceSize = size
			for _, f := range def.StaticFields {
				ft, err := ParseFieldType(f.Descriptor)
				if err != nil {
					return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
				}
				c.StaticField = append(c.StaticField, &Field{
					Owner: c, Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags,
					Static: true, Index: len(c.StaticField), Type: ft, Constant: f.Constant,
				})
			}
			c.Statics = make([]Value, len(c.StaticField))
			layoutState[c] = 2
			return nil
		}
		return fmt.Errorf("catty: missing image definition for %s", c.Name)
	}
	for _, def := range r.Image.Classes {
		if err := layout(r.Classes[def.Name]); err != nil {
			return err
		}
	}
	for _, bc := range bootstrapInventory {
		if err := layout(r.Bootstrap[bc.name]); err != nil {
			return err
		}
	}
	for _, c := range r.primitiveOrder {
		if err := layout(c); err != nil {
			return err
		}
	}

	// Methods.
	for _, bc := range bootstrapInventory {
		c := r.Bootstrap[bc.name]
		for _, bm := range bc.methods {
			mt, err := ParseMethodType(bm.descriptor)
			if err != nil {
				return fmt.Errorf("catty: bootstrap method %s.%s: %w", c.Name, bm.name, err)
			}
			m := &Method{
				Owner: c, Name: bm.name, Descriptor: bm.descriptor, Flags: bm.flags,
				Engine: "intrinsic", Reason: "bootstrap intrinsic", Intrinsic: bm.intrinsic, Type: mt,
			}
			if err := c.addMethod(m); err != nil {
				return err
			}
		}
	}
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		for _, md := range def.Methods {
			if md == nil {
				return fmt.Errorf("catty: class %s contains a nil method", c.Name)
			}
			if md.Owner != "" && md.Owner != c.Name {
				return fmt.Errorf("catty: method %s.%s%s declares owner %s", c.Name, md.Name, md.Descriptor, md.Owner)
			}
			mt, err := ParseMethodType(md.Descriptor)
			if err != nil {
				return fmt.Errorf("catty: method %s.%s: %w", c.Name, md.Name, err)
			}
			m := &Method{
				Owner: c, Name: md.Name, Descriptor: md.Descriptor, Flags: md.Flags,
				Engine: md.Engine, Reason: md.Reason, MaxStack: md.MaxStack, MaxLocals: md.MaxLocals,
				CodeSize: md.CodeSize, Refs: md.Refs, Insns: md.Insns, Handlers: md.Handlers,
				Entry: md.Entry, Type: mt,
			}
			switch md.Engine {
			case "aot":
				if m.Entry == nil {
					return fmt.Errorf("catty: AOT method %s.%s%s has no compiled entry", c.Name, md.Name, md.Descriptor)
				}
				if len(m.Insns) != 0 || m.CodeSize != 0 {
					return fmt.Errorf("catty: AOT method %s.%s%s also retains an interpreter body", c.Name, md.Name, md.Descriptor)
				}
			case "interp":
				if len(m.Insns) == 0 {
					return fmt.Errorf("catty: interpreted method %s.%s%s has no instructions", c.Name, md.Name, md.Descriptor)
				}
				if m.CodeSize <= 0 {
					return fmt.Errorf("catty: interpreted method %s.%s%s has no retained code size", c.Name, md.Name, md.Descriptor)
				}
			case "abstract":
				if m.Flags&AccAbstract == 0 {
					return fmt.Errorf("catty: method %s.%s%s is marked abstract in the image without ACC_ABSTRACT", c.Name, md.Name, md.Descriptor)
				}
			default:
				return fmt.Errorf("catty: method %s.%s%s has unknown engine %q", c.Name, md.Name, md.Descriptor, md.Engine)
			}
			for i := range m.Insns {
				if m.Insns[i].Ref >= int32(len(m.Refs)) {
					return fmt.Errorf("catty: method %s.%s%s instruction pc %d has invalid ref index %d",
						c.Name, md.Name, md.Descriptor, m.Insns[i].PC, m.Insns[i].Ref)
				}
			}
			if m.Engine == "interp" {
				m.pcIndex = map[int32]int{}
				for i := range m.Insns {
					m.pcIndex[m.Insns[i].PC] = i
				}
			}
			if err := c.addMethod(m); err != nil {
				return err
			}
		}
	}

	// Resolve parameter classes now that all class shells, primitives and
	// bootstrap classes exist.
	for _, c := range r.allClasses() {
		for _, m := range c.Methods {
			m.ParamClass = make([]*Class, len(m.Type.Params))
			for i, p := range m.Type.Params {
				pc, err := r.classForType(p)
				if err != nil {
					return fmt.Errorf("catty: method %s.%s%s parameter %d: %w", c.Name, m.Name, m.Descriptor, i, err)
				}
				m.ParamClass[i] = pc
			}
		}
		for _, f := range c.StaticField {
			if _, err := r.classForType(f.Type); err != nil {
				return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
			}
		}
		for _, f := range c.InstanceField {
			if _, err := r.classForType(f.Type); err != nil {
				return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
			}
		}
	}

	// Canonical mirrors.
	classClass := r.mustClass("java/lang/Class")
	for _, c := range r.classOrder {
		c.Mirror = &Object{Class: classClass, Payload: c, Identity: r.newIdentity()}
	}
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		c.Mirror = &Object{Class: classClass, Payload: c, Identity: r.newIdentity()}
	}

	// Bootstrap static storage.
	for _, bc := range bootstrapInventory {
		c := r.Bootstrap[bc.name]
		for i, f := range bc.fields {
			if f.initial != nil {
				c.Statics[i] = f.initial(r)
			} else {
				c.Statics[i] = zeroValueFor(f.descriptor)
			}
		}
	}
	// Application static storage: preparation initializes each slot to its type
	// default. ConstantValue attributes are applied by the initialization
	// procedure of JVMS 5.5, not during installation.
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		for i, f := range c.StaticField {
			if i < len(c.Statics) {
				c.Statics[i] = zeroValueFor(f.Descriptor)
			}
		}
	}

	r.installed = true
	return nil
}

func (r *Runtime) constantValue(c *Constant) (Value, error) {
	switch c.Kind {
	case "int":
		return IntValue(int32(uint32(c.Bits))), nil
	case "long":
		return LongValue(int64(c.Bits)), nil
	case "float":
		return FloatBitsValue(uint32(c.Bits)), nil
	case "double":
		return DoubleBitsValue(c.Bits), nil
	case "string":
		return RefValue(r.NewStringUnits(c.Units, true)), nil
	}
	return Value{}, fmt.Errorf("unsupported ConstantValue kind %q", c.Kind)
}

func zeroValueFor(desc string) Value {
	if desc == "" {
		return RefValue(nil)
	}
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		return IntValue(0)
	case 'J':
		return LongValue(0)
	case 'F':
		return FloatValue(0)
	case 'D':
		return DoubleValue(0)
	}
	return RefValue(nil)
}

func (r *Runtime) lookupClassShell(name string) (*Class, bool) {
	if c, ok := r.Classes[name]; ok {
		return c, true
	}
	if c, ok := r.Bootstrap[name]; ok {
		return c, true
	}
	if desc, ok := PrimitiveDescriptor(name); ok {
		if c, ok := r.primitives[desc]; ok {
			return c, true
		}
	}
	return nil, false
}

func (r *Runtime) allClasses() []*Class {
	out := make([]*Class, 0, len(r.classOrder)+len(r.Classes))
	out = append(out, r.classOrder...)
	names := make([]string, 0, len(r.Classes))
	for n := range r.Classes {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		out = append(out, r.Classes[n])
	}
	return out
}

func (r *Runtime) checkHierarchies() error {
	const (
		white = 0
		grey  = 1
		black = 2
	)
	state := map[*Class]uint8{}
	var visit func(c *Class) error
	visit = func(c *Class) error {
		switch state[c] {
		case grey:
			return fmt.Errorf("catty: class hierarchy cycle at %s", c.Name)
		case black:
			return nil
		}
		state[c] = grey
		if c.Super != nil {
			if err := visit(c.Super); err != nil {
				return err
			}
			if c.Super.IsPrimitive || c.Super.IsArray {
				return fmt.Errorf("catty: class %s has invalid superclass %s", c.Name, c.Super.Name)
			}
		}
		for _, in := range c.Interfaces {
			if err := visit(in); err != nil {
				return err
			}
			if !in.IsInterface() {
				return fmt.Errorf("catty: class %s implements non-interface %s", c.Name, in.Name)
			}
		}
		state[c] = black
		return nil
	}
	for _, c := range r.classOrder {
		if c.IsPrimitive {
			continue
		}
		if err := visit(c); err != nil {
			return err
		}
	}
	for _, def := range r.Image.Classes {
		if err := visit(r.Classes[def.Name]); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) newIdentity() uint64 {
	r.nextIdentity++
	return r.nextIdentity
}

// Class returns an installed class by internal name.
func (r *Runtime) Class(name string) (*Class, error) {
	c, ok := r.lookupClassShell(name)
	if !ok {
		return nil, fmt.Errorf("catty: unknown class %s", name)
	}
	return c, nil
}

func (r *Runtime) mustClass(name string) *Class {
	c, ok := r.lookupClassShell(name)
	if !ok {
		panic("catty: missing required class " + name)
	}
	return c
}

// PrimitiveClass returns the canonical primitive class for a descriptor
// letter.
func (r *Runtime) PrimitiveClass(desc byte) *Class {
	c, ok := r.primitives[desc]
	if !ok {
		panic(fmt.Sprintf("catty: missing primitive class %q", desc))
	}
	return c
}

// PrimitiveMirror returns the canonical mirror of a primitive type.
func (r *Runtime) PrimitiveMirror(desc byte) *Object {
	return r.PrimitiveClass(desc).Mirror
}

// classForType resolves a parsed type to its runtime class.
func (r *Runtime) classForType(t Type) (*Class, error) {
	switch t.Kind {
	case 'V':
		return nil, nil
	case 'I', 'J', 'F', 'D':
		desc := t.Descriptor
		if desc == "" || len(desc) != 1 {
			return nil, fmt.Errorf("malformed primitive descriptor %q", desc)
		}
		c, ok := r.primitives[desc[0]]
		if !ok {
			return nil, fmt.Errorf("unknown primitive descriptor %q", desc)
		}
		return c, nil
	case 'L', 'A':
		return r.classForName(t.DescriptorOrName())
	default:
		return nil, fmt.Errorf("unsupported type %q", t.Descriptor)
	}
}

// DescriptorOrName returns an internal name for class/array types.
func (t Type) DescriptorOrName() string {
	if t.Kind == 'L' {
		return t.ClassName
	}
	return t.Descriptor
}

// classForName resolves an internal class name or array descriptor.
func (r *Runtime) classForName(name string) (*Class, error) {
	if strings.HasPrefix(name, "[") {
		return r.arrayClassFromDescriptor(name)
	}
	if c, ok := r.Classes[name]; ok {
		return c, nil
	}
	if c, ok := r.Bootstrap[name]; ok {
		return c, nil
	}
	if desc, ok := PrimitiveDescriptor(name); ok {
		if c, ok := r.primitives[desc]; ok {
			return c, nil
		}
	}
	return nil, fmt.Errorf("unknown class %s", name)
}

func (r *Runtime) arrayClassFromDescriptor(desc string) (*Class, error) {
	if c, ok := r.arrayClasses[desc]; ok {
		return c, nil
	}
	if !strings.HasPrefix(desc, "[") {
		return nil, fmt.Errorf("not an array descriptor: %s", desc)
	}
	if _, err := ParseFieldType(desc); err != nil {
		return nil, fmt.Errorf("malformed array descriptor %s: %w", desc, err)
	}
	componentDesc := desc[1:]
	var component *Class
	var err error
	switch componentDesc[0] {
	case '[':
		component, err = r.arrayClassFromDescriptor(componentDesc)
	case 'L':
		if !strings.HasSuffix(componentDesc, ";") {
			return nil, fmt.Errorf("malformed array descriptor: %s", desc)
		}
		component, err = r.classForName(componentDesc[1 : len(componentDesc)-1])
	default:
		primitive, ok := r.primitives[componentDesc[0]]
		if !ok {
			return nil, fmt.Errorf("unsupported array descriptor: %s", desc)
		}
		component = primitive
	}
	if err != nil {
		return nil, err
	}
	c := &Class{
		Name:       desc,
		Descriptor: desc,
		IsArray:    true,
		Component:  component,
		Super:      r.mustClass("java/lang/Object"),
		Interfaces: []*Class{r.mustClass("java/lang/Cloneable"), r.mustClass("java/io/Serializable")},
		Flags:      AccPublic | AccFinal,
		State:      classInitialized,
		Loader:     r.Loader,
	}
	classClass := r.mustClass("java/lang/Class")
	c.Mirror = &Object{Class: classClass, Payload: c, Identity: r.newIdentity()}
	r.arrayClasses[desc] = c
	return c, nil
}

// NewInstance allocates an instance of c with zero-initialized slots.
func (r *Runtime) NewInstance(c *Class) *Object {
	o := &Object{
		Class:    c,
		Fields:   make([]Value, c.InstanceSize),
		Identity: r.newIdentity(),
	}
	for i := range o.Fields {
		desc := ""
		if i < len(c.InstanceKinds) {
			desc = c.InstanceKinds[i]
		}
		o.Fields[i] = zeroValueFor(desc)
	}
	if r.AllocHook != nil {
		r.AllocHook(o)
	}
	return o
}

// NewArray allocates a reference array whose component is component.
func (r *Runtime) NewArray(component *Class, length int) *Object {
	ac, err := r.arrayClassFor(component)
	if err != nil {
		panic(err)
	}
	o := &Object{
		Class:    ac,
		Array:    &ArrayData{Component: component, Kind: 'L', Elems: make([]Value, length)},
		Identity: r.newIdentity(),
	}
	return o
}

// NewArrayOfClass allocates an array of the given array class. The element
// storage is chosen from the component type.
func (r *Runtime) NewArrayOfClass(arrayClass *Class, length int) *Object {
	component := arrayClass.Component
	a := &ArrayData{Component: component}
	if !component.IsPrimitive {
		a.Kind = 'L'
		a.Elems = make([]Value, length)
	} else {
		switch component.Descriptor[0] {
		case 'Z', 'B':
			a.Kind = component.Descriptor[0]
			a.Bytes = make([]int8, length)
		case 'C':
			a.Kind = 'C'
			a.Chars = make([]uint16, length)
		case 'S':
			a.Kind = 'S'
			a.Shorts = make([]int16, length)
		case 'I':
			a.Kind = 'I'
			a.Ints = make([]int32, length)
		case 'J':
			a.Kind = 'J'
			a.Longs = make([]int64, length)
		case 'F':
			a.Kind = 'F'
			a.Floats = make([]float32, length)
		case 'D':
			a.Kind = 'D'
			a.Doubles = make([]float64, length)
		default:
			panic("catty: unsupported array component " + component.Descriptor)
		}
	}
	return &Object{Class: arrayClass, Array: a, Identity: r.newIdentity()}
}

func (r *Runtime) arrayClassFor(component *Class) (*Class, error) {
	return r.arrayClassFromDescriptor(arrayDescriptorOf(component))
}

func arrayDescriptorOf(component *Class) string {
	if component.IsArray {
		return "[" + component.Descriptor
	}
	if component.IsPrimitive {
		return "[" + component.Descriptor
	}
	return "[L" + component.Name + ";"
}

// IsAssignable implements reference/array assignability.
func (r *Runtime) IsAssignable(from, to *Class) bool {
	if from == nil || to == nil {
		return from == nil && to == nil
	}
	if from == to {
		return true
	}
	if from.IsPrimitive || to.IsPrimitive {
		return false
	}
	if to.Name == "java/lang/Object" {
		return true
	}
	if from.IsArray {
		if to.IsArray {
			if from.Component.IsPrimitive || to.Component.IsPrimitive {
				return from.Component == to.Component
			}
			return r.IsAssignable(from.Component, to.Component)
		}
		if to.IsInterface() {
			return r.implementsInterface(from, to)
		}
		return false
	}
	for c := from.Super; c != nil; c = c.Super {
		if c == to {
			return true
		}
	}
	if to.IsInterface() {
		return r.implementsInterface(from, to)
	}
	return false
}

func (r *Runtime) implementsInterface(from, to *Class) bool {
	seen := map[*Class]bool{}
	var visit func(c *Class) bool
	visit = func(c *Class) bool {
		if c == nil || seen[c] {
			return false
		}
		seen[c] = true
		for _, in := range c.Interfaces {
			if in == to || visit(in) {
				return true
			}
		}
		if c.Super != nil {
			return visit(c.Super)
		}
		return false
	}
	return visit(from)
}

// IsInterface reports whether c has ACC_INTERFACE.
func (c *Class) IsInterface() bool { return c != nil && c.Flags&AccInterface != 0 }

// PackageName returns the runtime package of a class ("" for the default
// package).
func (c *Class) PackageName() string {
	if c == nil || c.IsArray || c.IsPrimitive {
		return ""
	}
	if i := strings.LastIndex(c.Name, "/"); i >= 0 {
		return c.Name[:i]
	}
	return ""
}

func (c *Class) addMethod(m *Method) error {
	key := m.Name + m.Descriptor
	if _, ok := c.MethodIndex[key]; ok {
		return fmt.Errorf("catty: duplicate method %s.%s%s", c.Name, m.Name, m.Descriptor)
	}
	c.MethodIndex[key] = m
	c.Methods = append(c.Methods, m)
	return nil
}

// NewStringUnits creates a Java string from UTF-16 units.
func (r *Runtime) NewStringUnits(units []uint16, intern bool) *Object {
	if intern {
		key := unitsKey(units)
		if o, ok := r.strings[key]; ok {
			return o
		}
		o := &Object{
			Class:    r.mustClass("java/lang/String"),
			Payload:  &stringData{Units: append([]uint16(nil), units...)},
			Identity: r.newIdentity(),
		}
		r.strings[key] = o
		return o
	}
	return &Object{
		Class:    r.mustClass("java/lang/String"),
		Payload:  &stringData{Units: append([]uint16(nil), units...)},
		Identity: r.newIdentity(),
	}
}

// NewStringGo creates a Java string from Go text.
func (r *Runtime) NewStringGo(s string) *Object {
	return r.NewStringUnits(utf16.Encode([]rune(s)), false)
}

type stringData struct {
	Units []uint16
}

func unitsKey(units []uint16) string {
	b := make([]byte, len(units)*2)
	for i, u := range units {
		b[i*2] = byte(u >> 8)
		b[i*2+1] = byte(u)
	}
	return string(b)
}

// StringText converts a Java String object to Go text.
func StringText(o *Object) (string, bool) {
	if o == nil {
		return "", false
	}
	d, ok := o.Payload.(*stringData)
	if !ok {
		return "", false
	}
	return string(utf16.Decode(d.Units)), true
}

// StringUnits returns a Java String's UTF-16 code units.
func StringUnits(o *Object) ([]uint16, bool) {
	if o == nil {
		return nil, false
	}
	d, ok := o.Payload.(*stringData)
	if !ok {
		return nil, false
	}
	return d.Units, true
}

// NewStringArray builds a Java String[] argument vector.
func (r *Runtime) NewStringArray(args []string) *Object {
	strClass := r.mustClass("java/lang/String")
	arr := r.NewArray(strClass, len(args))
	for i, a := range args {
		arr.Array.Elems[i] = RefValue(r.NewStringGo(a))
	}
	return arr
}

// boxData is the payload of every boxed primitive value.
type boxData struct {
	Kind byte // primitive descriptor letter
	Bits uint64
}

func boxedClassName(desc byte) string {
	switch desc {
	case 'Z':
		return "java/lang/Boolean"
	case 'B':
		return "java/lang/Byte"
	case 'C':
		return "java/lang/Character"
	case 'S':
		return "java/lang/Short"
	case 'I':
		return "java/lang/Integer"
	case 'J':
		return "java/lang/Long"
	case 'F':
		return "java/lang/Float"
	case 'D':
		return "java/lang/Double"
	}
	return ""
}

func primitiveDescOfBox(class string) byte {
	switch class {
	case "java/lang/Boolean":
		return 'Z'
	case "java/lang/Byte":
		return 'B'
	case "java/lang/Character":
		return 'C'
	case "java/lang/Short":
		return 'S'
	case "java/lang/Integer":
		return 'I'
	case "java/lang/Long":
		return 'J'
	case "java/lang/Float":
		return 'F'
	case "java/lang/Double":
		return 'D'
	}
	return 0
}

// BoxValue boxes v as the boxed class of the given primitive descriptor.
func (r *Runtime) BoxValue(desc byte, v Value) *Object {
	switch desc {
	case 'Z':
		b := AsInt(v) != 0
		if o, ok := r.boxedBools[b]; ok {
			return o
		}
		o := r.boxedObject(desc, uint64(boolBits(b)))
		r.boxedBools[b] = o
		return o
	case 'B':
		b := int8(AsInt(v))
		if o, ok := r.boxedBytes[b]; ok {
			return o
		}
		o := r.boxedObject(desc, uint64(uint8(b)))
		r.boxedBytes[b] = o
		return o
	case 'C':
		c := uint16(AsInt(v))
		if o, ok := r.boxedChars[c]; ok {
			return o
		}
		o := r.boxedObject(desc, uint64(c))
		r.boxedChars[c] = o
		return o
	case 'S':
		s := int16(AsInt(v))
		if o, ok := r.boxedShorts[s]; ok {
			return o
		}
		o := r.boxedObject(desc, uint64(uint16(s)))
		r.boxedShorts[s] = o
		return o
	case 'I':
		return r.BoxedInteger(AsInt(v))
	case 'J':
		l := AsLong(v)
		if o, ok := r.boxedLongs[l]; ok {
			return o
		}
		o := r.boxedObject(desc, uint64(l))
		r.boxedLongs[l] = o
		return o
	case 'F':
		return r.boxedObject(desc, uint64(AsFloatBits(v)))
	case 'D':
		return r.boxedObject(desc, AsDoubleBits(v))
	}
	panic("catty: cannot box descriptor " + string(desc))
}

func boolBits(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func (r *Runtime) boxedObject(desc byte, bits uint64) *Object {
	name := boxedClassName(desc)
	o := r.NewInstance(r.mustClass(name))
	o.Payload = &boxData{Kind: desc, Bits: bits}
	return o
}

// BoxValueOf boxes v into the class identified by a boxed class name.
func (r *Runtime) BoxValueOf(class string, v Value) *Object {
	desc := primitiveDescOfBox(class)
	return r.BoxValue(desc, v)
}

// UnboxValue reads a boxed primitive value.
func UnboxValue(o *Object) (boxData, bool) {
	if o == nil {
		return boxData{}, false
	}
	d, ok := o.Payload.(*boxData)
	if !ok {
		return boxData{}, false
	}
	return *d, true
}

// BoxedValue converts a box data payload into the engine Value.
func (d boxData) Value() Value {
	switch d.Kind {
	case 'Z', 'B', 'C', 'S', 'I':
		return IntValue(int32(uint32(d.Bits)))
	case 'J':
		return LongValue(int64(d.Bits))
	case 'F':
		return FloatBitsValue(uint32(d.Bits))
	case 'D':
		return DoubleBitsValue(d.Bits)
	}
	return Value{}
}

// BoxedInteger returns a java.lang.Integer for v. Java guarantees identity for
// the standard cache range; values outside it receive a fresh box.
func (r *Runtime) BoxedInteger(v int32) *Object {
	if v >= -128 && v <= 127 {
		if o, ok := r.boxedInts[v]; ok {
			return o
		}
	}
	o := r.boxedObject('I', uint64(uint32(v)))
	if v >= -128 && v <= 127 {
		r.boxedInts[v] = o
	}
	return o
}

// IntegerValue reads a boxed int.
func IntegerValue(o *Object) (int32, bool) {
	d, ok := UnboxValue(o)
	if !ok {
		return 0, false
	}
	switch d.Kind {
	case 'Z', 'B', 'C', 'S', 'I':
		return int32(uint32(d.Bits)), true
	}
	return 0, false
}

// ThrowableData carries the cause and Java frames captured at construction.
type ThrowableData struct {
	Message string
	Cause   *Object
	Frames  []FrameInfo
}

// FrameInfo is one captured Java frame.
type FrameInfo struct {
	Method *Method
	PC     int32
}

// NewThrowable allocates a throwable of the named bootstrap class and captures
// the current Java stack. A nil thread captures no frames.
func (r *Runtime) NewThrowable(t *Thread, className string) *Object {
	c, err := r.classForName(className)
	if err != nil {
		panic("catty: missing throwable class " + className)
	}
	o := r.NewInstance(c)
	td := &ThrowableData{}
	if t != nil {
		td.Frames = t.FramesSnapshot()
	}
	o.Payload = td
	return o
}

// ThrowableOf returns throwable payload data, creating it if necessary.
func ThrowableOf(o *Object) *ThrowableData {
	if o == nil {
		return nil
	}
	if d, ok := o.Payload.(*ThrowableData); ok {
		return d
	}
	d := &ThrowableData{}
	o.Payload = d
	return d
}

// CauseOf returns a throwable's cause (nil when absent).
func CauseOf(o *Object) *Object {
	if d := ThrowableOf(o); d != nil {
		return d.Cause
	}
	return nil
}

// SetCause records a throwable cause.
func SetCause(o, cause *Object) {
	ThrowableOf(o).Cause = cause
}

func (t *Thread) FramesSnapshot() []FrameInfo {
	out := make([]FrameInfo, len(t.frames))
	copy(out, t.frames)
	return out
}

// ArrayValue reads one array element as an engine Value.
func ArrayValue(a *ArrayData, index int) Value {
	switch a.Kind {
	case 'L':
		return a.Elems[index]
	case 'Z', 'B':
		return IntValue(int32(a.Bytes[index]))
	case 'C':
		return IntValue(int32(a.Chars[index]))
	case 'S':
		return IntValue(int32(a.Shorts[index]))
	case 'I':
		return IntValue(a.Ints[index])
	case 'J':
		return LongValue(a.Longs[index])
	case 'F':
		return FloatValue(a.Floats[index])
	case 'D':
		return DoubleValue(a.Doubles[index])
	}
	panic("catty: unknown array storage kind")
}

// SetArrayValue writes one array element after masking it to the element
// width; callers have already performed reference assignability checks.
func SetArrayValue(a *ArrayData, index int, v Value) {
	switch a.Kind {
	case 'L':
		a.Elems[index] = v
	case 'Z':
		a.Bytes[index] = int8(AsInt(v) & 1)
	case 'B':
		a.Bytes[index] = int8(AsInt(v))
	case 'C':
		a.Chars[index] = uint16(AsInt(v))
	case 'S':
		a.Shorts[index] = int16(AsInt(v))
	case 'I':
		a.Ints[index] = AsInt(v)
	case 'J':
		a.Longs[index] = AsLong(v)
	case 'F':
		a.Floats[index] = AsFloat(v)
	case 'D':
		a.Doubles[index] = AsDouble(v)
	default:
		panic("catty: unknown array storage kind")
	}
}
