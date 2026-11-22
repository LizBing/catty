package rt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

// PluginError is a typed admission diagnostic crossing the runtime/admission
// boundary without importing the admission package. The adapter returns it and
// the runtime converts it into a Java PluginAdmissionException.
type PluginError struct {
	Code       string
	Message    string
	Entry      string
	Class      string
	Method     string
	Descriptor string
	PC         int
	Feature    string
	Limit      string
}

func (e *PluginError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "plugin admission failed: " + e.Code
}

// defaultAdmission is the process-wide adapter installed by a generated Gate 3
// application before the runtime is created.
var defaultAdmission func(*AdmissionRequest) (*AdmissionResult, error)

// SetAdmission installs the profile admission callback for runtimes created
// after this call. Generated Gate 3 applications register the narrow
// catty/plugin adapter before Install.
func SetAdmission(fn func(*AdmissionRequest) (*AdmissionResult, error)) {
	defaultAdmission = fn
}

// SetAdmission installs the callback on one runtime.
func (r *Runtime) SetAdmission(fn func(*AdmissionRequest) (*AdmissionResult, error)) {
	r.admit = fn
}

func pluginCode(err error) string {
	if pe, ok := err.(*PluginError); ok && pe.Code != "" {
		return pe.Code
	}
	return "VERIFY"
}

func (r *Runtime) pluginLimits() *PluginLimits {
	if r.Image != nil && r.Image.Limits != nil {
		return r.Image.Limits
	}
	if r.profile != nil && r.profile.Limits != nil {
		return r.profile.Limits
	}
	return Gate3PluginLimits
}

// loaderHandle returns the canonical ClassLoader mirror for l. The bootstrap
// loader is represented by Java null.
func (r *Runtime) loaderHandle(l *Loader) *Object {
	if l == nil || l == r.BootstrapLoader {
		return nil
	}
	if l.Mirror != nil {
		return l.Mirror
	}
	class := r.mustClass("java/lang/ClassLoader")
	l.Mirror = &Object{Class: class, Payload: l, Identity: r.newIdentity()}
	return l.Mirror
}

// NewPluginAdmissionException creates the Java throwable carrying a stable
// admission code.
func (r *Runtime) NewPluginAdmissionException(t *Thread, code, message string) *Object {
	o := r.NewThrowable(t, "catty/runtime/PluginAdmissionException")
	if o != nil {
		data := ThrowableOf(o)
		data.Code = code
		if message != "" {
			data.Message = message
		}
	}
	return o
}

// pluginCaller returns the first logical Java frame below intrinsic and
// reflection frames. Reflection and ordinary calls therefore enforce the same
// plugin-defined-caller rule.
func (t *Thread) pluginCaller() *Method {
	for i := len(t.frames) - 1; i >= 0; i-- {
		m := t.frames[i].Method
		if m == nil || m.Owner == nil {
			continue
		}
		if m.Engine == "intrinsic" && m.Owner.Name != "java/lang/reflect/Method" {
			continue
		}
		if m.Owner.Name == "java/lang/reflect/Method" {
			continue
		}
		return m
	}
	return nil
}

func (t *Thread) callerIsPlugin() bool {
	m := t.pluginCaller()
	if m == nil || m.Owner == nil || m.Owner.Loader == nil {
		return false
	}
	return m.Owner.Loader.Plugin
}

// intrinsicPluginsReadJar implements Plugins.readJar(String)byte[].
func intrinsicPluginsReadJar(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 || args[0].Kind != KindRef || args[0].Ref == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	path, ok := StringText(args[0].Ref)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	if t.callerIsPlugin() {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/SecurityException"))
	}
	limits := t.RT.pluginLimits()
	// Open without blocking so a FIFO or device cannot stall the JVM, then
	// verify the opened handle before reading. The actual byte read is bounded
	// to the ceiling plus one sentinel byte.
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/io/IOException"))
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/io/IOException"))
	}
	limit := limits.CompressedBytes
	if limit <= 0 {
		limit = Gate3PluginLimits.CompressedBytes
	}
	buf := make([]byte, 0, 64*1024)
	tmp := make([]byte, 64*1024)
	for int64(len(buf)) <= limit {
		want := int64(len(tmp))
		remaining := limit + 1 - int64(len(buf))
		if remaining < want {
			want = remaining
		}
		n, rerr := f.Read(tmp[:want])
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/io/IOException"))
		}
		if n == 0 {
			break
		}
	}
	if int64(len(buf)) > limit {
		return ThrowOutcome(t.RT.NewPluginAdmissionException(t, "LIMIT",
			fmt.Sprintf("compressed input exceeds %d bytes", limit)))
	}
	arrayClass := t.RT.arrayClassForComponent(t.RT.PrimitiveClass('B'))
	arr := t.RT.NewArrayOfClass(arrayClass, len(buf))
	for i, b := range buf {
		arr.Array.Bytes[i] = int8(b)
	}
	return OkValue(RefValue(arr))
}

