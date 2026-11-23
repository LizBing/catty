package rt

import "testing"

// testHiddenClass adds a manually built class to a loader. Runtime unit tests
// use these loaders to reach mixed-loader configurations that fixed
// parent-first visibility makes unreachable through the Java API.
func testHiddenClass(loader *Loader, name string, flags uint16, super *Class) *Class {
	c := &Class{Name: name, Flags: flags, Loader: loader, Super: super, MethodIndex: map[string]*Method{}}
	loader.Classes[name] = c
	return c
}

func testHiddenMethod(owner *Class, name, desc string, flags uint16) *Method {
	mt, err := ParseMethodType(desc)
	if err != nil {
		panic(err)
	}
	m := &Method{Owner: owner, Name: name, Descriptor: desc, Flags: flags, Type: mt}
	if err := owner.addMethod(m); err != nil {
		panic(err)
	}
	return m
}

func testPluginVoidDef(owner, name, desc string, maxLocals int) *MethodDef {
	return &MethodDef{
		Owner: owner, Name: name, Descriptor: desc, Flags: AccPublic,
		Engine: "interp", MaxStack: 0, MaxLocals: maxLocals, CodeSize: 1, EntryIndex: -1,
		Insns: []Insn{{Op: OpReturn, PC: 0, Len: 1, Ref: -1}},
	}
}

func testConstraintMethodType(t *testing.T, desc string) *MethodType {
	t.Helper()
	typ, err := ParseFieldType(desc)
	if err != nil {
		t.Fatal(err)
	}
	return &MethodType{Params: []Type{typ}, Return: Type{Kind: 'V', Descriptor: "V"}}
}

