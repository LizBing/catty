package build

import (
	"fmt"
	"runtime"
	"strings"

	"catty/internal/classfile"
	"catty/internal/verify"
	"catty/rt"
)

// Selection is one method's engine choice and evidence.
type Selection struct {
	Engine     string
	Reason     string
	EntryIndex int
	Symbol     string
	Heights    map[int32]int
	Shapes     map[int32][]int
}

type aotEntry struct {
	class   *classfile.Class
	method  *classfile.Method
	symbol  string
	heights map[int32]int
	shapes  map[int32][]int
}

// Build holds the staged build state.
type Build struct {
	opts       Options
	profile    *rt.Profile
	loaded     *Loaded
	ctx        *verify.Context
	selections map[*classfile.Method]*Selection
	aot        []*aotEntry
	image      *rt.Image
	report     *Report
	unresolved []UnresolvedRecord
}

func newBuild(opts Options) *Build {
	return &Build{opts: opts, selections: map[*classfile.Method]*Selection{}}
}

func (b *Build) prepare() error {
	if b.opts.Mode != "auto" && b.opts.Mode != "aot" && b.opts.Mode != "interp" {
		return fmt.Errorf("catty: unsupported mode %q (want auto, aot or interp)", b.opts.Mode)
	}
	if b.opts.Main == "" {
		return fmt.Errorf("catty: --main is required")
	}
	if b.opts.ReportPath == "" || b.opts.OutputPath == "" {
		return fmt.Errorf("catty: --report and --output are required")
	}
	switch {
	case b.opts.EnableLifecycle:
		// Gate 5 includes Gate 4 threading and Gate 3 plugin admission.
		b.profile = rt.Gate5Profile
	case b.opts.EnableThreads:
		// Gate 4 includes Gate 3 plugin admission.
		b.profile = rt.Gate4Profile
	case b.opts.EnablePlugins:
		b.profile = rt.Gate3Profile
	default:
		b.profile = rt.Gate2Profile
	}
	loaded, err := LoadClasspathFor(b.opts.Classpath, b.profile)
	if err != nil {
		return err
	}
	b.loaded = loaded
	b.ctx = &verify.Context{Classes: loaded.ByName, Profile: b.profile.ID}
	if err := b.validateClasses(); err != nil {
		return err
	}
	for _, c := range loaded.Classes {
		if err := verify.VerifyClass(c, b.ctx); err != nil {
			return fmt.Errorf("%s: %w", c.Origin, err)
		}
	}
	if err := b.validateMain(); err != nil {
		return err
	}
	var policy *Policy
	if b.opts.PolicyPath != "" {
		if b.opts.Mode != "auto" {
			return fmt.Errorf("catty: --policy is only valid with --mode auto")
		}
		p, err := LoadPolicy(b.opts.PolicyPath)
		if err != nil {
			return err
		}
		policy = p
	}
	if err := b.validatePolicy(policy); err != nil {
		return err
	}
	if err := b.selectEngines(policy); err != nil {
		return err
	}
	if err := b.resolveReferences(); err != nil {
		return err
	}
	if err := b.buildImage(); err != nil {
		return err
	}
	b.buildReport()
	return nil
}

func (b *Build) validateClasses() error {
	for _, c := range b.loaded.Classes {
		if c.Flags&0x8000 != 0 {
			return fmt.Errorf("%s: module classes are unsupported", c.Origin)
		}
		if c.Name == "java/lang/Object" {
			return fmt.Errorf("%s: application class %s is reserved", c.Origin, c.Name)
		}
		if c.Super == "" {
			return fmt.Errorf("%s: class %s has no superclass", c.Origin, c.Name)
		}
		if !b.classAvailable(c.Super) {
			return fmt.Errorf("%s: class %s has unknown superclass %s", c.Origin, c.Name, c.Super)
		}
		for _, in := range c.Interfaces {
			if !b.classAvailable(in) {
				return fmt.Errorf("%s: class %s has unknown interface %s", c.Origin, c.Name, in)
			}
		}
	}
	// Application hierarchy cycles.
	const (
		white = 0
		grey  = 1
		black = 2
	)
	state := map[string]uint8{}
	var visit func(name string) error
	visit = func(name string) error {
		switch state[name] {
		case grey:
			return fmt.Errorf("catty: class hierarchy cycle at %s", name)
		case black:
			return nil
		}
		state[name] = grey
		if c := b.loaded.ByName[name]; c != nil && c.Super != "" {
			if b.loaded.ByName[c.Super] != nil {
				if err := visit(c.Super); err != nil {
					return err
				}
			}
		}
		state[name] = black
		return nil
	}
	for _, c := range b.loaded.Classes {
		if err := visit(c.Name); err != nil {
			return err
		}
	}
	return nil
}

