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

func TestFailedDescriptorValidationLeavesArrayCacheUntouched(t *testing.T) {
	r := gate3Runtime(t)
	before := len(r.arrayClasses) + len(r.primitiveArrayClass)
	def := returnIntDef("ext/Bad", "app/Base")
	def.InstanceFields = []FieldDef{{Name: "value", Descriptor: "[Lext/Missing;"}}
	if _, err := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{def}, Defined: []string{"ext/Bad"}}); err == nil {
		t.Fatal("missing array component was committed")
	}
	after := len(r.arrayClasses) + len(r.primitiveArrayClass)
	if after != before {
		t.Fatalf("failed admission changed the array intern table: %d -> %d", before, after)
	}
	if len(r.pluginLoaders) != 0 {
		t.Fatal("failed descriptor validation published a loader")
	}
}