// intrinsicPluginsRegister implements
// Plugins.register(byte[],ClassLoader)ClassLoader.
func intrinsicPluginsRegister(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 2 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[0].Kind != KindRef || args[0].Ref == nil || args[0].Ref.Array == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	limit := t.RT.pluginLimits().CompressedBytes
	if limit > 0 && int64(args[0].Ref.Array.Len()) > limit {
		return ThrowOutcome(t.RT.NewPluginAdmissionException(t, "LIMIT",
			fmt.Sprintf("compressed input exceeds %d bytes", limit)))
	}
	parent, thrown := t.pluginParent(args[1])
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if t.RT.admit == nil {
		return ThrowOutcome(t.RT.NewPluginAdmissionException(t, "PROFILE",
			"plugin admission is not enabled"))
	}
	// The candidate owns a private copy of the caller's bytes.
	jar := make([]byte, args[0].Ref.Array.Len())
	for i := range jar {
		jar[i] = byte(args[0].Ref.Array.Bytes[i])
	}
	digest := sha256.Sum256(jar)
	res, err := t.RT.admit(&AdmissionRequest{Bytes: jar, Parent: parent, Profile: t.RT.profile.ID})
	if err != nil {
		pe, _ := err.(*PluginError)
		code := pluginCode(err)
		msg := err.Error()
		if tr := t.trace(); tr != nil {
			tr.admissionRejected(t.ID, pe, digest[:], parent)
		}
		return ThrowOutcome(t.RT.NewPluginAdmissionException(t, code, msg))
	}
	if res == nil {
		return ThrowOutcome(t.RT.NewPluginAdmissionException(t, "VERIFY", "admission produced no definitions"))
	}
	if res.Digest == "" {
		res.Digest = hex.EncodeToString(digest[:])
	}
	loader, thrown := t.RT.installPlugin(parent, res)
	if thrown != nil {
		if tr := t.trace(); tr != nil {
			code := "VERIFY"
			message := ""
			if data := ThrowableOf(thrown); data != nil {
				if data.Code != "" {
					code = data.Code
				}
				message = data.Message
			}
			tr.admissionRejected(t.ID, &PluginError{Code: code, Message: message}, digest[:], parent)
		}
		return ThrowOutcome(thrown)
	}
	if tr := t.trace(); tr != nil {
		tr.admission(t.ID, res, digest[:], parent, loader)
		for _, name := range res.Defined {
			if c, ok := loader.Local(name); ok {
				tr.classDefine(t.ID, c, loader)
			}
		}
	}
	return OkValue(RefValue(t.RT.loaderHandle(loader)))
}

func (t *Thread) pluginParent(v Value) (*Loader, *Object) {
	if v.Kind != KindRef {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
	}
	if v.Ref == nil {
		return t.RT.BootstrapLoader, nil
	}
	l, ok := v.Ref.Payload.(*Loader)
	if !ok || l == nil || l.Runtime != t.RT {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
	}
	if !t.RT.loaderPublished(l) {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalArgumentException")
	}
	return l, nil
}

func (r *Runtime) loaderPublished(l *Loader) bool {
	if l == r.BootstrapLoader || l == r.Loader {
		return true
	}
	for _, p := range r.pluginLoaders {
		if p == l {
			return true
		}
	}
	return false
}