func (b *Build) classAvailable(name string) bool {
	if strings.HasPrefix(name, "[") {
		_, err := rt.ParseFieldType(name)
		return err == nil
	}
	if rt.BootstrapHasClassFor(b.profile.ID, name) {
		return true
	}
	_, ok := b.loaded.ByName[name]
	return ok
}

func (b *Build) validateMain() error {
	internal := strings.ReplaceAll(b.opts.Main, ".", "/")
	c, ok := b.loaded.ByName[internal]
	if !ok {
		return fmt.Errorf("catty: entry class %s was not found on the classpath", b.opts.Main)
	}
	if c.Flags&rt.AccInterface != 0 {
		return fmt.Errorf("catty: entry class %s is an interface and has no main method", b.opts.Main)
	}
	for _, m := range c.Methods {
		if m.Name != "main" || m.Descriptor != "([Ljava/lang/String;)V" {
			continue
		}
		if m.Flags&rt.AccPublic == 0 || m.Flags&rt.AccStatic == 0 {
			return fmt.Errorf("catty: entry class %s main([Ljava/lang/String;)V must be public static", b.opts.Main)
		}
		return nil
	}
	return fmt.Errorf("catty: entry class %s has no public static main([Ljava/lang/String;)V", b.opts.Main)
}

func (b *Build) findDeclaredMethod(key MethodKey) (*classfile.Method, bool) {
	c, ok := b.loaded.ByName[key.Owner]
	if !ok {
		return nil, false
	}
	for _, m := range c.Methods {
		if m.Name == key.Name && m.Descriptor == key.Descriptor {
			return m, true
		}
	}
	return nil, false
}

func (b *Build) validatePolicy(p *Policy) error {
	if p == nil {
		return nil
	}
	for _, key := range p.Keys() {
		m, ok := b.findDeclaredMethod(key)
		if !ok {
			return fmt.Errorf("catty: policy selects unknown (missing) method %s.%s%s",
				key.Owner, key.Name, key.Descriptor)
		}
		if m.Flags&rt.AccAbstract != 0 {
			return fmt.Errorf("catty: policy target %s.%s%s is abstract; a concrete method is required (impossible selection)",
				key.Owner, key.Name, key.Descriptor)
		}
		engine, _ := p.Engine(key)
		if engine == "aot" {
			if pc, has := aotUnsupported(m); has {
				return fmt.Errorf("catty: AOT lowering is impossible for %s.%s%s: irem has no AOT lowering (pc %d, offset %d, bci %d)",
					key.Owner, key.Name, key.Descriptor, pc, pc, pc)
			}
		}
	}
	return nil
}

func aotUnsupported(m *classfile.Method) (int32, bool) {
	for _, insn := range m.Insns {
		if insn.Op == rt.OpIrem {
			return insn.PC, true
		}
	}
	return 0, false
}

