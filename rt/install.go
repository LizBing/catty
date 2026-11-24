package rt

import (
	"fmt"
	"sort"
	"strings"
)

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
	profile, err := ProfileForImage(r.Image.Profile, r.Image.ABI)
	if err != nil {
		return err
	}
	if r.Image.Plugins != profile.Plugins {
		return fmt.Errorf("catty: image profile %q plugin flag %v is inconsistent", r.Image.Profile, r.Image.Plugins)
	}
	if r.Image.Threads != profile.Threads {
		return fmt.Errorf("catty: image profile %q threading flag %v is inconsistent", r.Image.Profile, r.Image.Threads)
	}
	if r.Image.Lifecycle != profile.Lifecycle {
		return fmt.Errorf("catty: image profile %q lifecycle flag %v is inconsistent", r.Image.Profile, r.Image.Lifecycle)
	}
	if profile.Plugins && r.Image.Limits == nil {
		r.Image.Limits = profile.Limits
	}
	if len(r.Image.Opcodes) > 0 {
		for _, name := range r.Image.Opcodes {
			if !admittedOpcodeName(profile, name) {
				return fmt.Errorf("catty: image opcode inventory contains unsupported opcode %q", name)
			}
		}
	}
	if len(r.Image.Capabilities) > 0 {
		canonical := map[string]bool{}
		for _, name := range ProfileCapabilities(profile.ID) {
			canonical[name] = true
		}
		for _, name := range r.Image.Capabilities {
			if !canonical[name] {
				return fmt.Errorf("catty: image capability %q is not published by profile %q", name, profile.ID)
			}
		}
	}
	if err := r.validateRetainedThreading(profile); err != nil {
		return err
	}
	r.profile = profile
	r.bootstrapDefs = profile.bootstrapByName()

	if err := r.installBootstrap(profile); err != nil {
		return err
	}
	// Cache the Thread class identity so lifecycle hook validation can run
	// under the lifecycle guard without acquiring the metadata mutex.
	if profile.Threads {
		r.threadClass = r.Bootstrap["java/lang/Thread"]
	}
	if err := r.installImageClasses(); err != nil {
		return err
	}
	if err := r.linkAll(profile); err != nil {
		return err
	}
	if err := r.checkHierarchies(); err != nil {
		return err
	}
	if err := r.layoutAll(profile); err != nil {
		return err
	}
	if err := r.installAllMethods(profile); err != nil {
		return err
	}
	if err := r.resolveAllParameterClasses(); err != nil {
		return err
	}
	if err := r.buildMirrors(); err != nil {
		return err
	}
	if profile.Lifecycle {
		// The Runtime singleton is created before any Java code can ask for it.
		// It has no public constructor in the bootstrap inventory.
		r.runtimeObject = r.NewInstance(r.Bootstrap["java/lang/Runtime"])
	}
	if err := r.initBootstrapStatics(profile); err != nil {
		return err
	}
	r.initializeImageStatics()
	r.installed = true
	return nil
}

// admittedOpcodeName reports whether the profile publishes name as an admitted
// opcode.
func admittedOpcodeName(profile *Profile, name string) bool {
	for op, n := range opcodeNames {
		if n == name {
			return AdmitsOpcodeFor(profile.ID, op)
		}
	}
	return false
}

// validateRetainedThreading checks the actual retained image content against
// the declared profile. A supplied image may omit its opcode inventory, so the
// monitor opcodes and synchronized method flags are also checked directly;
// earlier profiles must keep rejecting both even for hand-authored images.
func (r *Runtime) validateRetainedThreading(profile *Profile) error {
	for _, def := range r.Image.Classes {
		if def == nil {
			continue
		}
		for _, md := range def.Methods {
			if md == nil {
				continue
			}
			if md.Flags&AccSynchronized != 0 {
				if !profile.Threads {
					return fmt.Errorf("catty: image method %s.%s%s is synchronized, which profile %q does not admit",
						def.Name, md.Name, md.Descriptor, profile.ID)
				}
				switch {
				case md.Flags&AccNative != 0:
					return fmt.Errorf("catty: image method %s.%s%s is synchronized native, which is unsupported",
						def.Name, md.Name, md.Descriptor)
				case md.Name == "<init>":
					return fmt.Errorf("catty: image instance constructor %s.%s%s must not be synchronized",
						def.Name, md.Name, md.Descriptor)
				case md.Flags&AccAbstract != 0:
					return fmt.Errorf("catty: image abstract method %s.%s%s must not be synchronized",
						def.Name, md.Name, md.Descriptor)
				case def.Flags&AccInterface != 0 && md.Name != "<clinit>":
					return fmt.Errorf("catty: image interface method %s.%s%s must not be synchronized",
						def.Name, md.Name, md.Descriptor)
				}
			}
			if profile.Threads && md.Name == "<clinit>" {
				if md.Flags&AccStatic == 0 || md.Descriptor != "()V" {
					return fmt.Errorf("catty: image class initializer %s.<clinit>%s must be static ()V",
						def.Name, md.Descriptor)
				}
				if md.Flags&AccAbstract != 0 {
					return fmt.Errorf("catty: image class initializer %s.<clinit> must not be abstract", def.Name)
				}
			}
			for _, insn := range md.Insns {
				if insn.Op != OpMonitorenter && insn.Op != OpMonitorexit {
					continue
				}
				if !profile.Threads {
					return fmt.Errorf("catty: image method %s.%s%s retains %s, which profile %q does not admit",
						def.Name, md.Name, md.Descriptor, OpcodeName(insn.Op), profile.ID)
				}
			}
		}
	}
	return nil
}

