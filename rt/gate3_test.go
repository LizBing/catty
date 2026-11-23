package rt

import (
	"fmt"
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

func TestLoadingConstraintPendingInitiation(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	typ, err := ParseFieldType("Lext/T;")
	if err != nil {
		t.Fatal(err)
	}
	mt := &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
	// Both endpoints are unresolved, so the equality is staged pending without
	// forcing a descriptor leaf lookup.
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("pending equality was reported as a conflict")
	}
	if len(r.initiations) != 0 || len(r.constraints) != 1 {
		t.Fatalf("pending staging published initialization state: init=%d constraints=%d",
			len(r.initiations), len(r.constraints))
	}
	if r.recordInitiation(a, "ext/T", ta) {
		t.Fatal("first initiation reported a conflict")
	}
	keyB := initiationKey{Loader: b, Name: "ext/T"}
	if !r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("later lookup did not observe the pending constraint")
	}
	if _, ok := r.initiations[keyB]; ok {
		t.Fatal("conflicting late initiation was published")
	}
}

func TestLoadingConstraintTransitiveComponent(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	c := &Loader{ID: 4, Name: "constraint-c", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	object := r.mustClass("java/lang/Object")
	tc := &Class{Name: "ext/T", Loader: c, Super: object, MethodIndex: map[string]*Method{}}
	c.Classes["ext/T"] = tc
	typ, _ := ParseFieldType("Lext/T;")
	mt := &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if r.imposeMethodConstraints(a, b, mt) || r.imposeMethodConstraints(b, c, mt) {
		t.Fatal("pending transitive equality reported a conflict")
	}
	if len(r.constraints) != 2 {
		t.Fatalf("transitive staging published %d edges, want 2", len(r.constraints))
	}
	if r.recordInitiation(a, "ext/T", ta) {
		t.Fatal("first initiation reported a conflict")
	}
	// C reaches A only through the still-unresolved B endpoint. The connected
	// component check must follow both equality edges and reject C's class.
	if !r.recordInitiation(c, "ext/T", tc) {
		t.Fatal("transitive conflict through an unresolved endpoint was missed")
	}
	keyC := initiationKey{Loader: c, Name: "ext/T"}
	if _, ok := r.initiations[keyC]; ok {
		t.Fatal("conflicting transitive initiation was published")
	}
	if _, ok := r.initiations[initiationKey{Loader: a, Name: "ext/T"}]; !ok {
		t.Fatal("prior initiation was lost by the refused operation")
	}
	if len(r.constraints) != 2 {
		t.Fatal("refused initiation changed the constraint graph")
	}
	_ = tb
}

func TestLoadingConstraintLateLookupPublishesSharedClass(t *testing.T) {
	r, a, b, _, _ := constraintRuntime(t)
	object := r.mustClass("java/lang/Object")
	shared := &Class{Name: "app/Shared", Loader: r.Loader, Super: object, MethodIndex: map[string]*Method{}}
	r.Loader.Classes["app/Shared"] = shared
	r.Classes["app/Shared"] = shared
	typ, _ := ParseFieldType("Lapp/Shared;")
	mt := &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("pending shared-class equality reported a conflict")
	}
	if r.recordInitiation(a, "app/Shared", shared) || r.recordInitiation(b, "app/Shared", shared) {
		t.Fatal("shared parent-visible class conflicted")
	}
	if r.initiations[initiationKey{Loader: a, Name: "app/Shared"}] != shared ||
		r.initiations[initiationKey{Loader: b, Name: "app/Shared"}] != shared {
		t.Fatal("shared initiations are not loader-specific")
	}
}

func TestLoadingConstraintConflictPreservesPriorState(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	if r.recordInitiation(a, "ext/T", ta) || r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("unconstrained initiations conflicted")
	}
	beforeInit := len(r.initiations)
	beforeConstraints := len(r.constraints)
	typ, _ := ParseFieldType("Lext/T;")
	mt := &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if !r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("resolved distinct classes did not conflict at staging")
	}
	if len(r.initiations) != beforeInit || len(r.constraints) != beforeConstraints {
		t.Fatal("refused staging changed published constraint state")
	}
}