// installPlugin commits one validated candidate. All definitions are built in
// a detached loader first; nothing is published until every step succeeds. A
// failed commit restores the array intern caches and identity cursor so the
// runtime transaction leaves the pre-registration snapshot intact.
func (r *Runtime) installPlugin(parent *Loader, res *AdmissionResult) (*Loader, *Object) {
	arrays := make(map[*Class]*Class, len(r.arrayClasses))
	for k, v := range r.arrayClasses {
		arrays[k] = v
	}
	primitiveArrays := make(map[byte]*Class, len(r.primitiveArrayClass))
	for k, v := range r.primitiveArrayClass {
		primitiveArrays[k] = v
	}
	identity := r.nextIdentity
	loader, thrown := r.installPluginTx(parent, res)
	if thrown != nil {
		r.arrayClasses = arrays
		r.primitiveArrayClass = primitiveArrays
		r.nextIdentity = identity
	}
	return loader, thrown
}

func (r *Runtime) installPluginTx(parent *Loader, res *AdmissionResult) (*Loader, *Object) {
	limits := r.pluginLimits()
	if limits.Loaders > 0 && len(r.pluginLoaders)+1 > limits.Loaders {
		return nil, r.NewPluginAdmissionException(nil, "LIMIT",
			fmt.Sprintf("published plugin loader limit %d reached", limits.Loaders))
	}
	if limits.CumulativeClassBytes > 0 && r.pluginClassBytes+res.ClassBytes > limits.CumulativeClassBytes {
		return nil, r.NewPluginAdmissionException(nil, "LIMIT",
			fmt.Sprintf("cumulative admitted class bytes limit %d reached", limits.CumulativeClassBytes))
	}
	loader := &Loader{
		ID:      r.nextLoaderID + 1,
		Name:    fmt.Sprintf("plugin-%d", len(r.pluginLoaders)+1),
		Runtime: r,
		Parent:  parent,
		Classes: map[string]*Class{},
		Plugin:  true,
		Digest:  res.Digest,
	}
	defs := map[string]*ClassDef{}
	created := make([]*Class, 0, len(res.Classes))
	seen := map[string]bool{}
	for _, def := range res.Classes {
		if def == nil || def.Name == "" {
			return nil, r.NewPluginAdmissionException(nil, "CLASS_FORMAT", "plugin contains an unnamed class")
		}
		if seen[def.Name] {
			return nil, r.NewPluginAdmissionException(nil, "CLASS_FORMAT", "duplicate plugin class definition "+def.Name)
		}
		seen[def.Name] = true
		if _, ok := loader.Parent.Lookup(def.Name); ok {
			return nil, r.NewPluginAdmissionException(nil, "CLASS_FORMAT", "effective plugin class shadows a parent definition: "+def.Name)
		}
		if reservedImageName(r.profile, def.Name) {
			return nil, r.NewPluginAdmissionException(nil, "PACKAGE", "reserved class name "+def.Name)
		}
		c := &Class{
			Name:        def.Name,
			Flags:       def.Flags,
			MethodIndex: map[string]*Method{},
			NestHost:    def.NestHost,
			NestMembers: append([]string(nil), def.NestMembers...),
			Loader:      loader,
		}
		loader.Classes[def.Name] = c
		defs[def.Name] = def
		created = append(created, c)
	}
	// Link every definition before any method is materialized.
	for _, c := range created {
		def := defs[c.Name]
		if def.Super != "" {
			s, ok := loader.Lookup(def.Super)
			if !ok {
				return nil, r.NewPluginAdmissionException(nil, "DEPENDENCY", "missing superclass "+def.Super)
			}
			c.Super = s
		}
		for _, in := range def.Interfaces {
			i, ok := loader.Lookup(in)
			if !ok {
				return nil, r.NewPluginAdmissionException(nil, "DEPENDENCY", "missing interface "+in)
			}
			c.Interfaces = append(c.Interfaces, i)
		}
	}
	// Hierarchy cycles and final-superclass constraints.
	if err := r.checkPluginHierarchies(created); err != nil {
		return nil, r.NewPluginAdmissionException(nil, "VERIFY", err.Error())
	}
	for _, c := range created {
		if err := r.layoutDefinition(c, defs[c.Name]); err != nil {
			return nil, r.NewPluginAdmissionException(nil, "VERIFY", err.Error())
		}
	}
	for _, c := range created {
		def := defs[c.Name]
		for _, md := range def.Methods {
			if md.Engine != "interp" && md.Engine != "abstract" {
				return nil, r.NewPluginAdmissionException(nil, "PROFILE", "plugin method is not interpreted")
			}
			if err := r.installMethodDef(c, md); err != nil {
				return nil, r.NewPluginAdmissionException(nil, "VERIFY", err.Error())
			}
		}
	}
	for _, c := range created {
		if err := r.resolveClassMetadata(c, false); err != nil {
			return nil, r.NewPluginAdmissionException(nil, "DEPENDENCY", err.Error())
		}
	}
	if err := r.checkOverrideConstraints(created); err != nil {
		return nil, r.NewPluginAdmissionException(nil, "VERIFY", err.Error())
	}
	for _, c := range created {
		c.Mirror = &Object{Class: r.mustClass("java/lang/Class"), Payload: c, Identity: r.newIdentity()}
		for i := range c.StaticField {
			c.Statics[i] = zeroValueFor(c.StaticField[i].Descriptor)
		}
	}
	r.nextLoaderID = loader.ID
	r.pluginLoaders = append(r.pluginLoaders, loader)
	r.pluginClassBytes += res.ClassBytes
	return loader, nil
}