func (r *Runtime) installBootstrap(profile *Profile) error {
	for _, bc := range profile.bootstrapInventory() {
		c := &Class{
			Name:        bc.name,
			Flags:       bc.flags | AccPublic,
			MethodIndex: map[string]*Method{},
			State:       classInitialized,
			Loader:      r.BootstrapLoader,
		}
		r.Bootstrap[bc.name] = c
		r.BootstrapLoader.Classes[bc.name] = c
		r.classOrder = append(r.classOrder, c)
		r.metrics.classLoaded()
	}
	for _, pc := range PrimitiveClasses {
		c := &Class{
			Name:        pc.Name,
			Descriptor:  pc.Descriptor,
			IsPrimitive: true,
			Flags:       AccPublic | AccFinal,
			MethodIndex: map[string]*Method{},
			State:       classInitialized,
			Loader:      r.BootstrapLoader,
		}
		r.primitives[pc.Descriptor[0]] = c
		r.classOrder = append(r.classOrder, c)
		r.primitiveOrder = append(r.primitiveOrder, c)
		r.metrics.classLoaded()
	}
	return nil
}

func reservedImageName(profile *Profile, name string) bool {
	if strings.HasPrefix(name, "java/") || strings.HasPrefix(name, "javax/") ||
		strings.HasPrefix(name, "jdk/") || strings.HasPrefix(name, "sun/") ||
		strings.HasPrefix(name, "com/sun/") || strings.HasPrefix(name, "org/w3c/") ||
		strings.HasPrefix(name, "org/xml/") {
		return true
	}
	if profile != nil && profile.Plugins && (name == "catty/runtime" || strings.HasPrefix(name, "catty/runtime/")) {
		return true
	}
	return false
}