func TestLoadingConstraintPrimitiveDescriptorsStayOutsideGraph(t *testing.T) {
	r, a, b, _, _ := constraintRuntime(t)
	intType, _ := ParseFieldType("I")
	arrayType, _ := ParseFieldType("[[I")
	void := Type{Kind: 'V', Descriptor: "V"}
	mt := &MethodType{Params: []Type{intType, arrayType}, Return: void}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("primitive descriptors conflicted")
	}
	if len(r.constraints) != 0 || len(r.initiations) != 0 {
		t.Fatal("primitive descriptors entered the equality graph")
	}
}

func TestLoadingConstraintArrayLeaf(t *testing.T) {
	r, a, b, _, _ := constraintRuntime(t)
	arrayType, _ := ParseFieldType("[[Lext/T;")
	mt := &MethodType{Params: []Type{arrayType}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("pending array leaf reported a conflict")
	}
	if len(r.constraints) != 1 {
		t.Fatalf("array leaf staged %d edges, want 1", len(r.constraints))
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
	// Admission-time conflicts are visible only when both sides have already
	// initiated the leaf; pending equalities stay lazy until a later lookup.
	if r.recordInitiation(a, "ext/T", ta) || r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("initial initiations conflicted")
	}
	tx := r.newConstraintTx()
	if err := r.checkOverrideConstraints([]*Class{cand}, tx); err == nil {
		t.Fatal("override preparation did not report the leaf-type conflict")
	}
	if len(r.constraints) != 0 {
		t.Fatal("rejected override preparation published constraints")
	}
}

func TestLoadingConstraintArrayLeafConflict(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	arrayType, _ := ParseFieldType("[[Lext/T;")
	mt := &MethodType{Params: []Type{arrayType}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("pending array leaf reported a conflict")
	}
	if r.recordInitiation(a, "ext/T", ta) {
		t.Fatal("first array leaf initiation conflicted")
	}
	if !r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("distinct array leaf classes did not conflict on later initiation")
	}
	if _, ok := r.initiations[initiationKey{Loader: b, Name: "ext/T"}]; ok {
		t.Fatal("conflicting array leaf initiation was published")
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
	if r.recordInitiation(a, "ext/T", ta) || r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("initial initiations conflicted")
	}
	tx := r.newConstraintTx()
	if err := r.checkOverrideConstraints([]*Class{cand}, tx); err == nil {
		t.Fatal("inherited interface implementation did not trigger preparation constraints")
	}
	if len(r.constraints) != 0 {
		t.Fatal("rejected inherited-interface preparation published constraints")
	}
}

func TestWorkLimitRefusalLeavesRunningLoaderUsable(t *testing.T) {
	r := gate3Runtime(t)
	def := returnIntDef("ext/Existing", "app/Base")
	loader, thrown := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: []*ClassDef{def}, Defined: []string{"ext/Existing"},
	})
	if thrown != nil {
		t.Fatal(thrown)
	}
	existing, ok := loader.Local("ext/Existing")
	if !ok || existing.Mirror == nil {
		t.Fatal("existing plugin class is not usable")
	}
	r.admit = func(*AdmissionRequest) (*AdmissionResult, error) {
		return nil, &PluginError{Code: "LIMIT", Message: "injected verification work limit"}
	}
	th := r.MainThread()
	arr, thrown := th.NewPrimitiveArray(nil, 0, 8, 0)
	if thrown != nil {
		t.Fatal(thrown)
	}
	out := intrinsicPluginsRegister(th, &Method{Owner: existing}, nil,
		[]Value{RefValue(arr), RefValue(nil)})
	data := ThrowableOf(out.Thrown)
	if data == nil || data.Code != "LIMIT" {
		t.Fatalf("expected LIMIT refusal, got %#v", out.Thrown)
	}
	if len(r.pluginLoaders) != 1 {
		t.Fatalf("refusal changed published loaders: %d", len(r.pluginLoaders))
	}
	if current, ok := loader.Local("ext/Existing"); !ok || current != existing || existing.Mirror == nil {
		t.Fatal("refusal made the existing loader unusable")
	}
}

func hierarchyChain(prefix string, first, count int, base string) []*ClassDef {
	defs := make([]*ClassDef, 0, count)
	for i := 0; i < count; i++ {
		n := first + i
		super := base
		if n > 0 {
			super = fmt.Sprintf("%s/%03d", prefix, n-1)
		} else if super == "" {
			super = "java/lang/Object"
		}
		defs = append(defs, &ClassDef{Name: fmt.Sprintf("%s/%03d", prefix, n), Super: super})
	}
	return defs
}

