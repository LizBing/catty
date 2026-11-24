package rt

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf16"
)

// runtimeClock is the monotonic clock used for every Java timer wait. Tests
// inject a manual clock so timeout transitions are deterministic.
type runtimeClock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realRuntimeClock struct{}

func (realRuntimeClock) Now() time.Time                         { return time.Now() }
func (realRuntimeClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Runtime owns loader state, class metadata, static storage, mirrors, strings
// and the logical Java threads used by the launcher.
//
// Three runtime mutexes partition the shared runtime state. memMu orders every
// access to published Java instance/static slots and array elements. metaMu
// guards resolution and dispatch caches, loader publication and counters,
// loading constraints, lazy array/mirror interning and string/boxing caches.
// thrMu guards logical thread state, monitor ownership and wait sets. The three
// are never nested with one another and never held across Java execution,
// monitor/class-initialization waits or host I/O.
type Runtime struct {
	memMu  sync.Mutex
	metaMu sync.Mutex
	thrMu  sync.Mutex

	Image *Image
	// Classes is the application loader's class table. Bootstrap is the
	// bootstrap loader's class table. Both are retained for host-facing code
	// and tests; the loaders themselves are authoritative.
	Classes         map[string]*Class
	Bootstrap       map[string]*Class
	Loader          *Loader
	BootstrapLoader *Loader

	profile       *Profile
	bootstrapDefs map[string]*bsClass

	primitives          map[byte]*Class
	arrayClasses        map[*Class]*Class
	primitiveArrayClass map[byte]*Class
	strings             map[string]*Object
	boxedInts           map[int32]*Object
	boxedLongs          map[int64]*Object
	boxedShorts         map[int16]*Object
	boxedBytes          map[int8]*Object
	boxedChars          map[uint16]*Object
	boxedBools          map[bool]*Object
	classOrder          []*Class
	primitiveOrder      []*Class
	pluginLoaders       []*Loader
	pluginClassBytes    int64
	classLoaderMirror   *Object
	pluginMirror        *Object
	nextLoaderID        int

	resolutions map[resolutionKey]*Resolution
	dispatch    map[dispatchKey]*dispatchEntry
	initiations map[initiationKey]*Class
	constraints map[constraintKey]struct{}

	admit func(*AdmissionRequest) (*AdmissionResult, error)

	nextIdentity atomic.Uint64
	nextThread   int64
	installed    bool
	trace        *trace
	mainThread   *Thread

	liveNonDaemon int
	liveCond      *sync.Cond
	clock         runtimeClock

	// ClassInitWaitHook is a deterministic-test hook invoked immediately
	// before a competing thread parks on an in-progress class initialization.
	ClassInitWaitHook func(c *Class, waiter *Thread)
	threadAudit       []*Thread
	unnamedPeers      int64
	exitHook          func()

	// diagMu serializes uncaught-thread diagnostics.
	diagMu sync.Mutex

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
	Class  *Class
	Thrown *Object
}

// NewRuntime creates an uninstalled runtime for img.
func NewRuntime(img *Image) *Runtime {
	r := &Runtime{
		Image:               img,
		Classes:             map[string]*Class{},
		Bootstrap:           map[string]*Class{},
		primitives:          map[byte]*Class{},
		arrayClasses:        map[*Class]*Class{},
		primitiveArrayClass: map[byte]*Class{},
		strings:             map[string]*Object{},
		boxedInts:           map[int32]*Object{},
		boxedLongs:          map[int64]*Object{},
		boxedShorts:         map[int16]*Object{},
		boxedBytes:          map[int8]*Object{},
		boxedChars:          map[uint16]*Object{},
		boxedBools:          map[bool]*Object{},
		resolutions:         map[resolutionKey]*Resolution{},
		dispatch:            map[dispatchKey]*dispatchEntry{},
		initiations:         map[initiationKey]*Class{},
		constraints:         map[constraintKey]struct{}{},
		admit:               defaultAdmission,
	}
	r.liveCond = sync.NewCond(&r.thrMu)
	r.clock = realRuntimeClock{}
	r.BootstrapLoader = &Loader{ID: 0, Name: "bootstrap", Runtime: r, Classes: r.Bootstrap}
	r.Loader = &Loader{ID: 1, Name: "app", Runtime: r, Parent: r.BootstrapLoader, Classes: r.Classes}
	r.nextLoaderID = 1
	return r
}

// Installed reports whether Install completed.
func (r *Runtime) Installed() bool { return r.installed }

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
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.lookupClassShellLocked(name)
}