// TestLoadingConstraintMergedGroupsConflict reproduces the staging failure
// where joining two pending groups hides a conflict because both new endpoints
// are unresolved: A initiated X1, D initiated a distinct X2 with the same binary
// name, A = B and C = D are pending at B and C, and B = C would join the groups.
func TestLoadingConstraintMergedGroupsConflict(t *testing.T) {
	r, a, b, ta, _ := constraintRuntime(t)
	object := r.mustClass("java/lang/Object")
	// Unique loader IDs keep makeConstraintKey canonical; B and C are the
	// unresolved endpoints whose merge triggers the conflict.
	c := &Loader{ID: 4, Name: "constraint-c", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	d := &Loader{ID: 5, Name: "constraint-d", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	x2 := &Class{Name: "ext/T", Loader: d, Super: object, MethodIndex: map[string]*Method{}}
	d.Classes["ext/T"] = x2
	mt := testConstraintMethodType(t, "Lext/T;")

	if r.recordInitiation(a, "ext/T", ta) {
		t.Fatal("A's initiation conflicted")
	}
	if r.recordInitiation(d, "ext/T", x2) {
		t.Fatal("D's independent initiation conflicted")
	}
	if r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("A = B was rejected while both groups were disjoint")
	}
	if r.imposeMethodConstraints(c, d, mt) {
		t.Fatal("C = D was rejected while both groups were disjoint")
	}
	if len(r.constraints) != 2 || len(r.initiations) != 2 {
		t.Fatalf("pending groups are not disjoint: edges=%d initiations=%d", len(r.constraints), len(r.initiations))
	}

	// Both merged endpoints are unresolved. The operation must still inspect
	// the full merged equality component and reject the incompatible classes.
	if !r.imposeMethodConstraints(b, c, mt) {
		t.Fatal("joining two incompatible pending groups was accepted")
	}
	if len(r.constraints) != 2 {
		t.Fatalf("refused merge published %d edges, want 2", len(r.constraints))
	}
	if len(r.initiations) != 2 {
		t.Fatalf("refused merge published %d initiations, want 2", len(r.initiations))
	}
	keyB := initiationKey{Loader: b, Name: "ext/T"}
	keyC := initiationKey{Loader: c, Name: "ext/T"}
	if _, ok := r.constraints[makeConstraintKey(keyB, keyC)]; ok {
		t.Fatal("refused merge published the B = C edge")
	}
	// The prior groups stay usable: B may still resolve to A's class.
	if r.recordInitiation(b, "ext/T", ta) {
		t.Fatal("valid initiation after the refused merge conflicted")
	}
}

// TestLoadingConstraintMergedGroupsSameClass covers the accepted merge: both
// groups already resolve their binary name to the same class, so the joined
// component agrees and a later initiation must keep agreeing.
func TestLoadingConstraintMergedGroupsSameClass(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	c := &Loader{ID: 4, Name: "constraint-c", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	d := &Loader{ID: 5, Name: "constraint-d", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	mt := testConstraintMethodType(t, "Lext/T;")

	if r.recordInitiation(a, "ext/T", ta) || r.recordInitiation(d, "ext/T", ta) {
		t.Fatal("compatible initial initiations conflicted")
	}
	if r.imposeMethodConstraints(a, b, mt) || r.imposeMethodConstraints(c, d, mt) {
		t.Fatal("disjoint compatible groups were rejected")
	}
	if r.imposeMethodConstraints(b, c, mt) {
		t.Fatal("compatible merge was rejected")
	}
	if len(r.constraints) != 3 {
		t.Fatalf("accepted merge published %d edges, want 3", len(r.constraints))
	}
	// A later initiation that disagrees with the merged component is refused.
	if !r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("conflicting late initiation was accepted by the merged component")
	}
	keyB := initiationKey{Loader: b, Name: "ext/T"}
	if _, ok := r.initiations[keyB]; ok {
		t.Fatal("conflicting late initiation was published")
	}
	// The agreeing class completes the pending endpoints.
	if r.recordInitiation(b, "ext/T", ta) || r.recordInitiation(c, "ext/T", ta) {
		t.Fatal("agreeing late initiations conflicted")
	}
}

// TestLoadingConstraintStageRollsBackWholeDescriptor proves that a conflict on a
// later descriptor type leaves the edge staged for an earlier type unpublished.
func TestLoadingConstraintStageRollsBackWholeDescriptor(t *testing.T) {
	r, a, b, ta, tb := constraintRuntime(t)
	if r.recordInitiation(a, "ext/T", ta) || r.recordInitiation(b, "ext/T", tb) {
		t.Fatal("initial initiations conflicted")
	}
	first, err := ParseFieldType("Lext/First;")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ParseFieldType("Lext/T;")
	if err != nil {
		t.Fatal(err)
	}
	mt := &MethodType{Params: []Type{first, second}, Return: Type{Kind: 'V', Descriptor: "V"}}
	if !r.imposeMethodConstraints(a, b, mt) {
		t.Fatal("a conflict on the second descriptor type was accepted")
	}
	if len(r.constraints) != 0 {
		t.Fatalf("refused descriptor operation published %d edges", len(r.constraints))
	}
	if len(r.initiations) != 2 {
		t.Fatalf("refused descriptor operation changed initiations: %d", len(r.initiations))
	}
}

// prepFixture builds two manually linked loaders with a superclass declaration
// in l1 and an interface declaration in l2 whose leaf type ext/N resolves to
// distinct classes in each loader. The class chain cannot arise through the
// public API because parent-first visibility would hand both loaders the same
// class; the runtime unit test exercises the constraint transition directly.
type prepFixture struct {
	r     *Runtime
	l1    *Loader
	l2    *Loader
	base  *Class
	iface *Class
	x1    *Class
	x2    *Class
}

func newPrepFixture(t *testing.T) *prepFixture {
	t.Helper()
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	l2 := &Loader{ID: 2, Name: "prep-root", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	l1 := &Loader{ID: 3, Name: "prep-base", Runtime: r, Parent: l2, Classes: map[string]*Class{}}
	r.nextLoaderID = 3
	base := testHiddenClass(l1, "ext/Base", AccPublic, object)
	iface := testHiddenClass(l2, "ext/IFace", AccPublic|AccInterface|AccAbstract, object)
	testHiddenMethod(base, "p", "(Lext/Pending;)V", AccPublic)
	testHiddenMethod(base, "m", "(Lext/N;)V", AccPublic)
	testHiddenMethod(iface, "m", "(Lext/N;)V", AccPublic|AccAbstract)
	x1 := testHiddenClass(l1, "ext/N", AccPublic, object)
	x2 := testHiddenClass(l2, "ext/N", AccPublic, object)
	if r.recordInitiation(l1, "ext/N", x1) || r.recordInitiation(l2, "ext/N", x2) {
		t.Fatal("initial ext/N initiations conflicted")
	}
	return &prepFixture{r: r, l1: l1, l2: l2, base: base, iface: iface, x1: x1, x2: x2}
}

func pendingOverrideDefs() []*ClassDef {
	return []*ClassDef{{
		Name: "ext/PendingChild", Super: "ext/Base",
		Methods: []*MethodDef{testPluginVoidDef("ext/PendingChild", "p", "(Lext/Pending;)V", 1)},
	}}
}

func conflictingOverrideDef() *ClassDef {
	return &ClassDef{
		Name: "ext/ConflictChild", Super: "ext/Base", Interfaces: []string{"ext/IFace"},
		Methods: []*MethodDef{testPluginVoidDef("ext/ConflictChild", "m", "(Lext/N;)V", 1)},
	}
}

// TestPluginPreparationStagesPendingEqualities exercises the real preparation
// path. A pending reference equality and a pending reference-array equality are
// published with the accepted loader, enforced on a later initiation, and a
// primitive descriptor adds no edge.
func TestPluginPreparationStagesPendingEqualities(t *testing.T) {
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	l2 := &Loader{ID: 2, Name: "prep-root", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	l1 := &Loader{ID: 3, Name: "prep-base", Runtime: r, Parent: l2, Classes: map[string]*Class{}}
	r.nextLoaderID = 3
	base := testHiddenClass(l1, "ext/Base", AccPublic, object)
	testHiddenMethod(base, "m", "(Lext/Sig;)V", AccPublic)
	testHiddenMethod(base, "arr", "([[Lext/Arr;)V", AccPublic)
	testHiddenMethod(base, "prim", "(I)V", AccPublic)

	def := &ClassDef{
		Name: "ext/Child", Super: "ext/Base",
		Methods: []*MethodDef{
			testPluginVoidDef("ext/Child", "m", "(Lext/Sig;)V", 1),
			testPluginVoidDef("ext/Child", "arr", "([[Lext/Arr;)V", 1),
			testPluginVoidDef("ext/Child", "prim", "(I)V", 1),
		},
	}
	loader, thrown := r.installPlugin(l1, &AdmissionResult{
		Classes: []*ClassDef{def}, Defined: []string{"ext/Child"},
	})
	if thrown != nil {
		t.Fatalf("preparation rejected a pending equality: %v", thrown)
	}
	sigKey := makeConstraintKey(initiationKey{Loader: loader, Name: "ext/Sig"}, initiationKey{Loader: l1, Name: "ext/Sig"})
	arrKey := makeConstraintKey(initiationKey{Loader: loader, Name: "ext/Arr"}, initiationKey{Loader: l1, Name: "ext/Arr"})
	if _, ok := r.constraints[sigKey]; !ok {
		t.Fatal("preparation did not publish the pending reference equality")
	}
	if _, ok := r.constraints[arrKey]; !ok {
		t.Fatal("preparation did not publish the pending reference-array equality")
	}
	if len(r.constraints) != 2 {
		t.Fatalf("preparation published %d edges, want 2 (primitive descriptors must add none)", len(r.constraints))
	}
	if len(r.initiations) != 0 {
		t.Fatalf("preparation forced descriptor leaf lookups: %d initiations", len(r.initiations))
	}

	// A later initiation resolves one side; a different class on the other side
	// must then be refused without being published.
	childSig := &Class{Name: "ext/Sig", Loader: loader, Super: object, MethodIndex: map[string]*Method{}}
	parentSig := &Class{Name: "ext/Sig", Loader: l1, Super: object, MethodIndex: map[string]*Method{}}
	if r.recordInitiation(loader, "ext/Sig", childSig) {
		t.Fatal("first pending-side initiation conflicted")
	}
	if !r.recordInitiation(l1, "ext/Sig", parentSig) {
		t.Fatal("pending reference equality was not enforced on a later initiation")
	}
	if _, ok := r.initiations[initiationKey{Loader: l1, Name: "ext/Sig"}]; ok {
		t.Fatal("conflicting later initiation was published")
	}

	childArr := &Class{Name: "ext/Arr", Loader: loader, Super: object, MethodIndex: map[string]*Method{}}
	parentArr := &Class{Name: "ext/Arr", Loader: l1, Super: object, MethodIndex: map[string]*Method{}}
	if r.recordInitiation(loader, "ext/Arr", childArr) {
		t.Fatal("first pending array-leaf initiation conflicted")
	}
	if !r.recordInitiation(l1, "ext/Arr", parentArr) {
		t.Fatal("pending array-leaf equality was not enforced on a later initiation")
	}
	if _, ok := r.initiations[initiationKey{Loader: l1, Name: "ext/Arr"}]; ok {
		t.Fatal("conflicting later array-leaf initiation was published")
	}
}

// TestPluginPreparationPendingEqualityDiscardedOnRejection proves that a
// pending equality staged for an earlier class disappears when a later class in
// the same candidate conflicts and rejects the whole admission.
func TestPluginPreparationPendingEqualityDiscardedOnRejection(t *testing.T) {
	f := newPrepFixture(t)
	beforeInit := len(f.r.initiations)
	beforeEdges := len(f.r.constraints)

	defs := append(pendingOverrideDefs(), conflictingOverrideDef())
	_, thrown := f.r.installPlugin(f.l1, &AdmissionResult{
		Classes: defs, Defined: []string{"ext/PendingChild", "ext/ConflictChild"},
	})
	if thrown == nil {
		t.Fatal("conflicting preparation was accepted")
	}
	if data := ThrowableOf(thrown); data == nil || data.Code != "VERIFY" {
		t.Fatalf("rejected preparation did not report VERIFY: %#v", data)
	}
	if len(f.r.constraints) != beforeEdges {
		t.Fatalf("rejected candidate published %d pending edges, want %d",
			len(f.r.constraints), beforeEdges)
	}
	if len(f.r.initiations) != beforeInit {
		t.Fatalf("rejected candidate changed initiations: %d -> %d",
			beforeInit, len(f.r.initiations))
	}
	if len(f.r.pluginLoaders) != 0 || f.r.nextLoaderID != 3 {
		t.Fatal("rejected candidate was published")
	}

	// Control: the first class alone does stage a pending equality through the
	// same preparation path, so the rejected combined candidate discarded it.
	loader, thrown := f.r.installPlugin(f.l1, &AdmissionResult{
		Classes: pendingOverrideDefs(), Defined: []string{"ext/PendingChild"},
	})
	if thrown != nil {
		t.Fatalf("control admission failed: %v", thrown)
	}
	key := makeConstraintKey(initiationKey{Loader: loader, Name: "ext/Pending"}, initiationKey{Loader: f.l1, Name: "ext/Pending"})
	if _, ok := f.r.constraints[key]; !ok {
		t.Fatal("control admission did not publish the pending equality")
	}
	if len(f.r.constraints) != 1 {
		t.Fatalf("control admission published %d edges, want 1", len(f.r.constraints))
	}
}