func hierarchyNames(defs []*ClassDef) []string {
	out := make([]string, 0, len(defs))
	for _, d := range defs {
		out = append(out, d.Name)
	}
	return out
}

func throwableCode(o *Object) string {
	if data := ThrowableOf(o); data != nil {
		return data.Code
	}
	return ""
}

func TestPluginHierarchyDepthExactBoundary(t *testing.T) {
	for _, tc := range []struct {
		edges int
		ok    bool
	}{
		{128, true},
		{129, false},
	} {
		r := gate3Runtime(t)
		defs := hierarchyChain("depth", 0, tc.edges, "")
		_, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: defs, Defined: hierarchyNames(defs)})
		if tc.ok {
			if thrown != nil {
				t.Fatalf("%d edges rejected: %v", tc.edges, thrown)
			}
			if len(r.pluginLoaders) != 1 {
				t.Fatalf("%d edges did not publish one loader", tc.edges)
			}
			continue
		}
		if code := throwableCode(thrown); code != "LIMIT" {
			t.Fatalf("%d edges: expected LIMIT, got %q", tc.edges, code)
		}
		if len(r.pluginLoaders) != 0 || r.nextLoaderID != 1 {
			t.Fatalf("%d edges published state after refusal", tc.edges)
		}
	}

	// Ancestor-first ClassDef order must not hide a longer path through an
	// already completed node.
	r := gate3Runtime(t)
	defs := hierarchyChain("order", 0, 129, "")
	for i, j := 0, len(defs)-1; i < j; i, j = i+1, j-1 {
		defs[i], defs[j] = defs[j], defs[i]
	}
	_, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: defs, Defined: hierarchyNames(defs)})
	if code := throwableCode(thrown); code != "LIMIT" {
		t.Fatalf("ancestor-first 129-edge chain was not refused: %q", code)
	}
}

func TestPluginHierarchyDepthSpansParentLoader(t *testing.T) {
	r := gate3Runtime(t)
	parentDefs := hierarchyChain("span", 0, 120, "")
	parent, thrown := r.installPlugin(r.Loader, &AdmissionResult{
		Classes: parentDefs, Defined: hierarchyNames(parentDefs)})
	if thrown != nil {
		t.Fatal(thrown)
	}
	childDefs := []*ClassDef{}
	for n := 120; n <= 127; n++ {
		childDefs = append(childDefs, &ClassDef{
			Name: fmt.Sprintf("span/%03d", n), Super: fmt.Sprintf("span/%03d", n-1)})
	}
	child, thrown := r.installPlugin(parent, &AdmissionResult{
		Classes: childDefs, Defined: hierarchyNames(childDefs)})
	if thrown != nil {
		t.Fatalf("128-edge path spanning a parent loader was rejected: %v", thrown)
	}
	grandDefs := []*ClassDef{{Name: "span/128", Super: "span/127"}}
	before := len(r.pluginLoaders)
	_, thrown = r.installPlugin(child, &AdmissionResult{
		Classes: grandDefs, Defined: hierarchyNames(grandDefs)})
	if code := throwableCode(thrown); code != "LIMIT" {
		t.Fatalf("129-edge path spanning parent loaders: got %q", code)
	}
	if len(r.pluginLoaders) != before {
		t.Fatal("refused parent-spanning admission published a loader")
	}
	if _, ok := child.Local("span/127"); !ok {
		t.Fatal("refusal damaged the existing parent loader")
	}
	if _, ok := parent.Local("span/119"); !ok {
		t.Fatal("refusal damaged the first parent loader")
	}
}

func TestPluginHierarchyCycleRejected(t *testing.T) {
	r := gate3Runtime(t)
	defs := []*ClassDef{
		{Name: "cycle/A", Super: "cycle/B"},
		{Name: "cycle/B", Super: "cycle/A"},
	}
	_, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: defs, Defined: hierarchyNames(defs)})
	if code := throwableCode(thrown); code != "VERIFY" {
		t.Fatalf("cycle admission code %q", code)
	}
	if len(r.pluginLoaders) != 0 {
		t.Fatal("cycle admission published a loader")
	}
}