// checkOverrideConstraints rejects an override whose reference descriptor
// leaf types resolve to different classes through the two defining loaders.
// Primitive descriptors and same-loader declarations impose no constraint.
// The candidate is compared with every inherited class and interface
// declaration with the same name and descriptor, including default methods
// inherited through the interface closure.
func (r *Runtime) checkOverrideConstraints(created []*Class) error {
	for _, c := range created {
		for _, m := range c.Methods {
			if m.Flags&(AccStatic|AccPrivate) != 0 || m.Name == "<init>" || m.Name == "<clinit>" {
				continue
			}
			for _, p := range inheritedMethods(c, m.Name+m.Descriptor) {
				if p.Owner == nil || p.Owner.Loader == nil || p.Owner.Loader == c.Loader {
					continue
				}
				if !r.methodOverrides(m, p, map[*Method]bool{}) {
					continue
				}
				if len(p.Type.Params) != len(m.Type.Params) ||
					p.Type.Return.Descriptor != m.Type.Return.Descriptor {
					continue
				}
				if ca, cb, name, conflict := r.leafTypeConflict(c.Loader, p.Owner.Loader, m.Type.Params, p.Type.Params); conflict {
					return fmt.Errorf("loader constraint conflict: %s.%s%s parameter leaf %s resolves to %s and %s",
						c.Name, m.Name, m.Descriptor, name, ca.Name, cb.Name)
				}
				if m.Type.Return.IsReference() {
					if ca, cb, name, conflict := r.leafTypeConflict(c.Loader, p.Owner.Loader,
						[]Type{m.Type.Return}, []Type{p.Type.Return}); conflict {
						return fmt.Errorf("loader constraint conflict: %s.%s%s return leaf %s resolves to %s and %s",
							c.Name, m.Name, m.Descriptor, name, ca.Name, cb.Name)
					}
				}
			}
		}
	}
	return nil
}

// inheritedMethods returns every distinct same-signature method declaration in
// the superclass chain and interface closure of c.
func inheritedMethods(c *Class, key string) []*Method {
	var out []*Method
	seenMethod := map[*Method]bool{}
	seenClass := map[*Class]bool{}
	add := func(p *Method) {
		if p != nil && !seenMethod[p] {
			seenMethod[p] = true
			out = append(out, p)
		}
	}
	var visitInterface func(in *Class)
	visitInterface = func(in *Class) {
		if in == nil || seenClass[in] {
			return
		}
		seenClass[in] = true
		add(in.MethodIndex[key])
		for _, parent := range in.Interfaces {
			visitInterface(parent)
		}
	}
	for sup := c.Super; sup != nil; sup = sup.Super {
		add(sup.MethodIndex[key])
		for _, in := range sup.Interfaces {
			visitInterface(in)
		}
	}
	for _, in := range c.Interfaces {
		visitInterface(in)
	}
	return out
}