func (b *Build) selectEngines(p *Policy) error {
	for _, c := range b.loaded.Classes {
		for _, m := range c.Methods {
			key := MethodKey{Owner: c.Name, Name: m.Name, Descriptor: m.Descriptor}
			sel := &Selection{EntryIndex: -1}
			forced, hasForced := p.Engine(key)
			switch {
			case m.Flags&rt.AccAbstract != 0:
				sel.Engine = "abstract"
				sel.Reason = "abstract declaration"
			case hasForced:
				switch forced {
				case "interp":
					sel.Engine = "interp"
					sel.Reason = "policy: forced interp"
				case "aot":
					if pc, bad := aotUnsupported(m); bad {
						return fmt.Errorf("catty: AOT lowering is impossible for %s.%s%s: irem has no AOT lowering (pc %d, offset %d, bci %d)",
							c.Name, m.Name, m.Descriptor, pc, pc, pc)
					}
					sel.Engine = "aot"
					sel.Reason = "policy: forced aot"
				}
			case b.opts.Mode == "interp":
				sel.Engine = "interp"
				sel.Reason = "interp: retained bytecode (interp mode)"
			case b.opts.Mode == "aot":
				if pc, bad := aotUnsupported(m); bad {
					return fmt.Errorf("catty: AOT lowering is impossible for %s.%s%s: irem has no AOT lowering (pc %d, offset %d, bci %d)",
						c.Name, m.Name, m.Descriptor, pc, pc, pc)
				}
				sel.Engine = "aot"
				sel.Reason = "aot: lowered to Go"
			default: // auto
				if pc, bad := aotUnsupported(m); bad {
					sel.Engine = "interp"
					sel.Reason = fmt.Sprintf("interp: irem has no AOT lowering at pc %d (offset %d, bci %d)", pc, pc, pc)
				} else {
					sel.Engine = "aot"
					sel.Reason = "aot: lowered to Go"
				}
			}
			if sel.Engine == "aot" {
				heights, ok := b.ctx.HeightsOf(m)
				if !ok {
					return fmt.Errorf("catty: %s.%s%s has no verified stack heights", c.Name, m.Name, m.Descriptor)
				}
				shapes, ok := b.ctx.ShapesOf(m)
				if !ok {
					return fmt.Errorf("catty: %s.%s%s has no verified stack shapes", c.Name, m.Name, m.Descriptor)
				}
				sel.Heights = heights
				sel.Shapes = shapes
				sel.EntryIndex = len(b.aot)
				sel.Symbol = fmt.Sprintf("cattyMethod%d", sel.EntryIndex)
				b.aot = append(b.aot, &aotEntry{class: c, method: m, symbol: sel.Symbol, heights: heights, shapes: shapes})
			}
			b.selections[m] = sel
		}
	}
	return nil
}

// resolveReferences validates every symbolic reference. A reference to an
// unavailable bootstrap facility is a build error. An application reference
// whose absence does not prevent verification is retained in the report and
// resolved lazily by the executing instruction, exactly as Java resolves it.
func (b *Build) resolveReferences() error {
	seenUnresolved := map[string]bool{}
	for _, c := range b.loaded.Classes {
		for _, m := range c.Methods {
			pcOf := map[int32]int32{}
			for _, insn := range m.Insns {
				if insn.Ref >= 0 {
					if _, ok := pcOf[insn.Ref]; !ok {
						pcOf[insn.Ref] = insn.PC
					}
				}
			}
			for index, ref := range m.Refs {
				unresolved, err := b.checkReference(c, m, ref)
				if err != nil {
					return err
				}
				if unresolved == nil {
					continue
				}
				unresolved.PC = int(pcOf[int32(index)])
				key := fmt.Sprintf("%s.%s%s|%d", unresolved.Owner, unresolved.Name, unresolved.Descriptor, unresolved.PC)
				if seenUnresolved[key] {
					continue
				}
				seenUnresolved[key] = true
				b.unresolved = append(b.unresolved, *unresolved)
			}
			for _, h := range m.Handlers {
				if h.CatchType == "" {
					continue
				}
				if !b.classAvailable(h.CatchType) {
					return fmt.Errorf("%s: %s.%s%s catch type %s is unavailable", c.Origin, c.Name, m.Name, m.Descriptor, h.CatchType)
				}
				if !b.isThrowable(h.CatchType) {
					return fmt.Errorf("%s: %s.%s%s catch type %s is not a Throwable", c.Origin, c.Name, m.Name, m.Descriptor, h.CatchType)
				}
			}
		}
	}
	return nil
}