func (r *Runtime) installImageClasses() error {
	seen := map[string]bool{}
	for _, def := range r.Image.Classes {
		if def == nil || def.Name == "" {
			return fmt.Errorf("catty: image contains an unnamed class definition")
		}
		if strings.HasPrefix(def.Name, "[") {
			return fmt.Errorf("catty: image class %q has an array name", def.Name)
		}
		if _, ok := r.Bootstrap[def.Name]; ok || reservedImageName(r.profile, def.Name) {
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
		r.Loader.Classes[def.Name] = c
		r.metrics.classLoaded()
	}
	return nil
}

func (r *Runtime) linkAll(profile *Profile) error {
	for _, bc := range profile.bootstrapInventory() {
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
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		if def.Super != "" {
			s, ok := r.Loader.Lookup(def.Super)
			if !ok {
				return fmt.Errorf("catty: class %s has unknown superclass %s", def.Name, def.Super)
			}
			c.Super = s
		}
		for _, in := range def.Interfaces {
			i, ok := r.Loader.Lookup(in)
			if !ok {
				return fmt.Errorf("catty: class %s has unknown interface %s", def.Name, in)
			}
			c.Interfaces = append(c.Interfaces, i)
		}
	}
	return nil
}

func (r *Runtime) layoutAll(profile *Profile) error {
	state := map[*Class]uint8{}
	var layout func(c *Class) error
	layout = func(c *Class) error {
		switch state[c] {
		case 1:
			return fmt.Errorf("catty: class hierarchy cycle at %s", c.Name)
		case 2:
			return nil
		}
		state[c] = 1
		if c.Super != nil {
			if err := layout(c.Super); err != nil {
				return err
			}
		}
		if c.IsPrimitive {
			c.Statics = nil
			c.InstanceSize = 0
			state[c] = 2
			return nil
		}
		if bc, ok := r.bootstrapDefs[c.Name]; ok {
			if err := r.layoutBootstrapClass(c, bc); err != nil {
				return err
			}
			state[c] = 2
			return nil
		}
		var def *ClassDef
		for _, d := range r.Image.Classes {
			if d.Name == c.Name {
				def = d
				break
			}
		}
		if def == nil {
			return fmt.Errorf("catty: missing image definition for %s", c.Name)
		}
		if err := r.layoutDefinition(c, def); err != nil {
			return err
		}
		state[c] = 2
		return nil
	}
	for _, bc := range profile.bootstrapInventory() {
		if err := layout(r.Bootstrap[bc.name]); err != nil {
			return err
		}
	}
	for _, c := range r.primitiveOrder {
		if err := layout(c); err != nil {
			return err
		}
	}
	for _, def := range r.Image.Classes {
		if err := layout(r.Classes[def.Name]); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) layoutBootstrapClass(c *Class, bc *bsClass) error {
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
	c.InstanceSize = size
	c.Statics = make([]Value, len(c.StaticField))
	return nil
}

// layoutDefinition assigns declared field offsets and static slots for one
// class. The caller has already laid out the superclass.
func (r *Runtime) layoutDefinition(c *Class, def *ClassDef) error {
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
	c.StaticField = nil
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
	return nil
}

func (r *Runtime) installAllMethods(profile *Profile) error {
	for _, bc := range profile.bootstrapInventory() {
		c := r.Bootstrap[bc.name]
		for _, bm := range bc.methods {
			mt, err := ParseMethodType(bm.descriptor)
			if err != nil {
				return fmt.Errorf("catty: bootstrap method %s.%s: %w", c.Name, bm.name, err)
			}
			engine, reason := "intrinsic", "bootstrap intrinsic"
			if bm.flags&AccAbstract != 0 {
				engine, reason = "abstract", "bootstrap abstract declaration"
			}
			m := &Method{
				Owner: c, Name: bm.name, Descriptor: bm.descriptor, Flags: bm.flags,
				Engine: engine, Reason: reason, Intrinsic: bm.intrinsic, Type: mt,
			}
			if err := c.addMethod(m); err != nil {
				return err
			}
			r.metrics.methodPublished(engine)
		}
	}
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		for _, md := range def.Methods {
			if err := r.installMethodDef(c, md); err != nil {
				return err
			}
			r.metrics.methodPublished(md.Engine)
		}
	}
	return nil
}

func (r *Runtime) installMethodDef(c *Class, md *MethodDef) error {
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
	return c.addMethod(m)
}

func (r *Runtime) resolveAllParameterClasses() error {
	for _, c := range r.allClasses() {
		if err := r.resolveClassMetadata(c, true); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) resolveClassMetadata(c *Class, strict bool) error {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	return r.resolveClassMetadataLocked(c, strict)
}

func (r *Runtime) resolveClassMetadataLocked(c *Class, strict bool) error {
	for _, m := range c.Methods {
		m.ParamClass = make([]*Class, len(m.Type.Params))
		for i, p := range m.Type.Params {
			pc, err := r.classForTypeInLoaderLocked(c.Loader, p)
			if err != nil {
				if strict {
					return fmt.Errorf("catty: method %s.%s%s parameter %d: %w", c.Name, m.Name, m.Descriptor, i, err)
				}
				continue
			}
			m.ParamClass[i] = pc
		}
	}
	if strict {
		for _, f := range c.StaticField {
			if _, err := r.classForTypeInLoaderLocked(c.Loader, f.Type); err != nil {
				return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
			}
		}
		for _, f := range c.InstanceField {
			if _, err := r.classForTypeInLoaderLocked(c.Loader, f.Type); err != nil {
				return fmt.Errorf("catty: field %s.%s: %w", c.Name, f.Name, err)
			}
		}
	}
	return nil
}

func (r *Runtime) buildMirrors() error {
	classClass := r.mustClass("java/lang/Class")
	makeMirror := func(c *Class) {
		if c.Mirror == nil {
			c.Mirror = &Object{Class: classClass, Payload: c, Identity: r.newIdentity()}
		}
	}
	for _, c := range r.classOrder {
		makeMirror(c)
	}
	for _, def := range r.Image.Classes {
		makeMirror(r.Classes[def.Name])
	}
	return nil
}

func (r *Runtime) initBootstrapStatics(profile *Profile) error {
	for _, bc := range profile.bootstrapInventory() {
		c := r.Bootstrap[bc.name]
		for i, f := range bc.fields {
			if f.initial != nil {
				c.Statics[i] = f.initial(r)
			} else {
				c.Statics[i] = zeroValueFor(f.descriptor)
			}
		}
	}
	return nil
}

func (r *Runtime) initializeImageStatics() {
	for _, def := range r.Image.Classes {
		c := r.Classes[def.Name]
		for i, f := range c.StaticField {
			if i < len(c.Statics) {
				c.Statics[i] = zeroValueFor(f.Descriptor)
			}
		}
	}
}

// sortedClassNames is a deterministic helper shared with plugin installation.
func sortedClassNames(classes map[string]*Class) []string {
	names := make([]string, 0, len(classes))
	for n := range classes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
