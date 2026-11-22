package rt

import (
	"testing"
)

func gate3Runtime(t *testing.T) *Runtime {
	t.Helper()
	img := &Image{
		ABI:     Gate3ABI,
		Profile: Gate3ProfileID,
		Plugins: true,
		Limits:  Gate3PluginLimits,
		Classes: []*ClassDef{
			{Name: "app/Base", Super: "java/lang/Object"},
			{Name: "app/Step", Super: "java/lang/Object", Flags: AccInterface | AccAbstract},
		},
	}
	r := NewRuntime(img)
	if err := r.Install(); err != nil {
		t.Fatal(err)
	}
	return r
}

func returnIntDef(name, super string) *ClassDef {
	return &ClassDef{
		Name: name, Super: super,
		Methods: []*MethodDef{{
			Owner: name, Name: "value", Descriptor: "()I", Flags: AccPublic | AccStatic,
			Engine: "interp", MaxStack: 1, MaxLocals: 0, CodeSize: 2, EntryIndex: -1,
			Insns: []Insn{{Op: OpIconst0, PC: 0, Len: 1, Ref: -1}, {Op: OpIreturn, PC: 1, Len: 1, Ref: -1}},
		}},
	}
}

func TestLoaderIdentityAndArrays(t *testing.T) {
	r := gate3Runtime(t)
	a, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/Worker", "app/Base")}, Defined: []string{"ext/Worker"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/Worker", "app/Base")}, Defined: []string{"ext/Worker"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ca, ok := a.Local("ext/Worker")
	if !ok {
		t.Fatal("first plugin class missing")
	}
	cb, ok := b.Local("ext/Worker")
	if !ok {
		t.Fatal("second plugin class missing")
	}
	if ca == cb || ca.Mirror == cb.Mirror {
		t.Fatal("identical bytes produced shared class identity")
	}
	if a == b || a.ID == b.ID {
		t.Fatal("plugin loaders share identity")
	}
	if a.ID != 2 || b.ID != 3 {
		t.Fatalf("plugin loader IDs are not stable and sequential: %d, %d", a.ID, b.ID)
	}
	if ca.Super != r.Classes["app/Base"] {
		t.Fatal("plugin superclass was not delegated to the application loader")
	}
	if r.arrayClassForComponent(ca) == r.arrayClassForComponent(cb) {
		t.Fatal("sibling array classes were interned together")
	}
	if r.arrayClassForComponent(ca) != r.arrayClassForComponent(ca) {
		t.Fatal("array identity is not stable")
	}
	if r.arrayClassForComponent(r.PrimitiveClass('I')) != r.arrayClassForComponent(r.PrimitiveClass('I')) {
		t.Fatal("primitive array identity is not stable")
	}
	if r.samePackage(ca, r.Classes["app/Base"]) {
		t.Fatal("sibling textual package compared equal across loaders")
	}
}

func TestPluginQuotaAndRollback(t *testing.T) {
	r := gate3Runtime(t)
	limits := *Gate3PluginLimits
	limits.Loaders = 1
	r.Image.Limits = &limits
	if _, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/One", "app/Base")}, Defined: []string{"ext/One"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/Two", "app/Base")}, Defined: []string{"ext/Two"},
	}); err == nil {
		t.Fatal("loader quota was not enforced")
	}
	loaders := len(r.pluginLoaders)
	bytesBefore := r.pluginClassBytes
	idBefore := r.nextLoaderID
	if _, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/Bad", "missing/Super")}, Defined: []string{"ext/Bad"},
	}); err == nil {
		t.Fatal("missing superclass was committed")
	}
	if len(r.pluginLoaders) != loaders || r.pluginClassBytes != bytesBefore || r.nextLoaderID != idBefore {
		t.Fatal("failed admission mutated loader registry or quotas")
	}
}

func TestPluginMethodRunsThroughInterpreter(t *testing.T) {
	r := gate3Runtime(t)
	loader, err := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{returnIntDef("ext/Worker", "app/Base")}, Defined: []string{"ext/Worker"},
	})
	if err != nil {
		t.Fatal(err)
	}
	c, _ := loader.Local("ext/Worker")
	m := c.MethodIndex["value()I"]
	if m == nil || m.Engine != "interp" {
		t.Fatalf("plugin method is not interpreted: %#v", m)
	}
	th := r.MainThread()
	out := th.Call(m, nil, nil)
	if out.Thrown != nil {
		t.Fatalf("plugin method threw: %v", out.Thrown)
	}
	if AsInt(out.Value) != 0 {
		t.Fatalf("unexpected value %d", AsInt(out.Value))
	}
}