func (b *Build) checkReference(c *classfile.Class, m *classfile.Method, ref rt.Ref) (*UnresolvedRecord, error) {
	switch ref.Kind {
	case rt.RefInt, rt.RefLong, rt.RefFloat, rt.RefDouble, rt.RefString:
		return nil, nil
	case rt.RefClass:
		if !b.classAvailable(ref.Class) {
			return nil, fmt.Errorf("%s: %s.%s%s class constant %s is unavailable",
				c.Origin, c.Name, m.Name, m.Descriptor, ref.Class)
		}
		return nil, nil
	case rt.RefField:
		if _, err := rt.ParseFieldType(ref.Descriptor); err != nil {
			return nil, fmt.Errorf("%s: %s.%s%s field %s.%s: %w", c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, err)
		}
		if rt.BootstrapHasClassFor(b.profile.ID, ref.Owner) {
			if !rt.BootstrapHasFieldFor(b.profile.ID, ref.Owner, ref.Name, ref.Descriptor) {
				return nil, fmt.Errorf("%s: %s.%s%s needs unavailable bootstrap field %s.%s:%s",
					c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, ref.Descriptor)
			}
			return nil, nil
		}
		if reservedClassName(b.profile, ref.Owner) {
			return nil, fmt.Errorf("%s: %s.%s%s needs unavailable bootstrap field %s.%s:%s",
				c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, ref.Descriptor)
		}
		if b.fieldExists(ref.Owner, ref.Name, ref.Descriptor) {
			return nil, nil
		}
		return &UnresolvedRecord{
			Kind: "field", Owner: ref.Owner, Name: ref.Name, Descriptor: ref.Descriptor,
			Referrer: c.Name, ReferrerMethod: m.Name, ReferrerDescriptor: m.Descriptor,
			ClassMissing: b.loaded.ByName[ref.Owner] == nil,
		}, nil
	case rt.RefMethod:
		if _, err := rt.ParseMethodType(ref.Descriptor); err != nil {
			return nil, fmt.Errorf("%s: %s.%s%s method %s.%s: %w", c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, err)
		}
		if rt.BootstrapHasClassFor(b.profile.ID, ref.Owner) {
			if !rt.BootstrapHasMethodFor(b.profile.ID, ref.Owner, ref.Name, ref.Descriptor) {
				return nil, fmt.Errorf("%s: %s.%s%s needs unavailable bootstrap method %s.%s%s",
					c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, ref.Descriptor)
			}
			return nil, nil
		}
		if reservedClassName(b.profile, ref.Owner) {
			return nil, fmt.Errorf("%s: %s.%s%s needs unavailable bootstrap method %s.%s%s",
				c.Origin, c.Name, m.Name, m.Descriptor, ref.Owner, ref.Name, ref.Descriptor)
		}
		if b.methodExists(ref.Owner, ref.Name, ref.Descriptor) {
			return nil, nil
		}
		return &UnresolvedRecord{
			Kind: "method", Owner: ref.Owner, Name: ref.Name, Descriptor: ref.Descriptor,
			Referrer: c.Name, ReferrerMethod: m.Name, ReferrerDescriptor: m.Descriptor,
			ClassMissing: b.loaded.ByName[ref.Owner] == nil,
		}, nil
	}
	return nil, nil
}

func (b *Build) fieldExists(owner, name, descriptor string) bool {
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(class string) bool {
		if class == "" || seen[class] {
			return false
		}
		seen[class] = true
		if rt.BootstrapHasFieldFor(b.profile.ID, class, name, descriptor) {
			return true
		}
		c := b.loaded.ByName[class]
		if c == nil {
			return false
		}
		for _, f := range c.Fields {
			if f.Name == name && f.Descriptor == descriptor {
				return true
			}
		}
		for _, in := range c.Interfaces {
			if visit(in) {
				return true
			}
		}
		return visit(c.Super)
	}
	return visit(owner)
}

func (b *Build) methodExists(owner, name, descriptor string) bool {
	if rt.BootstrapHasClass(owner) {
		return rt.BootstrapHasMethod(owner, name, descriptor)
	}
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(class string) bool {
		if class == "" || seen[class] {
			return false
		}
		seen[class] = true
		if rt.BootstrapHasMethodFor(b.profile.ID, class, name, descriptor) {
			return true
		}
		c := b.loaded.ByName[class]
		if c == nil {
			return false
		}
		for _, m := range c.Methods {
			if m.Name == name && m.Descriptor == descriptor {
				return true
			}
		}
		for _, in := range c.Interfaces {
			if visit(in) {
				return true
			}
		}
		return visit(c.Super)
	}
	return visit(owner)
}

