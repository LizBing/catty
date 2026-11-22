package verify

import (
	"fmt"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// The Gate 2 verifier tracks the full primitive and reference verification
// types. Every value counts one logical operand entry, while long and double
// occupy two Java local and operand slots. Stack map frames are authoritative
// at branch, switch and handler entries; each normal edge must be assignable to
// the declared frame. The analysis is a worklist over verified states and
// produces the logical stack height at every reachable instruction for the
// generated Go emitter.

const (
	vtTop byte = iota
	vtInt
	vtLong
	vtFloat
	vtDouble
	vtRef
	vtNull
	vtUninit
	vtUninitThis
)

type vtype struct {
	kind  byte
	class string
	off   int32
}

func (v vtype) width() int {
	if v.kind == vtLong || v.kind == vtDouble {
		return 2
	}
	return 1
}

func (v vtype) isWide() bool { return v.width() == 2 }

func (v vtype) String() string {
	switch v.kind {
	case vtTop:
		return "top"
	case vtInt:
		return "int"
	case vtLong:
		return "long"
	case vtFloat:
		return "float"
	case vtDouble:
		return "double"
	case vtNull:
		return "null"
	case vtUninitThis:
		return "uninitializedThis"
	case vtUninit:
		return fmt.Sprintf("uninitialized@%d", v.off)
	case vtRef:
		return v.class
	default:
		return "?"
	}
}

func vtypeFromType(t rt.Type) (vtype, error) {
	switch t.Kind {
	case 'I':
		return vtype{kind: vtInt}, nil
	case 'J':
		return vtype{kind: vtLong}, nil
	case 'F':
		return vtype{kind: vtFloat}, nil
	case 'D':
		return vtype{kind: vtDouble}, nil
	case 'L':
		return vtype{kind: vtRef, class: t.ClassName}, nil
	case 'A':
		return vtype{kind: vtRef, class: t.Descriptor}, nil
	case 'V':
		return vtype{}, fmt.Errorf("void is not a value type")
	default:
		return vtype{}, fmt.Errorf("unsupported type %q", t.Descriptor)
	}
}

func vtypeFromFrame(v classfile.VType) (vtype, error) {
	switch v.Kind {
	case classfile.VTTop:
		return vtype{kind: vtTop}, nil
	case classfile.VTInteger:
		return vtype{kind: vtInt}, nil
	case classfile.VTLong:
		return vtype{kind: vtLong}, nil
	case classfile.VTFloat:
		return vtype{kind: vtFloat}, nil
	case classfile.VTDouble:
		return vtype{kind: vtDouble}, nil
	case classfile.VTNull:
		return vtype{kind: vtNull}, nil
	case classfile.VTUninitializedThis:
		return vtype{kind: vtUninitThis}, nil
	case classfile.VTObject:
		return vtype{kind: vtRef, class: v.Class}, nil
	case classfile.VTUninitialized:
		return vtype{kind: vtUninit, off: v.Offset}, nil
	default:
		return vtype{}, fmt.Errorf("unsupported stack map verification type %d", v.Kind)
	}
}

type frameState struct {
	locals []vtype
	stack  []vtype
}

func (s *frameState) push(v vtype) { s.stack = append(s.stack, v) }

func (s *frameState) pop() (vtype, error) {
	if len(s.stack) == 0 {
		return vtype{}, fmt.Errorf("operand stack underflow")
	}
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v, nil
}

func (s *frameState) popKind(kind byte, what string) (vtype, error) {
	v, err := s.pop()
	if err != nil {
		return vtype{}, err
	}
	if v.kind != kind {
		return vtype{}, fmt.Errorf("expected %s on the operand stack, found %s", what, v)
	}
	return v, nil
}

func (s *frameState) popRef(allowUninit bool) (vtype, error) {
	v, err := s.pop()
	if err != nil {
		return vtype{}, err
	}
	if v.kind == vtRef || v.kind == vtNull {
		return v, nil
	}
	if allowUninit && (v.kind == vtUninit || v.kind == vtUninitThis) {
		return v, nil
	}
	return vtype{}, fmt.Errorf("expected a reference on the operand stack, found %s", v)
}

func (s *frameState) slots() int {
	n := 0
	for _, v := range s.stack {
		n += v.width()
	}
	return n
}

func initialState(c *classfile.Class, m *classfile.Method, mt *rt.MethodType) (*frameState, error) {
	locals := make([]vtype, m.MaxLocals)
	for i := range locals {
		locals[i] = vtype{kind: vtTop}
	}
	slot := 0
	if m.Flags&rt.AccStatic == 0 {
		if slot >= len(locals) {
			return nil, fmt.Errorf("max_locals %d cannot hold the receiver", m.MaxLocals)
		}
		if m.Name == "<init>" {
			locals[0] = vtype{kind: vtUninitThis}
		} else {
			locals[0] = vtype{kind: vtRef, class: c.Name}
		}
		slot = 1
	}
	for i, p := range mt.Params {
		vt, err := vtypeFromType(p)
		if err != nil {
			return nil, err
		}
		at := slot + mt.ParamSlot[i]
		if at+vt.width() > len(locals) {
			return nil, fmt.Errorf("parameter %d does not fit in max_locals %d", i, m.MaxLocals)
		}
		locals[at] = vt
	}
	return &frameState{locals: locals}, nil
}

func frameStateFrom(f *classfile.Frame, maxLocals int) (*frameState, error) {
	locals := make([]vtype, maxLocals)
	for i := range locals {
		locals[i] = vtype{kind: vtTop}
	}
	slot := 0
	for _, lv := range f.Locals {
		vt, err := vtypeFromFrame(lv)
		if err != nil {
			return nil, err
		}
		if slot+vt.width() > maxLocals {
			return nil, fmt.Errorf("stack map local slots exceed max_locals %d", maxLocals)
		}
		locals[slot] = vt
		slot += vt.width()
	}
	stack := make([]vtype, 0, len(f.Stack))
	for _, sv := range f.Stack {
		vt, err := vtypeFromFrame(sv)
		if err != nil {
			return nil, err
		}
		stack = append(stack, vt)
	}
	return &frameState{locals: locals, stack: stack}, nil
}

func copyState(s *frameState) *frameState {
	return &frameState{
		locals: append([]vtype(nil), s.locals...),
		stack:  append([]vtype(nil), s.stack...),
	}
}

func statesEqual(a, b *frameState) bool {
	if len(a.locals) != len(b.locals) || len(a.stack) != len(b.stack) {
		return false
	}
	for i := range a.locals {
		if a.locals[i] != b.locals[i] {
			return false
		}
	}
	for i := range a.stack {
		if a.stack[i] != b.stack[i] {
			return false
		}
	}
	return true
}

func checkAssignableState(ctx *Context, from, to *frameState) error {
	if len(from.stack) != len(to.stack) {
		return fmt.Errorf("operand stack depth %d differs from frame depth %d", len(from.stack), len(to.stack))
	}
	for i := range to.locals {
		want := to.locals[i]
		if want.kind == vtTop {
			continue
		}
		if i >= len(from.locals) || from.locals[i].kind == vtTop {
			return fmt.Errorf("local %d is uninitialized but the stack map expects %s", i, want)
		}
		if !assignable(ctx, from.locals[i], want) {
			return fmt.Errorf("local %d is %s, not assignable to %s", i, from.locals[i], want)
		}
	}
	for i := range to.stack {
		if !assignable(ctx, from.stack[i], to.stack[i]) {
			return fmt.Errorf("operand %d is %s, not assignable to %s", i, from.stack[i], to.stack[i])
		}
	}
	return nil
}

func assignable(ctx *Context, from, to vtype) bool {
	if from == to {
		return true
	}
	switch to.kind {
	case vtTop:
		return false
	case vtInt:
		return from.kind == vtInt
	case vtLong:
		return from.kind == vtLong
	case vtFloat:
		return from.kind == vtFloat
	case vtDouble:
		return from.kind == vtDouble
	case vtNull:
		return from.kind == vtNull
	case vtUninitThis, vtUninit:
		return from.kind == to.kind && from.off == to.off
	case vtRef:
		switch from.kind {
		case vtNull:
			return true
		case vtRef:
			return classAssignable(ctx, from.class, to.class)
		}
	}
	return false
}

func classAssignable(ctx *Context, from, to string) bool {
	if from == to {
		return true
	}
	if to == "java/lang/Object" {
		return true
	}
	if strings.HasPrefix(from, "[") {
		if to == "java/lang/Cloneable" || to == "java/io/Serializable" {
			return true
		}
		if strings.HasPrefix(to, "[") {
			fc, fok := arrayComponent(from)
			tc, tok := arrayComponent(to)
			if !fok || !tok {
				return false
			}
			if isPrimitiveName(fc) || isPrimitiveName(tc) {
				return fc == tc
			}
			return classAssignable(ctx, fc, tc)
		}
		return false
	}
	if strings.HasPrefix(to, "[") {
		return false
	}
	if classIsInterface(ctx, to) {
		if classIsInterface(ctx, from) && interfaceExtends(ctx, from, to) {
			return true
		}
		return classImplements(ctx, from, to)
	}
	cur := from
	seen := map[string]bool{}
	for cur != "" && !seen[cur] {
		seen[cur] = true
		info, ok := classInfo(ctx, cur)
		if !ok {
			return false
		}
		if info.super == to {
			return true
		}
		cur = info.super
	}
	return false
}

type classInfoData struct {
	super      string
	interfaces []string
	flags      uint16
}

func classInfo(ctx *Context, name string) (classInfoData, bool) {
	if c, ok := ctx.Classes[name]; ok {
		return classInfoData{super: c.Super, interfaces: c.Interfaces, flags: c.Flags}, true
	}
	if bc, ok := rt.BootstrapClassDef(name); ok {
		return classInfoData{super: bc.Super, interfaces: bc.Interfaces, flags: bc.Flags}, true
	}
	return classInfoData{}, false
}

func classIsInterface(ctx *Context, name string) bool {
	info, ok := classInfo(ctx, name)
	return ok && info.flags&rt.AccInterface != 0
}

func classImplements(ctx *Context, from, target string) bool {
	seen := map[string]bool{}
	var visitInterface func(string) bool
	visitInterface = func(name string) bool {
		if name == "" || seen[name] {
			return false
		}
		seen[name] = true
		if name == target {
			return true
		}
		info, ok := classInfo(ctx, name)
		if !ok {
			return false
		}
		for _, in := range info.interfaces {
			if in == target || visitInterface(in) {
				return true
			}
		}
		return false
	}
	for name := from; name != ""; {
		info, ok := classInfo(ctx, name)
		if !ok {
			return false
		}
		for _, in := range info.interfaces {
			if visitInterface(in) {
				return true
			}
		}
		name = info.super
	}
	return false
}

func interfaceExtends(ctx *Context, from, target string) bool {
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(name string) bool {
		if name == "" || seen[name] {
			return false
		}
		seen[name] = true
		if name == target {
			return true
		}
		info, ok := classInfo(ctx, name)
		if !ok {
			return false
		}
		for _, in := range info.interfaces {
			if visit(in) {
				return true
			}
		}
		return false
	}
	return visit(from)
}

func isPrimitiveName(name string) bool {
	if len(name) != 1 {
		return false
	}
	switch name[0] {
	case 'Z', 'B', 'C', 'S', 'I', 'J', 'F', 'D':
		return true
	}
	return false
}

// arrayComponent splits an array descriptor into its component: a primitive
// descriptor letter, a nested array descriptor, or a class internal name.
func arrayComponent(desc string) (string, bool) {
	if !strings.HasPrefix(desc, "[") || len(desc) < 2 {
		return "", false
	}
	component := desc[1:]
	switch component[0] {
	case '[':
		return component, true
	case 'L':
		if !strings.HasSuffix(component, ";") || len(component) < 3 {
			return "", false
		}
		return component[1 : len(component)-1], true
	case 'Z', 'B', 'C', 'S', 'I', 'J', 'F', 'D':
		if len(component) == 1 {
			return component, true
		}
		return "", false
	}
	return "", false
}

func arrayDescriptor(component string) string {
	if len(component) > 0 && component[0] == '[' {
		return "[" + component
	}
	if isPrimitiveName(component) {
		return "[" + component
	}
	return "[L" + component + ";"
}

func arrayComponentType(v vtype) (vtype, error) {
	if v.kind == vtNull {
		return vtype{kind: vtRef, class: "java/lang/Object"}, nil
	}
	if v.kind != vtRef || !strings.HasPrefix(v.class, "[") {
		return vtype{}, fmt.Errorf("expected an array reference, found %s", v)
	}
	component, ok := arrayComponent(v.class)
	if !ok {
		return vtype{}, fmt.Errorf("malformed array type %s", v.class)
	}
	if isPrimitiveName(component) {
		switch component {
		case "Z", "B", "C", "S", "I":
			return vtype{kind: vtInt}, nil
		case "J":
			return vtype{kind: vtLong}, nil
		case "F":
			return vtype{kind: vtFloat}, nil
		case "D":
			return vtype{kind: vtDouble}, nil
		}
	}
	return vtype{kind: vtRef, class: component}, nil
}

// typeChecker carries the per-instruction verification environment.
type typeChecker struct {
	ctx    *Context
	class  *classfile.Class
	method *classfile.Method
	mt     *rt.MethodType
	state  *frameState
	pc     int32
	insn   rt.Insn
}

func (t *typeChecker) push(v vtype) error {
	t.state.push(v)
	if t.state.slots() > t.method.MaxStack {
		return fmt.Errorf("operand stack exceeds max_stack %d", t.method.MaxStack)
	}
	return nil
}

func (t *typeChecker) pop() (vtype, error) { return t.state.pop() }

func (t *typeChecker) popOf(kind byte, what string) (vtype, error) {
	return t.state.popKind(kind, what)
}

func (t *typeChecker) popRef(allowUninit bool) (vtype, error) {
	return t.state.popRef(allowUninit)
}

func (t *typeChecker) ref(idx int32) (rt.Ref, error) {
	if idx < 0 || int(idx) >= len(t.method.Refs) {
		return rt.Ref{}, fmt.Errorf("missing constant reference")
	}
	return t.method.Refs[idx], nil
}

func (t *typeChecker) fieldType(ref rt.Ref) (vtype, error) {
	ft, err := rt.ParseFieldType(ref.Descriptor)
	if err != nil {
		return vtype{}, err
	}
	return vtypeFromType(ft)
}

func (t *typeChecker) methodType(ref rt.Ref) (*rt.MethodType, error) {
	return rt.ParseMethodType(ref.Descriptor)
}

func (t *typeChecker) classExists(name string) bool { return t.ctx.ClassExists(name) }

func (t *typeChecker) storeLocal(idx int32, v vtype) error {
	i := int(idx)
	if i < 0 || i >= len(t.state.locals) {
		return fmt.Errorf("local index %d out of range for max_locals %d", idx, len(t.state.locals))
	}
	if i+v.width() > len(t.state.locals) {
		return fmt.Errorf("local index %d does not fit a %s in max_locals %d", idx, v, len(t.state.locals))
	}
	if i > 0 && t.state.locals[i-1].isWide() {
		t.state.locals[i-1] = vtype{kind: vtTop}
	}
	t.state.locals[i] = v
	if v.width() == 2 {
		t.state.locals[i+1] = vtype{kind: vtTop}
	}
	return nil
}

func (t *typeChecker) loadLocal(idx int32, want vtype) error {
	i := int(idx)
	if i < 0 || i >= len(t.state.locals) {
		return fmt.Errorf("local index %d out of range for max_locals %d", idx, len(t.state.locals))
	}
	got := t.state.locals[i]
	if got.kind != want.kind {
		if got.kind == vtTop && i > 0 && t.state.locals[i-1].isWide() {
			return fmt.Errorf("local index %d is the unusable second slot of the %s at local %d", idx, t.state.locals[i-1], idx-1)
		}
		return fmt.Errorf("local index %d holds %s, not %s", idx, got, want)
	}
	if want.width() == 2 {
		if i+1 >= len(t.state.locals) {
			return fmt.Errorf("local index %d does not fit a %s in max_locals %d", idx, want, len(t.state.locals))
		}
		if t.state.locals[i+1].kind != vtTop {
			return fmt.Errorf("local index %d: second slot of the %s is not unusable (found %s)", idx, want, t.state.locals[i+1])
		}
	}
	return t.push(got)
}

func (t *typeChecker) popArray(wantPrimitive string) (vtype, error) {
	v, err := t.popRef(false)
	if err != nil {
		return vtype{}, err
	}
	if v.kind == vtNull {
		return v, nil
	}
	if !strings.HasPrefix(v.class, "[") {
		return vtype{}, fmt.Errorf("expected an array reference, found %s", v)
	}
	if wantPrimitive != "" {
		if v.class != "["+wantPrimitive {
			return vtype{}, fmt.Errorf("expected a %s[] array, found %s", primitiveLabel(wantPrimitive), v)
		}
	}
	return v, nil
}

func primitiveLabel(desc string) string {
	switch desc {
	case "Z":
		return "boolean"
	case "B":
		return "byte"
	case "C":
		return "char"
	case "S":
		return "short"
	case "I":
		return "int"
	case "J":
		return "long"
	case "F":
		return "float"
	case "D":
		return "double"
	}
	return desc
}

func (t *typeChecker) pushReturn(mt *rt.MethodType) error {
	if mt.Return.Kind == 'V' {
		return nil
	}
	vt, err := vtypeFromType(mt.Return)
	if err != nil {
		return err
	}
	return t.push(vt)
}

// successor is one verified outgoing edge.
type successor struct {
	pc    int32
	state *frameState
}

func insertStack(s []vtype, at int, vals ...vtype) []vtype {
	out := make([]vtype, 0, len(s)+len(vals))
	out = append(out, s[:at]...)
	out = append(out, vals...)
	out = append(out, s[at:]...)
	return out
}

func (t *typeChecker) stackManipulation() error {
	s := t.state.stack
	n := len(s)
	switch t.insn.Op {
	case rt.OpPop:
		v, err := t.pop()
		if err != nil {
			return err
		}
		if v.width() != 1 {
			return fmt.Errorf("pop of a category-two value %s", v)
		}
	case rt.OpPop2:
		if n == 0 {
			return fmt.Errorf("pop2 with an empty stack")
		}
		if s[n-1].width() == 2 {
			t.state.stack = s[:n-1]
			return nil
		}
		if n < 2 || s[n-2].width() != 1 {
			return fmt.Errorf("pop2 requires two category-one values or one category-two value")
		}
		t.state.stack = s[:n-2]
	case rt.OpDup:
		if n < 1 || s[n-1].width() != 1 {
			return fmt.Errorf("dup requires a category-one value")
		}
		t.state.stack = insertStack(s, n, s[n-1])
	case rt.OpDupX1:
		if n < 2 || s[n-1].width() != 1 || s[n-2].width() != 1 {
			return fmt.Errorf("dup_x1 requires two category-one values")
		}
		v1, v2 := s[n-1], s[n-2]
		t.state.stack = append(s[:n-2], v1, v2, v1)
	case rt.OpDupX2:
		if n < 2 {
			return fmt.Errorf("dup_x2 with an empty stack")
		}
		v1 := s[n-1]
		if v1.width() != 1 {
			return fmt.Errorf("dup_x2 requires a category-one top value")
		}
		if s[n-2].width() == 2 {
			v2 := s[n-2]
			t.state.stack = append(s[:n-2], v1, v2, v1)
			return nil
		}
		if n < 3 || s[n-3].width() != 1 {
			return fmt.Errorf("dup_x2 requires three category-one values or one category-two and one category-one value")
		}
		v2, v3 := s[n-2], s[n-3]
		t.state.stack = append(s[:n-3], v1, v3, v2, v1)
	case rt.OpDup2:
		if n < 1 {
			return fmt.Errorf("dup2 with an empty stack")
		}
		if s[n-1].width() == 2 {
			t.state.stack = insertStack(s, n, s[n-1])
			return nil
		}
		if n < 2 || s[n-2].width() != 1 {
			return fmt.Errorf("dup2 requires two category-one values or one category-two value")
		}
		v1, v2 := s[n-1], s[n-2]
		t.state.stack = append(s[:n-2], v2, v1, v2, v1)
	case rt.OpDup2X1:
		if n < 2 {
			return fmt.Errorf("dup2_x1 with an empty stack")
		}
		v1 := s[n-1]
		if v1.width() == 2 {
			if s[n-2].width() != 1 {
				return fmt.Errorf("dup2_x1 form is not legal for these stack widths")
			}
			v2 := s[n-2]
			t.state.stack = append(s[:n-2], v1, v2, v1)
			return nil
		}
		if n < 3 || s[n-2].width() != 1 || s[n-3].width() != 1 {
			return fmt.Errorf("dup2_x1 requires three category-one values or one category-two and one category-one value")
		}
		v2, v3 := s[n-2], s[n-3]
		t.state.stack = append(s[:n-3], v2, v1, v3, v2, v1)
	case rt.OpDup2X2:
		if n < 2 {
			return fmt.Errorf("dup2_x2 with an empty stack")
		}
		v1 := s[n-1]
		if v1.width() == 2 {
			if s[n-2].width() == 2 {
				v2 := s[n-2]
				t.state.stack = append(s[:n-2], v1, v2, v1)
				return nil
			}
			if n < 3 || s[n-2].width() != 1 || s[n-3].width() != 1 {
				return fmt.Errorf("dup2_x2 form is not legal for these stack widths")
			}
			v2, v3 := s[n-2], s[n-3]
			t.state.stack = append(s[:n-3], v1, v3, v2, v1)
			return nil
		}
		if s[n-2].width() != 1 {
			return fmt.Errorf("dup2_x2 form is not legal for these stack widths")
		}
		v2 := s[n-2]
		if s[n-3].width() == 2 {
			v3 := s[n-3]
			t.state.stack = append(s[:n-3], v2, v1, v3, v2, v1)
			return nil
		}
		if n < 4 || s[n-3].width() != 1 || s[n-4].width() != 1 {
			return fmt.Errorf("dup2_x2 requires four category-one values or the category-two forms")
		}
		v3, v4 := s[n-3], s[n-4]
		t.state.stack = append(s[:n-4], v2, v1, v4, v3, v2, v1)
	case rt.OpSwap:
		if n < 2 || s[n-1].width() != 1 || s[n-2].width() != 1 {
			return fmt.Errorf("swap requires two category-one values")
		}
		v1, v2 := s[n-1], s[n-2]
		t.state.stack = append(s[:n-2], v1, v2)
	}
	return nil
}

func (t *typeChecker) binary(kind byte, what string) error {
	if _, err := t.popOf(kind, what); err != nil {
		return err
	}
	if _, err := t.popOf(kind, what); err != nil {
		return err
	}
	return t.push(vtype{kind: kind})
}

func (t *typeChecker) shift(valueKind byte, what string) error {
	if _, err := t.state.popKind(vtInt, "int shift distance"); err != nil {
		return err
	}
	if _, err := t.popOf(valueKind, what); err != nil {
		return err
	}
	return t.push(vtype{kind: valueKind})
}

func (t *typeChecker) convert(from, to vtype) error {
	if _, err := t.popOf(from.kind, from.String()); err != nil {
		return err
	}
	return t.push(to)
}

func (t *typeChecker) arrayLoad(elementKind byte, primitive string) error {
	if _, err := t.state.popKind(vtInt, "int array index"); err != nil {
		return err
	}
	if _, err := t.popArray(primitive); err != nil {
		return err
	}
	return t.push(vtype{kind: elementKind})
}

func (t *typeChecker) arrayStore(valueKind byte, what, primitive string) error {
	if _, err := t.popOf(valueKind, what); err != nil {
		return err
	}
	if _, err := t.state.popKind(vtInt, "int array index"); err != nil {
		return err
	}
	_, err := t.popArray(primitive)
	return err
}

func (t *typeChecker) aaload() error {
	if _, err := t.state.popKind(vtInt, "int array index"); err != nil {
		return err
	}
	arr, err := t.popRef(false)
	if err != nil {
		return err
	}
	component, err := arrayComponentType(arr)
	if err != nil {
		return err
	}
	return t.push(component)
}

func (t *typeChecker) aastore() error {
	value, err := t.popRef(false)
	if err != nil {
		return err
	}
	if _, err := t.state.popKind(vtInt, "int array index"); err != nil {
		return err
	}
	arr, err := t.popRef(false)
	if err != nil {
		return err
	}
	component, err := arrayComponentType(arr)
	if err != nil {
		return err
	}
	if component.kind != vtRef {
		return fmt.Errorf("aastore into an array with primitive component type %s", component)
	}
	// JVMS 4.10.1.9 requires only that the value and the array component be
	// reference types; the component compatibility check is performed by the
	// aastore instruction at run time (ArrayStoreException).
	_ = value
	return nil
}

func (t *typeChecker) applyField() error {
	ref, err := t.ref(t.insn.Ref)
	if err != nil {
		return err
	}
	if ref.Kind != rt.RefField {
		return fmt.Errorf("instruction does not reference a field")
	}
	ft, err := t.fieldType(ref)
	if err != nil {
		return err
	}
	switch t.insn.Op {
	case rt.OpGetstatic:
		return t.push(ft)
	case rt.OpPutstatic:
		value, err := t.pop()
		if err != nil {
			return err
		}
		if !assignable(t.ctx, value, ft) {
			return fmt.Errorf("putstatic value %s is not assignable to %s", value, ft)
		}
		return nil
	case rt.OpGetfield:
		recv, err := t.popRef(false)
		if err != nil {
			return err
		}
		if !assignable(t.ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
			return fmt.Errorf("getfield receiver %s is not assignable to %s", recv, ref.Owner)
		}
		return t.push(ft)
	case rt.OpPutfield:
		value, err := t.pop()
		if err != nil {
			return err
		}
		if !assignable(t.ctx, value, ft) {
			return fmt.Errorf("putfield value %s is not assignable to %s", value, ft)
		}
		recv, err := t.popRef(true)
		if err != nil {
			return err
		}
		if recv.kind == vtUninitThis {
			if ref.Owner != t.class.Name {
				return fmt.Errorf("putfield on uninitializedThis may only write a field declared by %s, not %s", t.class.Name, ref.Owner)
			}
			return nil
		}
		if recv.kind != vtRef && recv.kind != vtNull {
			return fmt.Errorf("putfield receiver is not a reference: %s", recv)
		}
		if !assignable(t.ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
			return fmt.Errorf("putfield receiver %s is not assignable to %s", recv, ref.Owner)
		}
		return nil
	}
	return fmt.Errorf("unexpected field opcode 0x%02x", t.insn.Op)
}

func (t *typeChecker) applyInvoke() error {
	ref, err := t.ref(t.insn.Ref)
	if err != nil {
		return err
	}
	if ref.Kind != rt.RefMethod {
		return fmt.Errorf("instruction does not reference a method")
	}
	mt, err := t.methodType(ref)
	if err != nil {
		return err
	}
	if t.insn.Op == rt.OpInvokeinterface {
		if got := int(t.insn.B); got != mt.ParamSlots+1 {
			return fmt.Errorf("invokeinterface count %d does not match %d argument slots", got, mt.ParamSlots+1)
		}
	}
	for i := len(mt.Params) - 1; i >= 0; i-- {
		want, err := vtypeFromType(mt.Params[i])
		if err != nil {
			return err
		}
		got, err := t.pop()
		if err != nil {
			return err
		}
		if !assignable(t.ctx, got, want) {
			return fmt.Errorf("argument %d is %s, not assignable to %s", i, got, want)
		}
	}
	if t.insn.Op == rt.OpInvokestatic {
		if ref.Name == "<init>" || ref.Name == "<clinit>" {
			return fmt.Errorf("static invocation of %s", ref.Name)
		}
		return t.pushReturn(mt)
	}
	recv, err := t.popRef(true)
	if err != nil {
		return err
	}
	if ref.Name == "<clinit>" {
		return fmt.Errorf("invocation of a class initializer")
	}
	if ref.Name == "<init>" {
		if t.insn.Op != rt.OpInvokespecial {
			return fmt.Errorf("constructor invocation requires invokespecial")
		}
		if mt.Return.Kind != 'V' {
			return fmt.Errorf("constructor descriptor must return void")
		}
		var initialized string
		switch recv.kind {
		case vtUninitThis:
			if ref.Owner != t.class.Name && ref.Owner != t.class.Super {
				return fmt.Errorf("uninitializedThis may only call the constructor of %s or its direct superclass, not %s", t.class.Name, ref.Owner)
			}
			initialized = t.class.Name
		case vtUninit:
			if ref.Owner != recv.class {
				return fmt.Errorf("constructor %s does not match allocated class %s", ref.Owner, recv.class)
			}
			initialized = recv.class
		default:
			return fmt.Errorf("constructor invocation requires an uninitialized receiver, found %s", recv)
		}
		replaceUninitialized(t.state, recv, vtype{kind: vtRef, class: initialized})
		return nil
	}
	if recv.kind == vtNull {
		// A null receiver is statically assignable to any reference type.
	} else if recv.kind != vtRef {
		return fmt.Errorf("method receiver is not a reference: %s", recv)
	} else if !assignable(t.ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
		return fmt.Errorf("method receiver %s is not assignable to %s", recv, ref.Owner)
	}
	if ref.Name == "<init>" {
		return fmt.Errorf("unreachable")
	}
	return t.pushReturn(mt)
}

func replaceUninitialized(state *frameState, from, to vtype) {
	for i, v := range state.locals {
		if v.kind == from.kind && v.off == from.off {
			state.locals[i] = to
		}
	}
	for i, v := range state.stack {
		if v.kind == from.kind && v.off == from.off {
			state.stack[i] = to
		}
	}
}

// applyStep verifies one instruction and returns its successor states.
func (t *typeChecker) applyStep() (succs []successor, terminal bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("internal verifier error: %v", r)
		}
	}()
	pc, insn := t.pc, t.insn
	addTarget := func(target int32) {
		succs = append(succs, successor{pc: target, state: copyState(t.state)})
	}
	addNext := func() { addTarget(pc + insn.Len) }
	push := t.push

	switch insn.Op {
	case rt.OpNop:
		addNext()
	case rt.OpAconstNull:
		if err := push(vtype{kind: vtNull}); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2, rt.OpIconst3,
		rt.OpIconst4, rt.OpIconst5, rt.OpBipush, rt.OpSipush:
		if err := push(vtype{kind: vtInt}); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpLconst0, rt.OpLconst1:
		if err := push(vtype{kind: vtLong}); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpFconst0, rt.OpFconst1, rt.OpFconst2:
		if err := push(vtype{kind: vtFloat}); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpDconst0, rt.OpDconst1:
		if err := push(vtype{kind: vtDouble}); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpLdc, rt.OpLdcW, rt.OpLdc2W:
		ref, err := t.ref(insn.Ref)
		if err != nil {
			return nil, false, err
		}
		var value vtype
		switch ref.Kind {
		case rt.RefInt:
			if insn.Op == rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc2_w cannot load an int constant")
			}
			value = vtype{kind: vtInt}
		case rt.RefFloat:
			if insn.Op == rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc2_w cannot load a float constant")
			}
			value = vtype{kind: vtFloat}
		case rt.RefLong:
			if insn.Op != rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc cannot load a long constant; use ldc2_w")
			}
			value = vtype{kind: vtLong}
		case rt.RefDouble:
			if insn.Op != rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc cannot load a double constant; use ldc2_w")
			}
			value = vtype{kind: vtDouble}
		case rt.RefString:
			if insn.Op == rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc2_w cannot load a String constant")
			}
			value = vtype{kind: vtRef, class: "java/lang/String"}
		case rt.RefClass:
			if insn.Op == rt.OpLdc2W {
				return nil, false, fmt.Errorf("ldc2_w cannot load a Class constant")
			}
			if !t.classExists(ref.Class) {
				return nil, false, fmt.Errorf("ldc of unknown class %s", ref.Class)
			}
			value = vtype{kind: vtRef, class: "java/lang/Class"}
		default:
			return nil, false, fmt.Errorf("unsupported ldc constant kind %d", ref.Kind)
		}
		if err := push(value); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpIload, rt.OpFload, rt.OpLload, rt.OpDload, rt.OpAload:
		idx, _ := rt.LocalIndex(insn.Op, insn.A)
		switch insn.Op {
		case rt.OpIload:
			if err := t.loadLocal(idx, vtype{kind: vtInt}); err != nil {
				return nil, false, err
			}
		case rt.OpFload:
			if err := t.loadLocal(idx, vtype{kind: vtFloat}); err != nil {
				return nil, false, err
			}
		case rt.OpLload:
			if err := t.loadLocal(idx, vtype{kind: vtLong}); err != nil {
				return nil, false, err
			}
		case rt.OpDload:
			if err := t.loadLocal(idx, vtype{kind: vtDouble}); err != nil {
				return nil, false, err
			}
		case rt.OpAload:
			if idx < 0 || int(idx) >= len(t.state.locals) {
				return nil, false, fmt.Errorf("local index %d out of range for max_locals %d", idx, len(t.state.locals))
			}
			got := t.state.locals[idx]
			if got.kind != vtRef && got.kind != vtNull && got.kind != vtUninit && got.kind != vtUninitThis {
				if got.kind == vtTop && idx > 0 && t.state.locals[idx-1].isWide() {
					return nil, false, fmt.Errorf("aload local index %d is the unusable second slot of the %s at local %d", idx, t.state.locals[idx-1], idx-1)
				}
				return nil, false, fmt.Errorf("aload local index %d holds %s, not a reference", idx, got)
			}
			if err := push(got); err != nil {
				return nil, false, err
			}
		}
	case rt.OpIstore, rt.OpFstore, rt.OpLstore, rt.OpDstore, rt.OpAstore:
		idx, _ := rt.LocalIndex(insn.Op, insn.A)
		var v vtype
		switch insn.Op {
		case rt.OpIstore:
			v, err = t.popOf(vtInt, "int")
		case rt.OpFstore:
			v, err = t.popOf(vtFloat, "float")
		case rt.OpLstore:
			v, err = t.popOf(vtLong, "long")
		case rt.OpDstore:
			v, err = t.popOf(vtDouble, "double")
		case rt.OpAstore:
			v, err = t.popRef(true)
		}
		if err != nil {
			return nil, false, err
		}
		if err := t.storeLocal(idx, v); err != nil {
			return nil, false, err
		}
		addNext()
	case rt.OpIaload:
		err = t.arrayLoad(vtInt, "I")
	case rt.OpLaload:
		err = t.arrayLoad(vtLong, "J")
	case rt.OpFaload:
		err = t.arrayLoad(vtFloat, "F")
	case rt.OpDaload:
		err = t.arrayLoad(vtDouble, "D")
	case rt.OpBaload:
		if _, err := t.state.popKind(vtInt, "int array index"); err != nil {
			return nil, false, err
		}
		arr, e := t.popRef(false)
		if e != nil {
			return nil, false, e
		}
		if arr.kind != vtNull && arr.class != "[B" && arr.class != "[Z" {
			return nil, false, fmt.Errorf("expected a byte[] or boolean[] array, found %s", arr)
		}
		err = push(vtype{kind: vtInt})
	case rt.OpCaload:
		err = t.arrayLoad(vtInt, "C")
	case rt.OpSaload:
		err = t.arrayLoad(vtInt, "S")
	case rt.OpAaload:
		err = t.aaload()
	case rt.OpIastore:
		err = t.arrayStore(vtInt, "int", "I")
	case rt.OpLastore:
		err = t.arrayStore(vtLong, "long", "J")
	case rt.OpFastore:
		err = t.arrayStore(vtFloat, "float", "F")
	case rt.OpDastore:
		err = t.arrayStore(vtDouble, "double", "D")
	case rt.OpAastore:
		err = t.aastore()
	case rt.OpBastore:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		if _, e := t.state.popKind(vtInt, "int array index"); e != nil {
			return nil, false, e
		}
		arr, e := t.popRef(false)
		if e != nil {
			return nil, false, e
		}
		if arr.kind != vtNull && arr.class != "[B" && arr.class != "[Z" {
			return nil, false, fmt.Errorf("expected a byte[] or boolean[] array, found %s", arr)
		}
	case rt.OpCastore:
		err = t.arrayStore(vtInt, "int", "C")
	case rt.OpSastore:
		err = t.arrayStore(vtInt, "int", "S")
	case rt.OpPop, rt.OpPop2, rt.OpDup, rt.OpDupX1, rt.OpDupX2, rt.OpDup2,
		rt.OpDup2X1, rt.OpDup2X2, rt.OpSwap:
		err = t.stackManipulation()
	case rt.OpIadd, rt.OpIsub, rt.OpImul, rt.OpIdiv, rt.OpIrem, rt.OpIand, rt.OpIor, rt.OpIxor:
		err = t.binary(vtInt, "int")
	case rt.OpLadd, rt.OpLsub, rt.OpLmul, rt.OpLdiv, rt.OpLrem, rt.OpLand, rt.OpLor, rt.OpLxor:
		err = t.binary(vtLong, "long")
	case rt.OpFadd, rt.OpFsub, rt.OpFmul, rt.OpFdiv, rt.OpFrem:
		err = t.binary(vtFloat, "float")
	case rt.OpDadd, rt.OpDsub, rt.OpDmul, rt.OpDdiv, rt.OpDrem:
		err = t.binary(vtDouble, "double")
	case rt.OpIshl, rt.OpIshr, rt.OpIushr:
		err = t.shift(vtInt, "int")
	case rt.OpLshl, rt.OpLshr, rt.OpLushr:
		err = t.shift(vtLong, "long")
	case rt.OpIneg, rt.OpI2b, rt.OpI2c, rt.OpI2s:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtInt})
	case rt.OpLneg:
		if _, e := t.popOf(vtLong, "long"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtLong})
	case rt.OpFneg:
		if _, e := t.popOf(vtFloat, "float"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtFloat})
	case rt.OpDneg:
		if _, e := t.popOf(vtDouble, "double"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtDouble})
	case rt.OpIinc:
		idx, _ := rt.LocalIndex(rt.OpIinc, insn.A)
		if idx < 0 || int(idx) >= len(t.state.locals) {
			return nil, false, fmt.Errorf("iinc local index %d out of range for max_locals %d", idx, len(t.state.locals))
		}
		if t.state.locals[idx].kind != vtInt {
			return nil, false, fmt.Errorf("iinc local index %d holds %s, not int", idx, t.state.locals[idx])
		}
	case rt.OpI2l:
		err = t.convert(vtype{kind: vtInt}, vtype{kind: vtLong})
	case rt.OpI2f:
		err = t.convert(vtype{kind: vtInt}, vtype{kind: vtFloat})
	case rt.OpI2d:
		err = t.convert(vtype{kind: vtInt}, vtype{kind: vtDouble})
	case rt.OpL2i:
		err = t.convert(vtype{kind: vtLong}, vtype{kind: vtInt})
	case rt.OpL2f:
		err = t.convert(vtype{kind: vtLong}, vtype{kind: vtFloat})
	case rt.OpL2d:
		err = t.convert(vtype{kind: vtLong}, vtype{kind: vtDouble})
	case rt.OpF2i:
		err = t.convert(vtype{kind: vtFloat}, vtype{kind: vtInt})
	case rt.OpF2l:
		err = t.convert(vtype{kind: vtFloat}, vtype{kind: vtLong})
	case rt.OpF2d:
		err = t.convert(vtype{kind: vtFloat}, vtype{kind: vtDouble})
	case rt.OpD2i:
		err = t.convert(vtype{kind: vtDouble}, vtype{kind: vtInt})
	case rt.OpD2l:
		err = t.convert(vtype{kind: vtDouble}, vtype{kind: vtLong})
	case rt.OpD2f:
		err = t.convert(vtype{kind: vtDouble}, vtype{kind: vtFloat})
	case rt.OpLcmp:
		if _, e := t.popOf(vtLong, "long"); e != nil {
			return nil, false, e
		}
		if _, e := t.popOf(vtLong, "long"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtInt})
	case rt.OpFcmpl, rt.OpFcmpg:
		if _, e := t.popOf(vtFloat, "float"); e != nil {
			return nil, false, e
		}
		if _, e := t.popOf(vtFloat, "float"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtInt})
	case rt.OpDcmpl, rt.OpDcmpg:
		if _, e := t.popOf(vtDouble, "double"); e != nil {
			return nil, false, e
		}
		if _, e := t.popOf(vtDouble, "double"); e != nil {
			return nil, false, e
		}
		err = push(vtype{kind: vtInt})
	case rt.OpIfeq, rt.OpIfne, rt.OpIflt, rt.OpIfge, rt.OpIfgt, rt.OpIfle:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		addTarget(insn.A)
		addNext()
		return succs, false, nil
	case rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfIcmplt, rt.OpIfIcmpge, rt.OpIfIcmpgt, rt.OpIfIcmple:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		addTarget(insn.A)
		addNext()
		return succs, false, nil
	case rt.OpIfAcmpeq, rt.OpIfAcmpne:
		if _, e := t.popRef(true); e != nil {
			return nil, false, e
		}
		if _, e := t.popRef(true); e != nil {
			return nil, false, e
		}
		addTarget(insn.A)
		addNext()
		return succs, false, nil
	case rt.OpIfnull, rt.OpIfnonnull:
		if _, e := t.popRef(true); e != nil {
			return nil, false, e
		}
		addTarget(insn.A)
		addNext()
		return succs, false, nil
	case rt.OpGoto, rt.OpGotoW:
		addTarget(insn.A)
		return succs, false, nil
	case rt.OpTableswitch, rt.OpLookupswitch:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		for _, target := range insn.Targets {
			addTarget(target)
		}
		addTarget(insn.A)
		return succs, false, nil
	case rt.OpIreturn:
		if t.mt.Return.Kind != 'I' {
			return nil, false, fmt.Errorf("ireturn in a method returning %s", t.mt.Return.Descriptor)
		}
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		return nil, true, nil
	case rt.OpLreturn:
		if t.mt.Return.Kind != 'J' {
			return nil, false, fmt.Errorf("lreturn in a method returning %s", t.mt.Return.Descriptor)
		}
		if _, e := t.popOf(vtLong, "long"); e != nil {
			return nil, false, e
		}
		return nil, true, nil
	case rt.OpFreturn:
		if t.mt.Return.Kind != 'F' {
			return nil, false, fmt.Errorf("freturn in a method returning %s", t.mt.Return.Descriptor)
		}
		if _, e := t.popOf(vtFloat, "float"); e != nil {
			return nil, false, e
		}
		return nil, true, nil
	case rt.OpDreturn:
		if t.mt.Return.Kind != 'D' {
			return nil, false, fmt.Errorf("dreturn in a method returning %s", t.mt.Return.Descriptor)
		}
		if _, e := t.popOf(vtDouble, "double"); e != nil {
			return nil, false, e
		}
		return nil, true, nil
	case rt.OpReturn:
		if t.mt.Return.Kind != 'V' {
			return nil, false, fmt.Errorf("return in a method returning %s", t.mt.Return.Descriptor)
		}
		return nil, true, nil
	case rt.OpAreturn:
		v, e := t.popRef(false)
		if e != nil {
			return nil, false, e
		}
		want, e := vtypeFromType(t.mt.Return)
		if e != nil {
			return nil, false, e
		}
		if !assignable(t.ctx, v, want) {
			return nil, false, fmt.Errorf("areturn value %s is not assignable to %s", v, want)
		}
		return nil, true, nil
	case rt.OpGetstatic, rt.OpPutstatic, rt.OpGetfield, rt.OpPutfield:
		if e := t.applyField(); e != nil {
			return nil, false, e
		}
		addNext()
		return succs, false, nil
	case rt.OpInvokestatic, rt.OpInvokespecial, rt.OpInvokevirtual, rt.OpInvokeinterface:
		if e := t.applyInvoke(); e != nil {
			return nil, false, e
		}
		addNext()
		return succs, false, nil
	case rt.OpNew:
		ref, e := t.ref(insn.Ref)
		if e != nil {
			return nil, false, e
		}
		if strings.HasPrefix(ref.Class, "[") {
			return nil, false, fmt.Errorf("new of array class %s", ref.Class)
		}
		if !t.classExists(ref.Class) {
			return nil, false, fmt.Errorf("new of unknown class %s", ref.Class)
		}
		if e := push(vtype{kind: vtUninit, class: ref.Class, off: pc}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpNewarray:
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		desc, ok := rt.NewarrayDescriptor(insn.A)
		if !ok {
			return nil, false, fmt.Errorf("invalid newarray atype %d", insn.A)
		}
		if e := push(vtype{kind: vtRef, class: desc}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpAnewarray:
		ref, e := t.ref(insn.Ref)
		if e != nil {
			return nil, false, e
		}
		if _, e := t.popOf(vtInt, "int"); e != nil {
			return nil, false, e
		}
		if !t.classExists(ref.Class) {
			return nil, false, fmt.Errorf("anewarray of unknown component %s", ref.Class)
		}
		if e := push(vtype{kind: vtRef, class: arrayDescriptor(ref.Class)}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpMultianewarray:
		ref, e := t.ref(insn.Ref)
		if e != nil {
			return nil, false, e
		}
		if !strings.HasPrefix(ref.Class, "[") {
			return nil, false, fmt.Errorf("multianewarray class %s is not an array type", ref.Class)
		}
		rank := 0
		for rank < len(ref.Class) && ref.Class[rank] == '[' {
			rank++
		}
		dims := int(insn.B)
		if dims < 1 || dims > rank {
			return nil, false, fmt.Errorf("multianewarray dimensions %d do not match array type %s", dims, ref.Class)
		}
		if _, e := rt.ParseFieldType(ref.Class); e != nil {
			return nil, false, e
		}
		for i := 0; i < dims; i++ {
			if _, e := t.popOf(vtInt, "int"); e != nil {
				return nil, false, e
			}
		}
		if e := push(vtype{kind: vtRef, class: ref.Class}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpArraylength:
		v, e := t.popRef(false)
		if e != nil {
			return nil, false, e
		}
		if v.kind != vtNull && (!strings.HasPrefix(v.class, "[") || !isArrayDesc(v.class)) {
			return nil, false, fmt.Errorf("arraylength on non-array %s", v)
		}
		if e := push(vtype{kind: vtInt}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpAthrow:
		v, e := t.popRef(false)
		if e != nil {
			return nil, false, e
		}
		if v.kind == vtRef && !classAssignable(t.ctx, v.class, "java/lang/Throwable") {
			return nil, false, fmt.Errorf("athrow of %s which is not a Throwable", v.class)
		}
		return nil, true, nil
	case rt.OpCheckcast:
		ref, e := t.ref(insn.Ref)
		if e != nil {
			return nil, false, e
		}
		if _, e := t.popRef(false); e != nil {
			return nil, false, e
		}
		if !t.classExists(ref.Class) {
			return nil, false, fmt.Errorf("checkcast of unknown class %s", ref.Class)
		}
		if e := push(vtype{kind: vtRef, class: ref.Class}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpInstanceof:
		ref, e := t.ref(insn.Ref)
		if e != nil {
			return nil, false, e
		}
		if _, e := t.popRef(false); e != nil {
			return nil, false, e
		}
		if !t.classExists(ref.Class) {
			return nil, false, fmt.Errorf("instanceof of unknown class %s", ref.Class)
		}
		if e := push(vtype{kind: vtInt}); e != nil {
			return nil, false, e
		}
		addNext()
	case rt.OpInvokedynamic:
		return nil, false, fmt.Errorf("invokedynamic/bootstrap constants are unsupported")
	case rt.OpMonitorenter, rt.OpMonitorexit:
		return nil, false, fmt.Errorf("monitor instruction 0x%02x is unsupported", insn.Op)
	default:
		return nil, false, fmt.Errorf("unsupported opcode 0x%02x", insn.Op)
	}
	if err != nil {
		return nil, false, err
	}
	addNext()
	return succs, false, nil
}

func isArrayDesc(desc string) bool {
	_, err := rt.ParseFieldType(desc)
	return err == nil && strings.HasPrefix(desc, "[")
}