func TestProfilePairAndInventory(t *testing.T) {
	cases := []struct {
		profile string
		abi     int
		plugins bool
		ok      bool
	}{
		{Gate2ProfileID, Gate2ABI, false, true},
		{Gate3ProfileID, Gate3ABI, true, true},
		{Gate3ProfileID, Gate2ABI, true, false},
		{Gate2ProfileID, Gate3ABI, false, false},
		{"catty-unknown", Gate3ABI, true, false},
	}
	for _, c := range cases {
		img := &Image{ABI: c.abi, Profile: c.profile, Plugins: c.plugins}
		r := NewRuntime(img)
		err := r.Install()
		if c.ok && err != nil {
			t.Fatalf("%s/%d: unexpected install error: %v", c.profile, c.abi, err)
		}
		if !c.ok && err == nil {
			t.Fatalf("%s/%d: mixed or unknown profile pair was installed", c.profile, c.abi)
		}
	}
	if !BootstrapHasClassFor(Gate3ProfileID, "catty/runtime/Plugins") {
		t.Fatal("Gate 3 inventory misses Plugins")
	}
	if BootstrapHasClassFor(Gate2ProfileID, "catty/runtime/Plugins") {
		t.Fatal("Gate 2 inventory contains the Gate 3 API")
	}
	if !BootstrapHasMethodFor(Gate3ProfileID, "java/lang/Class", "forName",
		"(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;") {
		t.Fatal("Gate 3 inventory misses Class.forName")
	}
}

func TestResolutionCacheScopedToRuntimeClassAndCP(t *testing.T) {
	r := gate3Runtime(t)
	defs := []*ClassDef{
		{Name: "ext/Dep", Super: "java/lang/Object"},
		returnIntDef("ext/Link", "java/lang/Object"),
	}
	a, err := r.installPlugin(r.Loader, &AdmissionResult{Classes: defs, Defined: []string{"ext/Dep", "ext/Link"}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{
		{Name: "ext/Dep", Super: "java/lang/Object"}, returnIntDef("ext/Link", "java/lang/Object")},
		Defined: []string{"ext/Dep", "ext/Link"}})
	if err != nil {
		t.Fatal(err)
	}
	ca, _ := a.Local("ext/Link")
	cb, _ := b.Local("ext/Link")
	ma, mb := ca.MethodIndex["value()I"], cb.MethodIndex["value()I"]
	th := r.MainThread()
	_, thrownA1 := th.resolveFieldSymbol(ma, 1, "ext/Dep", "value", "I")
	_, thrownA2 := th.resolveFieldSymbol(ma, 1, "ext/Dep", "value", "I")
	_, thrownB1 := th.resolveFieldSymbol(mb, 1, "ext/Dep", "value", "I")
	if thrownA1 == nil || thrownB1 == nil {
		t.Fatal("missing field did not fail")
	}
	if thrownA1 != thrownA2 {
		t.Fatal("resolution failure is not sticky for one class/CP")
	}
	if thrownA1 == thrownB1 {
		t.Fatal("resolution cache is shared across runtime classes")
	}
	_, thrownA3 := th.resolveFieldSymbol(ma, 2, "ext/Dep", "value", "I")
	if thrownA3 == thrownA1 {
		t.Fatal("resolution cache ignores the constant-pool entry")
	}
	if len(r.resolutions) < 3 {
		t.Fatalf("expected scoped cache entries, got %d", len(r.resolutions))
	}
}

func TestDeferredDescriptorLeavesArrayCacheUntouched(t *testing.T) {
	r := gate3Runtime(t)
	before := len(r.arrayClasses) + len(r.primitiveArrayClass)
	def := returnIntDef("ext/Deferred", "app/Base")
	def.InstanceFields = []FieldDef{{Name: "value", Descriptor: "[Lext/Missing;"}}
	loader, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{def}, Defined: []string{"ext/Deferred"}})
	if thrown != nil {
		t.Fatalf("deferred field descriptor rejected: %v", thrown)
	}
	if loader == nil {
		t.Fatal("deferred descriptor did not publish a loader")
	}
	if after := len(r.arrayClasses) + len(r.primitiveArrayClass); after != before {
		t.Fatalf("deferred descriptor changed the array intern table: %d -> %d", before, after)
	}
	if c, ok := loader.Local("ext/Deferred"); !ok || c.Mirror == nil {
		t.Fatal("deferred class was not installed")
	}
}

func constraintRuntime(t *testing.T) (*Runtime, *Loader, *Loader, *Class, *Class) {
	t.Helper()
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	a := &Loader{ID: 2, Name: "constraint-a", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	b := &Loader{ID: 3, Name: "constraint-b", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	ta := &Class{Name: "ext/T", Loader: a, Super: object, MethodIndex: map[string]*Method{}}
	tb := &Class{Name: "ext/T", Loader: b, Super: object, MethodIndex: map[string]*Method{}}
	a.Classes["ext/T"] = ta
	b.Classes["ext/T"] = tb
	return r, a, b, ta, tb
}

func TestLoadingConstraintTransition(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	typ, err := ParseFieldType("Lext/T;")
	if err != nil {
		t.Fatal(err)
	}
	mt := &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if !r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("two distinct local leaf classes did not conflict")
	}
	if len(r.initiations) != 0 || len(r.constraints) != 0 {
		t.Fatalf("conflicting constraint published state: init=%d constraints=%d", len(r.initiations), len(r.constraints))
	}
	if r.recordInitiation(a, "ext/T", ta) {
		t.Fatal("first initiation reported a conflict")
	}
	keyB := initiationKey{Loader: b, Name: "ext/T"}
	r.constraints[makeConstraintKey(initiationKey{Loader: a, Name: "ext/T"}, keyB)] = struct{}{}
	if !r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("later lookup did not observe the pending constraint")
	}
	if _, ok := r.initiations[keyB]; ok {
		t.Fatal("conflicting late initiation was published")
	}
}