func (b *Build) isThrowable(name string) bool {
	for cur := name; cur != ""; {
		if cur == "java/lang/Throwable" {
			return true
		}
		if c := b.loaded.ByName[cur]; c != nil {
			cur = c.Super
			continue
		}
		if info, ok := rt.BootstrapClassDefFor(b.profile.ID, cur); ok {
			cur = info.Super
			continue
		}
		return false
	}
	return false
}

func (b *Build) buildImage() error {
	classes := make([]*rt.ClassDef, 0, len(b.loaded.Classes))
	for _, c := range b.loaded.Classes {
		def := &rt.ClassDef{
			Name:        c.Name,
			Super:       c.Super,
			Interfaces:  append([]string(nil), c.Interfaces...),
			Flags:       c.Flags,
			NestHost:    c.NestHost,
			NestMembers: append([]string(nil), c.NestMembers...),
		}
		for _, f := range c.Fields {
			fd := rt.FieldDef{Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags, Constant: f.Constant}
			if f.Flags&rt.AccStatic != 0 {
				def.StaticFields = append(def.StaticFields, fd)
			} else {
				def.InstanceFields = append(def.InstanceFields, fd)
			}
		}
		for _, m := range c.Methods {
			sel := b.selections[m]
			md := &rt.MethodDef{
				Owner:      c.Name,
				Name:       m.Name,
				Descriptor: m.Descriptor,
				Flags:      m.Flags,
				Engine:     sel.Engine,
				Reason:     sel.Reason,
				MaxStack:   m.MaxStack,
				MaxLocals:  m.MaxLocals,
				Handlers:   append([]rt.Handler(nil), m.Handlers...),
				Refs:       append([]rt.Ref(nil), m.Refs...),
				EntryIndex: -1,
			}
			switch sel.Engine {
			case "interp":
				md.CodeSize = m.CodeLength
				md.Insns = append([]rt.Insn(nil), m.Insns...)
			case "aot":
				md.EntryIndex = sel.EntryIndex
			case "abstract":
				md.EntryIndex = -1
			}
			def.Methods = append(def.Methods, md)
		}
		classes = append(classes, def)
	}
	image := &rt.Image{
		ABI:     b.profile.ABI,
		Profile: b.profile.ID,
		Plugins: b.profile.Plugins,
		Opcodes: rt.SupportedOpcodeNamesFor(b.profile.ID),
		Classes: classes,
		Main: rt.MethodKey{
			Owner:      strings.ReplaceAll(b.opts.Main, ".", "/"),
			Name:       "main",
			Descriptor: "([Ljava/lang/String;)V",
		},
	}
	image.Threads = b.profile.Threads
	image.Lifecycle = b.profile.Lifecycle
	if b.profile.Plugins {
		image.Capabilities = rt.ProfileCapabilities(b.profile.ID)
		image.Limits = b.profile.Limits
	}
	b.image = image
	return nil
}

// Report types mirror the published Gate 2 build contract.
type Report struct {
	Profile      string             `json:"profile"`
	ImageABI     int                `json:"imageABI"`
	Plugins      bool               `json:"plugins"`
	Threads      bool               `json:"threads"`
	Lifecycle    bool               `json:"lifecycle"`
	Capabilities []string           `json:"capabilities,omitempty"`
	PluginLimits *rt.PluginLimits   `json:"pluginLimits,omitempty"`
	Main         string             `json:"main"`
	Toolchains   Toolchains         `json:"toolchains"`
	Inputs       []Input            `json:"inputs"`
	Classes      []ClassRecord      `json:"classes"`
	Methods      []MethodRecord     `json:"methods"`
	Opcodes      []string           `json:"opcodes"`
	Bootstrap    BootstrapInventory `json:"bootstrap"`
	Unresolved   []UnresolvedRecord `json:"unresolved,omitempty"`
}

type Toolchains struct {
	Go   string `json:"go"`
	Java string `json:"java,omitempty"`
}

type ClassRecord struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type MethodRecord struct {
	Owner             string           `json:"owner"`
	Name              string           `json:"name"`
	Descriptor        string           `json:"descriptor"`
	Engine            string           `json:"engine"`
	Reason            string           `json:"reason"`
	RetainedCodeBytes int              `json:"retainedCodeBytes"`
	Generated         *GeneratedRecord `json:"generated,omitempty"`
}

