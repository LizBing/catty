package rt

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"
)

// Runtime owns loader state, class metadata, static storage, mirrors, strings
// and the single logical Java thread used by the Gate 1 launcher.
type Runtime struct {
	Image     *Image
	Classes   map[string]*Class
	Bootstrap map[string]*Class
	Loader    *Loader

	intClass     *Class
	arrayClasses map[string]*Class
	strings      map[string]*Object
	boxedInts    map[int32]*Object
	classOrder   []*Class

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

// NewRuntime creates an uninstalled runtime for img.
func NewRuntime(img *Image) *Runtime {
	r := &Runtime{
		Image:        img,
		Classes:      map[string]*Class{},
		Bootstrap:    map[string]*Class{},
		arrayClasses: map[string]*Class{},
		strings:      map[string]*Object{},
		boxedInts:    map[int32]*Object{},
	}
	r.Loader = &Loader{Name: "app", Runtime: r}
	return r
}

// Installed reports whether Install completed.
func (r *Runtime) Installed() bool { return r.installed }

// Install validates the image ABI/profile and materializes immutable class
// definitions, static storage, mirrors and method metadata. No Java
// initializer runs here.
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

	for _, bc := range bootstrapInventory {
		c := &Class{
			Name:        bc.name,
			Flags:       bc.flags,
			MethodIndex: map[string]*Method{},
			State:       classInitialized,
			Loader:      r.Loader,
		}
		r.Bootstrap[bc.name] = c
		r.classOrder = append(r.classOrder, c)
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
		if _, ok := r.Classes[def.Name]; ok {
			return fmt.Errorf("catty: duplicate class definition %q in image", def.Name)
		}
		seen[def.Name] = true
		c := &Class{
			Name:        def.Name,
			Flags:       def.Flags,
			MethodIndex: map[string]*Method{},
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
					Static: true, Index: len(c.StaticField), Type: ft,
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
				if len(m.Insns) != 0 && m.CodeSize != 0 {
					// Generated code must execute its own lowered control flow;
					// retaining an interpreter body for AOT methods is a build
					// error in this profile.
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

	// The primitive int class is needed while resolving parameter classes.
	r.intClass = &Class{
		Name:        "int",
		Descriptor:  "I",
		IsPrimitive: true,
		Flags:       AccPublic | AccFinal,
		State:       classInitialized,
		Loader:      r.Loader,
	}

	// Resolve parameter classes now that all class shells exist.
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
	r.intClass.Mirror = &Object{Class: classClass, Payload: r.intClass, Identity: r.newIdentity()}

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
	for _, c := range r.allClasses() {
		for i, f := range c.StaticField {
			if c.Statics[i].Kind == 0 && !c.IsPrimitive {
				c.Statics[i] = zeroValueFor(f.Descriptor)
			}
		}
	}

	r.installed = true
	return nil
}

func zeroValueFor(desc string) Value {
	if desc == "I" {
		return IntValue(0)
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

// IntClass returns the canonical primitive int mirror class.
func (r *Runtime) IntClass() *Class { return r.intClass }

// classForType resolves a parsed type to its runtime class.
func (r *Runtime) classForType(t Type) (*Class, error) {
	switch t.Kind {
	case 'V':
		return nil, nil
	case 'I':
		if r.intClass == nil {
			return nil, fmt.Errorf("primitive int class is not installed")
		}
		return r.intClass, nil
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
	return nil, fmt.Errorf("unknown class %s", name)
}

func (r *Runtime) arrayClassFromDescriptor(desc string) (*Class, error) {
	if c, ok := r.arrayClasses[desc]; ok {
		return c, nil
	}
	if !strings.HasPrefix(desc, "[") {
		return nil, fmt.Errorf("not an array descriptor: %s", desc)
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
	case 'I':
		return nil, fmt.Errorf("primitive array %s is outside the Gate 1 profile", desc)
	default:
		return nil, fmt.Errorf("unsupported array descriptor: %s", desc)
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
		o.Fields[i] = Value{}
	}
	if r.AllocHook != nil {
		r.AllocHook(o)
	}
	return o
}

// NewArray allocates a reference array.
func (r *Runtime) NewArray(component *Class, length int) *Object {
	ac, err := r.arrayClassFor(component)
	if err != nil {
		panic(err)
	}
	o := &Object{
		Class:    ac,
		Array:    &ArrayData{Component: component, Elems: make([]Value, length)},
		Identity: r.newIdentity(),
	}
	return o
}

func (r *Runtime) arrayClassFor(component *Class) (*Class, error) {
	desc := ""
	if component.IsArray {
		desc = "[" + component.Descriptor
	} else {
		desc = "[L" + component.Name + ";"
	}
	return r.arrayClassFromDescriptor(desc)
}

// IsAssignable implements the reference/array assignability subset.
func (r *Runtime) IsAssignable(from, to *Class) bool {
	if from == nil || to == nil {
		return from == nil || from == to
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

// BoxedInteger returns a java.lang.Integer for v. Java guarantees identity for
// the standard cache range; values outside it receive a fresh box.
func (r *Runtime) BoxedInteger(v int32) *Object {
	if v >= -128 && v <= 127 {
		if o, ok := r.boxedInts[v]; ok {
			return o
		}
	}
	o := r.NewInstance(r.mustClass("java/lang/Integer"))
	o.Payload = &integerData{Value: v}
	if v >= -128 && v <= 127 {
		r.boxedInts[v] = o
	}
	return o
}

type integerData struct {
	Value int32
}

// IntegerValue reads a boxed int.
func IntegerValue(o *Object) (int32, bool) {
	if o == nil {
		return 0, false
	}
	d, ok := o.Payload.(*integerData)
	if !ok {
		return 0, false
	}
	return d.Value, true
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