func TestLoadingConstraintPositiveAndArrayLeaf(t *testing.T) {
	r, a, b, _, _ := constraintRuntime(t)
	object := r.mustClass("java/lang/Object")
	shared := &Class{Name: "app/Shared", Loader: r.Loader, Super: object, MethodIndex: map[string]*Method{}}
	r.Loader.Classes["app/Shared"] = shared
	r.Classes["app/Shared"] = shared
	sharedType, _ := ParseFieldType("Lapp/Shared;")
	arrayType, _ := ParseFieldType("[Lapp/Shared;")
	mt := &MethodType{Params: []Type{sharedType, arrayType}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("shared parent-visible leaf types conflicted")
	}
	if len(r.constraints) == 0 {
		t.Fatal("successful constraints were not published")
	}
	// Both loaders must resolve the shared type; the array leaf is the same.
	if r.initiations[initiationKey{Loader: a, Name: "app/Shared"}] != shared ||
		r.initiations[initiationKey{Loader: b, Name: "app/Shared"}] != shared {
		t.Fatal("shared initiations are not loader-specific")
	}
}

func TestOverridePreparationUsesLeafTypes(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	object := r.mustClass("java/lang/Object")
	// Unit-level preparation: a candidate method in loader a declares ext/T in
	// its descriptor while an inherited declaration in loader b resolves the
	// same textual leaf to a different class. Parent-first visibility makes
	// this unreachable through normal public registration when a shared parent
	// defines the name; the constraint transition itself is tested here.
	declType := Type{Kind: 'L', ClassName: "ext/T", Descriptor: "Lext/T;"}
	parent := &Class{Name: "ext/Parent", Loader: b, Super: object, MethodIndex: map[string]*Method{}}
	p := &Method{Owner: parent, Name: "m", Descriptor: "(Lext/T;)V", Flags: AccPublic,
		Type: &MethodType{Params: []Type{declType}, Return: Type{Kind: 'V', Descriptor: "V"}}, ParamClass: []*Class{tb}}
	parent.addMethod(p)
	cand := &Class{Name: "ext/Candidate", Loader: a, Super: parent, MethodIndex: map[string]*Method{}}
	m := &Method{Owner: cand, Name: "m", Descriptor: "(Lext/T;)V", Flags: AccPublic,
		Type: &MethodType{Params: []Type{declType}, Return: Type{Kind: 'V', Descriptor: "V"}}, ParamClass: []*Class{ta}}
	cand.addMethod(m)
	if err := r.checkOverrideConstraints([]*Class{cand}); err == nil {
		t.Fatal("override preparation did not report the leaf-type conflict")
	}
}

func TestLoadingConstraintArrayLeafConflict(t *testing.T) {
	r, a, b, _, _ := constraintRuntime(t)
	arrayType, _ := ParseFieldType("[[Lext/T;")
	mt := &MethodType{Params: []Type{arrayType}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if !r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("distinct array leaf classes did not conflict")
	}
	if len(r.initiations) != 0 || len(r.constraints) != 0 {
		t.Fatal("conflicting array constraint published state")
	}
}

func TestOverridePreparationInheritedInterface(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	object := r.mustClass("java/lang/Object")
	iface := &Class{Name: "ext/IFace", Loader: b, Super: object,
		Flags: AccInterface | AccAbstract | AccPublic, MethodIndex: map[string]*Method{}}
	declType := Type{Kind: 'L', ClassName: "ext/T", Descriptor: "Lext/T;"}
	p := &Method{Owner: iface, Name: "m", Descriptor: "(Lext/T;)V", Flags: AccPublic | AccAbstract,
		Type: &MethodType{Params: []Type{declType}, Return: Type{Kind: 'V', Descriptor: "V"}}, ParamClass: []*Class{tb}}
	iface.addMethod(p)
	cand := &Class{Name: "ext/Impl", Loader: a, Super: object,
		Interfaces: []*Class{iface}, MethodIndex: map[string]*Method{}}
	m := &Method{Owner: cand, Name: "m", Descriptor: "(Lext/T;)V", Flags: AccPublic,
		Type: &MethodType{Params: []Type{declType}, Return: Type{Kind: 'V', Descriptor: "V"}}, ParamClass: []*Class{ta}}
	cand.addMethod(m)
	if err := r.checkOverrideConstraints([]*Class{cand}); err == nil {
		t.Fatal("inherited interface implementation did not trigger preparation constraints")
	}
}