type GeneratedRecord struct {
	File   string `json:"file"`
	Symbol string `json:"symbol"`
}

// BootstrapInventory publishes the supported bootstrap classes and members.
type BootstrapInventory struct {
	Classes []BootstrapClassRecord  `json:"classes"`
	Methods []BootstrapMemberRecord `json:"methods"`
	Fields  []BootstrapMemberRecord `json:"fields"`
}

type BootstrapClassRecord struct {
	Name       string   `json:"name"`
	Super      string   `json:"super,omitempty"`
	Interfaces []string `json:"interfaces,omitempty"`
	Flags      uint16   `json:"flags,omitempty"`
}

type BootstrapMemberRecord struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
}

// UnresolvedRecord reports one retained application reference that Java
// resolves lazily when its instruction executes.
type UnresolvedRecord struct {
	Kind               string `json:"kind"`
	Owner              string `json:"owner"`
	Name               string `json:"name"`
	Descriptor         string `json:"descriptor"`
	Referrer           string `json:"referrer"`
	ReferrerMethod     string `json:"referrerMethod"`
	ReferrerDescriptor string `json:"referrerDescriptor"`
	PC                 int    `json:"pc"`
	ClassMissing       bool   `json:"classMissing,omitempty"`
}

func (b *Build) buildReport() {
	r := &Report{
		Profile:   b.profile.ID,
		ImageABI:  b.profile.ABI,
		Plugins:   b.profile.Plugins,
		Threads:   b.profile.Threads,
		Lifecycle: b.profile.Lifecycle,
		Main:      b.opts.Main,
		Toolchains: Toolchains{
			Go: runtime.Version(),
		},
		Inputs:     append([]Input(nil), b.loaded.Inputs...),
		Opcodes:    rt.SupportedOpcodeNamesFor(b.profile.ID),
		Unresolved: append([]UnresolvedRecord(nil), b.unresolved...),
	}
	if b.profile.Plugins {
		r.Capabilities = rt.ProfileCapabilities(b.profile.ID)
		r.PluginLimits = b.profile.Limits
	}
	for _, c := range b.loaded.Classes {
		r.Classes = append(r.Classes, ClassRecord{Name: c.Name, Origin: c.Origin})
	}
	for _, c := range b.loaded.Classes {
		for _, m := range c.Methods {
			sel := b.selections[m]
			rec := MethodRecord{
				Owner:      c.Name,
				Name:       m.Name,
				Descriptor: m.Descriptor,
				Engine:     sel.Engine,
				Reason:     sel.Reason,
			}
			if sel.Engine == "interp" {
				rec.RetainedCodeBytes = m.CodeLength
			}
			if sel.Engine == "aot" && b.opts.GeneratedSource != "" {
				rec.Generated = &GeneratedRecord{File: generatedFileName, Symbol: sel.Symbol}
			}
			r.Methods = append(r.Methods, rec)
		}
	}
	for _, intrinsic := range rt.BootstrapMethodsFor(b.profile.ID) {
		r.Methods = append(r.Methods, MethodRecord{
			Owner:             intrinsic.Owner,
			Name:              intrinsic.Name,
			Descriptor:        intrinsic.Descriptor,
			Engine:            "intrinsic",
			Reason:            "bootstrap intrinsic",
			RetainedCodeBytes: 0,
		})
	}
	for _, class := range rt.BootstrapClassesFor(b.profile.ID) {
		r.Bootstrap.Classes = append(r.Bootstrap.Classes, BootstrapClassRecord{
			Name: class.Name, Super: class.Super, Interfaces: class.Interfaces, Flags: class.Flags,
		})
	}
	for _, method := range rt.BootstrapMethodsFor(b.profile.ID) {
		r.Bootstrap.Methods = append(r.Bootstrap.Methods, BootstrapMemberRecord{
			Owner: method.Owner, Name: method.Name, Descriptor: method.Descriptor,
		})
	}
	for _, field := range rt.BootstrapFieldsFor(b.profile.ID) {
		r.Bootstrap.Fields = append(r.Bootstrap.Fields, BootstrapMemberRecord{
			Owner: field.Owner, Name: field.Name, Descriptor: field.Descriptor,
		})
	}
	b.report = r
}