// leafTypeConflict resolves corresponding reference leaf types through two
// loaders. It returns the two classes, the leaf name and whether a conflict
// exists; unresolved types keep ordinary resolution timing and are skipped.
func (r *Runtime) leafTypeConflict(aLoader, bLoader *Loader, aTypes, bTypes []Type) (*Class, *Class, string, bool) {
	for i := range aTypes {
		name := leafClassName(aTypes[i])
		if name == "" || i >= len(bTypes) || leafClassName(bTypes[i]) != name {
			continue
		}
		a, errA := r.classForLoaderName(aLoader, name)
		if errA != nil {
			continue
		}
		b, errB := r.classForLoaderName(bLoader, name)
		if errB != nil {
			continue
		}
		if a != b {
			return a, b, name, true
		}
	}
	return nil, nil, "", false
}

func (r *Runtime) checkPluginHierarchies(created []*Class) error {
	state := map[*Class]uint8{}
	var visit func(c *Class, depth int) error
	visit = func(c *Class, depth int) error {
		if depth > 128 {
			return fmt.Errorf("hierarchy traversal depth exceeds 128 at %s", c.Name)
		}
		switch state[c] {
		case 1:
			return fmt.Errorf("class hierarchy cycle at %s", c.Name)
		case 2:
			return nil
		}
		state[c] = 1
		if c.Super != nil {
			if c.Super.IsPrimitive || c.Super.IsArray {
				return fmt.Errorf("class %s has invalid superclass %s", c.Name, c.Super.Name)
			}
			if c.Super.Flags&AccFinal != 0 {
				return fmt.Errorf("class %s extends final class %s", c.Name, c.Super.Name)
			}
			if err := visit(c.Super, depth+1); err != nil {
				return err
			}
		}
		for _, in := range c.Interfaces {
			if !in.IsInterface() {
				return fmt.Errorf("class %s implements non-interface %s", c.Name, in.Name)
			}
			if err := visit(in, depth+1); err != nil {
				return err
			}
		}
		state[c] = 2
		return nil
	}
	for _, c := range created {
		if err := visit(c, 0); err != nil {
			return err
		}
	}
	return nil
}

// intrinsicPluginAdmissionGetCode implements
// PluginAdmissionException.getCode()String.
func intrinsicPluginAdmissionGetCode(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	code := ""
	if data := ThrowableOf(receiver); data != nil {
		code = data.Code
	}
	return OkValue(RefValue(t.RT.NewStringGo(code)))
}

// intrinsicClassGetClassLoader implements Class.getClassLoader()ClassLoader.
func intrinsicClassGetClassLoader(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	c, ok := receiver.Payload.(*Class)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
	}
	loader := c.Loader
	if c.IsArray && c.Component != nil {
		loader = c.Component.Loader
	}
	return OkValue(RefValue(t.RT.loaderHandle(loader)))
}

// intrinsicClassForName implements
// Class.forName(String,boolean,ClassLoader)Class.
func intrinsicClassForName(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 3 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if args[0].Kind != KindRef || args[0].Ref == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	name, ok := StringText(args[0].Ref)
	if !ok {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassNotFoundException"))
	}
	// Class.forName takes a binary name; arrays use the binary form
	// [Lpkg.Type; whose component may be dotted.
	name = strings.ReplaceAll(name, ".", "/")
	initialize := AsInt(args[1]) != 0
	loader := t.RT.BootstrapLoader
	if args[2].Kind == KindRef && args[2].Ref != nil {
		l, ok := args[2].Ref.Payload.(*Loader)
		if !ok || l == nil || l.Runtime != t.RT || !t.RT.loaderPublished(l) {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassNotFoundException"))
		}
		loader = l
	}
	c, err := t.RT.classForLoaderName(loader, name)
	if err != nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassNotFoundException"))
	}
	if t.RT.recordClassInitiation(loader, name, c) {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/LinkageError"))
	}
	if initialize && !c.IsArray && !c.IsPrimitive {
		if thrown := t.EnsureInitialized(c); thrown != nil {
			return ThrowOutcome(thrown)
		}
	}
	return OkValue(RefValue(c.Mirror))
}

// helper used by loader code/tests
func loaderID(l *Loader) int {
	if l == nil {
		return 0
	}
	return l.ID
}