func (r *Runtime) lookupClassShellLocked(name string) (*Class, bool) {
	if strings.HasPrefix(name, "[") {
		if c, err := r.classForLoaderNameLocked(r.Loader, name); err == nil {
			return c, true
		}
		return nil, false
	}
	if desc, ok := PrimitiveDescriptor(name); ok {
		if c, ok := r.primitives[desc]; ok {
			return c, true
		}
	}
	return r.Loader.Lookup(name)
}

// arrayClassFromDescriptor resolves an array descriptor through the
// application loader. Symbolic uses route through classForLoaderName instead.
func (r *Runtime) arrayClassFromDescriptor(desc string) (*Class, error) {
	return r.classForLoaderName(r.Loader, desc)
}

func (r *Runtime) allClasses() []*Class {
	out := make([]*Class, 0, len(r.classOrder)+len(r.Classes))
	out = append(out, r.classOrder...)
	for _, n := range sortedClassNames(r.Classes) {
		out = append(out, r.Classes[n])
	}
	for _, l := range r.pluginLoaders {
		for _, n := range sortedClassNames(l.Classes) {
			out = append(out, l.Classes[n])
		}
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
	return r.nextIdentity.Add(1)
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

// mustClassLocked is mustClass for a caller that already holds the metadata
// mutex.
func (r *Runtime) mustClassLocked(name string) *Class {
	c, ok := r.lookupClassShellLocked(name)
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

// classForType resolves a parsed type through the application loader.
func (r *Runtime) classForType(t Type) (*Class, error) {
	return r.classForTypeInLoader(r.Loader, t)
}

// classForTypeInLoader resolves a parsed type through loader.
func (r *Runtime) classForTypeInLoader(loader *Loader, t Type) (*Class, error) {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.classForTypeInLoaderLocked(loader, t)
}

func (r *Runtime) classForTypeInLoaderLocked(loader *Loader, t Type) (*Class, error) {
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
		return r.classForLoaderNameLocked(loader, t.DescriptorOrName())
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
// classForName resolves a name through the application loader. Host-facing
// callers use this convenience binding; symbolic resolution from Java code uses
// the referring class's loader.
func (r *Runtime) classForName(name string) (*Class, error) {
	return r.classForLoaderName(r.Loader, name)
}

// classForLoaderName resolves an internal class name, primitive descriptor
// pseudo-name or array descriptor through the supplied loader. Arrays are
// interned by immediate component identity, so equal descriptors from two
// sibling loaders produce distinct array classes.
func (r *Runtime) classForLoaderName(loader *Loader, name string) (*Class, error) {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.classForLoaderNameLocked(loader, name)
}

// classForLoaderNameLocked is classForLoaderName for a caller holding metaMu.
func (r *Runtime) classForLoaderNameLocked(loader *Loader, name string) (*Class, error) {
	if loader == nil {
		loader = r.Loader
	}
	if strings.HasPrefix(name, "[") {
		componentDesc := name[1:]
		var component *Class
		var err error
		switch componentDesc[0] {
		case '[':
			component, err = r.classForLoaderNameLocked(loader, componentDesc)
		case 'L':
			if !strings.HasSuffix(componentDesc, ";") {
				return nil, fmt.Errorf("malformed array descriptor: %s", name)
			}
			component, err = r.classForLoaderNameLocked(loader, componentDesc[1:len(componentDesc)-1])
		default:
			c, ok := r.primitives[componentDesc[0]]
			if !ok {
				return nil, fmt.Errorf("unsupported array descriptor: %s", name)
			}
			component = c
		}
		if err != nil {
			return nil, err
		}
		return r.arrayClassForComponentLocked(component), nil
	}
	if desc, ok := PrimitiveDescriptor(name); ok {
		if c, ok := r.primitives[desc]; ok {
			return c, nil
		}
	}
	if c, ok := loader.Lookup(name); ok {
		return c, nil
	}
	return nil, fmt.Errorf("unknown class %s", name)
}

// arrayClassForComponent interns a reference array by immediate component
// identity and a primitive array by its canonical primitive class.
func (r *Runtime) arrayClassForComponent(component *Class) *Class {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.arrayClassForComponentLocked(component)
}

func (r *Runtime) arrayClassForComponentLocked(component *Class) *Class {
	if component == nil {
		return nil
	}
	if component.IsPrimitive {
		desc := component.Descriptor[0]
		if c, ok := r.primitiveArrayClass[desc]; ok {
			return c
		}
		c := r.newArrayClassLocked("["+component.Descriptor, component, r.BootstrapLoader)
		r.primitiveArrayClass[desc] = c
		return c
	}
	if c, ok := r.arrayClasses[component]; ok {
		return c
	}
	loader := component.Loader
	if loader == nil {
		loader = r.BootstrapLoader
	}
	c := r.newArrayClassLocked(arrayDescriptorOf(component), component, loader)
	r.arrayClasses[component] = c
	return c
}

func (r *Runtime) newArrayClassLocked(desc string, component *Class, loader *Loader) *Class {
	c := &Class{
		Name:       desc,
		Descriptor: desc,
		IsArray:    true,
		Component:  component,
		Super:      r.mustClassLocked("java/lang/Object"),
		Interfaces: []*Class{r.mustClassLocked("java/lang/Cloneable"), r.mustClassLocked("java/io/Serializable")},
		Flags:      AccPublic | AccFinal,
		State:      classInitialized,
		Loader:     loader,
	}
	classClass := r.mustClassLocked("java/lang/Class")
	c.Mirror = &Object{Class: classClass, Payload: c, Identity: r.newIdentity()}
	return c
}

// NewInstance allocates an instance of c with zero-initialized slots.
func (r *Runtime) NewInstance(c *Class) *Object {
	o := &Object{
		Class:    c,
		Fields:   make([]Value, c.InstanceSize),
		Identity: r.newIdentity(),
	}
	// Throwable payloads are frozen before publication: every allocation of a
	// throwable has its data attached here, and only constructors running before
	// publication may fill in the cause, message or frames.
	if isThrowableClass(c) {
		o.Payload = &ThrowableData{}
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
	if component == nil {
		return nil, fmt.Errorf("catty: nil array component")
	}
	return r.arrayClassForComponent(component), nil
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

// NewStringUnits creates a Java string from UTF-16 units. Interned strings are
// canonical and shared, so the cache transition is a metadata operation.
func (r *Runtime) NewStringUnits(units []uint16, intern bool) *Object {
	if !intern {
		return &Object{
			Class:    r.mustClass("java/lang/String"),
			Payload:  &stringData{Units: append([]uint16(nil), units...)},
			Identity: r.newIdentity(),
		}
	}
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.newStringUnitsLocked(units)
}

func (r *Runtime) newStringUnitsLocked(units []uint16) *Object {
	key := unitsKey(units)
	if o, ok := r.strings[key]; ok {
		return o
	}
	o := &Object{
		Class:    r.mustClassLocked("java/lang/String"),
		Payload:  &stringData{Units: append([]uint16(nil), units...)},
		Identity: r.newIdentity(),
	}
	r.strings[key] = o
	return o
}

// NewStringGo creates a Java string from Go text.
func (r *Runtime) NewStringGo(s string) *Object {
	return r.NewStringUnits(utf16.Encode([]rune(s)), false)
}

// NewStringGoInterned creates a canonical Java string from Go text.
func (r *Runtime) NewStringGoInterned(s string) *Object {
	return r.NewStringUnits(utf16.Encode([]rune(s)), true)
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
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.boxValueLocked(desc, v)
}

func (r *Runtime) boxValueLocked(desc byte, v Value) *Object {
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
		return r.boxedIntegerLocked(AsInt(v))
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
	o := r.NewInstance(r.mustClassLocked(name))
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
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.boxedIntegerLocked(v)
}

func (r *Runtime) boxedIntegerLocked(v int32) *Object {
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
	Code    string
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
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.newThrowableLocked(t, className)
}

// newThrowableLocked is NewThrowable for a caller holding metaMu. It performs no
// Java execution and acquires no other runtime lock.
func (r *Runtime) newThrowableLocked(t *Thread, className string) *Object {
	c, err := r.classForLoaderNameLocked(r.Loader, className)
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

// ThrowableOf returns throwable payload data. The payload is attached at
// allocation and never created after publication, so this read is race-free.
func ThrowableOf(o *Object) *ThrowableData {
	if o == nil {
		return nil
	}
	d, _ := o.Payload.(*ThrowableData)
	return d
}

// isThrowableClass walks the immutable superclass chain without taking a
// runtime lock; NewInstance may be called with the metadata mutex held.
func isThrowableClass(c *Class) bool {
	for x := c; x != nil; x = x.Super {
		if x.Name == "java/lang/Throwable" {
			return true
		}
	}
	return false
}

// CauseOf returns a throwable's cause (nil when absent).
func CauseOf(o *Object) *Object {
	if d := ThrowableOf(o); d != nil {
		return d.Cause
	}
	return nil
}

// SetCause records a throwable cause. Callers set causes during construction,
// before the throwable is published.
func SetCause(o, cause *Object) {
	if d := ThrowableOf(o); d != nil {
		d.Cause = cause
	}
}

func (t *Thread) FramesSnapshot() []FrameInfo {
	out := make([]FrameInfo, len(t.frames))
	copy(out, t.frames)
	return out
}

// fieldRead reads one published instance slot under the runtime memory mutex.
func (r *Runtime) fieldRead(o *Object, offset int) Value {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	return o.Fields[offset]
}

// fieldWrite writes one published instance slot under the runtime memory mutex.
func (r *Runtime) fieldWrite(o *Object, offset int, v Value) {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	o.Fields[offset] = v
}

// staticRead reads one published static slot under the runtime memory mutex.
func (r *Runtime) staticRead(c *Class, index int) Value {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	return c.Statics[index]
}

// staticWrite writes one published static slot under the runtime memory mutex.
func (r *Runtime) staticWrite(c *Class, index int, v Value) {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	c.Statics[index] = v
}

// arrayRead reads one published array element under the runtime memory mutex.
func (r *Runtime) arrayRead(a *ArrayData, index int) Value {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	return arrayValueLocked(a, index)
}

// arrayWrite writes one published array element under the runtime memory
// mutex. Callers have already performed the reference assignability check.
func (r *Runtime) arrayWrite(a *ArrayData, index int, v Value) {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	setArrayValueLocked(a, index, v)
}

// arrayLength reads the fixed array length. Length, kind and component are
// immutable after allocation and need no memory mutex.
func arrayLength(a *ArrayData) int { return a.Len() }

// arraySnapshotElems copies a reference array's element slots under the memory
// mutex. Reflection uses it to detach its argument view from concurrent Java
// mutation.
func (r *Runtime) arraySnapshotElems(a *ArrayData) []Value {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	out := make([]Value, len(a.Elems))
	copy(out, a.Elems)
	return out
}

// ArrayValue reads one array element as an engine Value.
func ArrayValue(a *ArrayData, index int) Value {
	return arrayValueLocked(a, index)
}

func arrayValueLocked(a *ArrayData, index int) Value {
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
	setArrayValueLocked(a, index, v)
}

func setArrayValueLocked(a *ArrayData, index int, v Value) {
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
