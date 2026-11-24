package rt

import "testing"

func dispatchClass(r *Runtime, name string, loader *Loader, super *Class, flags uint16, interfaces ...*Class) *Class {
	return &Class{
		Name: name, Loader: loader, Super: super, Interfaces: interfaces,
		Flags: flags, MethodIndex: map[string]*Method{},
	}
}

func dispatchMethod(owner *Class, name, descriptor string, flags uint16) *Method {
	m := &Method{Owner: owner, Name: name, Descriptor: descriptor, Flags: flags}
	if err := owner.addMethod(m); err != nil {
		panic(err)
	}
	return m
}

func TestSelectVirtualInheritedInterfaceImplementation(t *testing.T) {
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	base := dispatchClass(r, "app/Base", r.Loader, object, AccPublic)
	baseValue := dispatchMethod(base, "value", "()I", AccPublic)
	step := dispatchClass(r, "app/Step", r.Loader, object, AccInterface|AccAbstract|AccPublic)
	stepValue := dispatchMethod(step, "value", "()I", AccPublic|AccAbstract)
	child := dispatchClass(r, "ext/Child", &Loader{ID: 2, Runtime: r, Parent: r.Loader},
		base, AccPublic, step)

	selected := r.selectVirtual(child, stepValue)
	if selected == nil {
		t.Fatal("inherited interface implementation was not selected")
	}
	if selected != baseValue {
		t.Fatalf("selected %s.%s, want inherited app/Base.value", selected.Owner.Name, selected.Name)
	}

	// Once the receiver class declares its own implementation it wins.
	childValue := dispatchMethod(child, "value", "()I", AccPublic)
	selected = r.selectVirtual(child, stepValue)
	if selected != childValue {
		t.Fatalf("declared override not selected: target=%v", selected)
	}
}

func TestSelectVirtualPackagePrivateOverrideRules(t *testing.T) {
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	pkgBase := dispatchClass(r, "demo/Base", r.Loader, object, AccPublic)
	pkgValue := dispatchMethod(pkgBase, "value", "()I", 0)
	childLoader := &Loader{ID: 2, Runtime: r, Parent: r.Loader}
	child := dispatchClass(r, "demo/Child", childLoader, pkgBase, AccPublic)
	dispatchMethod(child, "value", "()I", AccPublic)

	selected := r.selectVirtual(child, pkgValue)
	if selected != pkgValue {
		t.Fatalf("package-private declaration was overridden across defining loaders: %s", selected.Owner.Name)
	}

	// A public declaration is legitimately overridden.
	publicBase := dispatchClass(r, "demo/PublicBase", r.Loader, object, AccPublic)
	publicValue := dispatchMethod(publicBase, "value", "()I", AccPublic)
	publicChild := dispatchClass(r, "demo/PublicChild", childLoader, publicBase, AccPublic)
	childOverride := dispatchMethod(publicChild, "value", "()I", AccPublic)
	selected = r.selectVirtual(publicChild, publicValue)
	if selected != childOverride {
		t.Fatalf("public override not selected: target=%v", selected)
	}
}

func TestSelectVirtualMissingInterfaceImplementation(t *testing.T) {
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	step := dispatchClass(r, "app/Step", r.Loader, object, AccInterface|AccAbstract|AccPublic)
	stepValue := dispatchMethod(step, "value", "()I", AccPublic|AccAbstract)
	orphan := dispatchClass(r, "ext/Orphan", &Loader{ID: 2, Runtime: r, Parent: r.Loader}, object, AccPublic, step)
	if selected := r.selectVirtual(orphan, stepValue); selected != nil {
		t.Fatalf("missing interface implementation selected %s", selected.Owner.Name)
	}
}

func interfaceCallRuntime(t *testing.T, baseFlags uint16) (*Runtime, *Thread, *Method, *Object) {
	t.Helper()
	r := gate3Runtime(t)
	object := r.mustClass("java/lang/Object")
	base := dispatchClass(r, "api/Base", r.Loader, object, AccPublic)
	baseValue := dispatchMethod(base, "value", "()I", baseFlags)
	baseValue.Engine = "aot"
	baseValue.MaxStack = 1
	baseValue.Entry = func(_ *Thread, _ *Method, _ *Object, _ []Value) Outcome {
		return OkValue(IntValue(31))
	}
	step := dispatchClass(r, "api/Step", r.Loader, object, AccInterface|AccAbstract|AccPublic)
	stepValue := dispatchMethod(step, "value", "()I", AccPublic|AccAbstract)
	r.Classes["api/Base"] = base
	r.Loader.Classes["api/Base"] = base
	r.Classes["api/Step"] = step
	r.Loader.Classes["api/Step"] = step
	plugin := &Loader{ID: 2, Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	child := dispatchClass(r, "ext/Child", plugin, base, AccPublic, step)
	plugin.Classes[child.Name] = child
	callerClass := dispatchClass(r, "app/Main", r.Loader, object, AccPublic)
	caller := dispatchMethod(callerClass, "call", "(Lapi/Step;)I", AccPublic|AccStatic)
	caller.Engine = "aot"
	_ = stepValue
	receiver := r.NewInstance(child)
	return r, r.MainThread(), caller, receiver
}

func TestInvokeInterfaceInheritedPublicImplementation(t *testing.T) {
	r, th, caller, receiver := interfaceCallRuntime(t, AccPublic)
	out := th.InvokeInterface(caller, 1, 7, RefValue(receiver), "api/Step", "value", "()I", nil)
	if out.Thrown != nil {
		t.Fatalf("public inherited implementation threw: %v", out.Thrown.Class)
	}
	if AsInt(out.Value) != 31 {
		t.Fatalf("unexpected result %d", AsInt(out.Value))
	}
	_ = r
}

func TestInvokeInterfaceRejectsPackagePrivateSelectedMethod(t *testing.T) {
	_, th, caller, receiver := interfaceCallRuntime(t, 0)
	out := th.InvokeInterface(caller, 1, 7, RefValue(receiver), "api/Step", "value", "()I", nil)
	if out.Thrown == nil || out.Thrown.Class == nil || out.Thrown.Class.Name != "java/lang/IllegalAccessError" {
		t.Fatalf("package-private selected method did not raise IllegalAccessError: %v", out.Thrown)
	}
}
