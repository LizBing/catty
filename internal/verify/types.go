package verify

import (
	"fmt"
	"sort"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// The Gate 1 verifier tracks category-one int and reference types, Java null,
// uninitializedThis, and each uninitialized object produced by `new`. Stack map
// frames are authoritative at branch and handler entries; each incoming normal
// edge must be assignable to the declared frame.

const (
	vtTop byte = iota
	vtInt
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

func (v vtype) String() string {
	switch v.kind {
	case vtTop:
		return "top"
	case vtInt:
		return "int"
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

func initialState(c *classfile.Class, m *classfile.Method, mt *rt.MethodType) *frameState {
	locals := make([]vtype, m.MaxLocals)
	for i := range locals {
		locals[i] = vtype{kind: vtTop}
	}
	slot := 0
	if m.Flags&rt.AccStatic == 0 {
		if m.Name == "<init>" {
			locals[0] = vtype{kind: vtUninitThis}
		} else {
			locals[0] = vtype{kind: vtRef, class: c.Name}
		}
		slot = 1
	}
	for _, p := range mt.Params {
		vt, err := vtypeFromType(p)
		if err != nil {
			return nil
		}
		if slot < len(locals) {
			locals[slot] = vt
		}
		slot++
	}
	return &frameState{locals: locals}
}

func frameStateFrom(f *classfile.Frame, maxLocals int) (*frameState, error) {
	locals := make([]vtype, maxLocals)
	for i := range locals {
		locals[i] = vtype{kind: vtTop}
	}
	for i, lv := range f.Locals {
		if i >= maxLocals {
			return nil, fmt.Errorf("stack map locals exceed max_locals")
		}
		vt, err := vtypeFromFrame(lv)
		if err != nil {
			return nil, err
		}
		locals[i] = vt
	}
	stack := make([]vtype, len(f.Stack))
	for i, sv := range f.Stack {
		vt, err := vtypeFromFrame(sv)
		if err != nil {
			return nil, err
		}
		stack[i] = vt
	}
	return &frameState{locals: locals, stack: stack}, nil
}

func copyState(s *frameState) *frameState {
	out := &frameState{
		locals: append([]vtype(nil), s.locals...),
		stack:  append([]vtype(nil), s.stack...),
	}
	return out
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

func (s *frameState) popInt() error {
	v, err := s.pop()
	if err != nil {
		return err
	}
	if v.kind != vtInt {
		return fmt.Errorf("expected int on the operand stack, found %s", v)
	}
	return nil
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
			return fmt.Errorf("local %d is %s but the stack map expects %s", i, "uninitialized", want)
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
			fc, fok := arrayComponentName(from)
			tc, tok := arrayComponentName(to)
			if !fok || !tok {
				return false
			}
			return classAssignable(ctx, fc, tc)
		}
		return false
	}
	if strings.HasPrefix(to, "[") {
		return false
	}
	if classIsInterface(ctx, to) {
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
	var visitClass func(string) bool
	visitInterface := func(string) bool { return false }
	var visitInterfaceRec func(string) bool
	visitInterfaceRec = func(name string) bool {
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
			if visitInterfaceRec(in) {
				return true
			}
		}
		return false
	}
	visitInterface = visitInterfaceRec
	visitClass = func(name string) bool {
		for name != "" {
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
	return visitClass(from)
}

func arrayComponentName(desc string) (string, bool) {
	if !strings.HasPrefix(desc, "[") {
		return "", false
	}
	component := desc[1:]
	if strings.HasPrefix(component, "L") && strings.HasSuffix(component, ";") {
		return component[1 : len(component)-1], true
	}
	if strings.HasPrefix(component, "[") {
		return component, true
	}
	return "", false
}

func arrayDescriptor(component string) string {
	if strings.HasPrefix(component, "[") {
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
	component := v.class[1:]
	switch {
	case strings.HasPrefix(component, "L") && strings.HasSuffix(component, ";"):
		return vtype{kind: vtRef, class: component[1 : len(component)-1]}, nil
	case strings.HasPrefix(component, "["):
		return vtype{kind: vtRef, class: component}, nil
	default:
		return vtype{}, fmt.Errorf("primitive arrays are outside the Gate 1 profile")
	}
}

// verifyTypes performs receiver/operand type checking against the verified
// stack-height map and declared stack-map frames.
func verifyTypes(c *classfile.Class, m *classfile.Method, mt *rt.MethodType, ctx *Context, heights map[int32]int, frameAt map[int32]*classfile.Frame) error {
	index := map[int32]int{}
	for i, insn := range m.Insns {
		index[insn.PC] = i
	}
	handlerPC := map[int32]bool{}
	for _, h := range m.Handlers {
		handlerPC[h.HandlerPC] = true
	}
	var pcs []int32
	for pc := range heights {
		pcs = append(pcs, pc)
	}
	sort.Slice(pcs, func(i, j int) bool { return pcs[i] < pcs[j] })

	state := initialState(c, m, mt)
	initial := copyState(state)
	stateValid := true
	first := true
	for _, pc := range pcs {
		i, ok := index[pc]
		if !ok {
			continue
		}
		if f := frameAt[pc]; f != nil {
			fs, err := frameStateFrom(f, m.MaxLocals)
			if err != nil {
				return fmt.Errorf("pc %d: %w", pc, err)
			}
			if !first && !handlerPC[pc] && stateValid {
				if err := checkAssignableState(ctx, state, fs); err != nil {
					return fmt.Errorf("pc %d: incoming state does not match stack map: %w", pc, err)
				}
			}
			if first {
				if err := checkAssignableState(ctx, initial, fs); err != nil {
					return fmt.Errorf("pc 0: initial state does not match stack map: %w", err)
				}
			}
			state = fs
			stateValid = true
		} else if !first && !stateValid {
			return fmt.Errorf("pc %d: reachable control-flow join has no stack map frame", pc)
		}
		first = false
		insn := m.Insns[i]
		after, terminal, err := applyTypes(ctx, c, m, mt, pc, insn, state)
		if err != nil {
			return fmt.Errorf("pc %d: %w", pc, err)
		}
		state = after
		if target, ok := classfile.BranchTarget(insn); ok {
			var targetState *frameState
			if target == 0 {
				targetState = initial
			} else if f := frameAt[target]; f != nil {
				targetState, err = frameStateFrom(f, m.MaxLocals)
				if err != nil {
					return fmt.Errorf("pc %d: %w", pc, err)
				}
			} else {
				return fmt.Errorf("pc %d: branch target %d has no stack map frame", pc, target)
			}
			if err := checkAssignableState(ctx, state, targetState); err != nil {
				return fmt.Errorf("pc %d: branch to %d: %w", pc, target, err)
			}
		}
		stateValid = !terminal
	}
	return nil
}

func applyTypes(ctx *Context, c *classfile.Class, m *classfile.Method, mt *rt.MethodType, pc int32, insn rt.Insn, state *frameState) (*frameState, bool, error) {
	switch insn.Op {
	case rt.OpAconstNull:
		state.push(vtype{kind: vtNull})
	case rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2, rt.OpIconst3, rt.OpIconst4, rt.OpIconst5,
		rt.OpBipush, rt.OpSipush:
		state.push(vtype{kind: vtInt})
	case rt.OpLdc, rt.OpLdcW:
		if insn.Ref < 0 || int(insn.Ref) >= len(m.Refs) {
			return nil, false, fmt.Errorf("missing ldc constant")
		}
		ref := m.Refs[insn.Ref]
		switch ref.Kind {
		case rt.RefInt:
			state.push(vtype{kind: vtInt})
		case rt.RefString:
			state.push(vtype{kind: vtRef, class: "java/lang/String"})
		case rt.RefClass:
			state.push(vtype{kind: vtRef, class: "java/lang/Class"})
		default:
			return nil, false, fmt.Errorf("unsupported ldc constant")
		}
	case rt.OpIload, rt.OpAload, 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d:
		idx, _ := localIndex(insn.Op, insn.A)
		if int(idx) >= len(state.locals) {
			return nil, false, fmt.Errorf("local index %d out of range", idx)
		}
		v := state.locals[idx]
		switch insn.Op {
		case rt.OpIload, 0x1a, 0x1b, 0x1c, 0x1d:
			if v.kind != vtInt {
				return nil, false, fmt.Errorf("iload of non-int local %s", v)
			}
		default:
			if v.kind != vtRef && v.kind != vtNull && v.kind != vtUninit && v.kind != vtUninitThis {
				return nil, false, fmt.Errorf("aload of non-reference local %s", v)
			}
		}
		state.push(v)
	case rt.OpIstore, rt.OpAstore, 0x3b, 0x3c, 0x3d, 0x3e, 0x4b, 0x4c, 0x4d, 0x4e:
		idx, _ := localIndex(insn.Op, insn.A)
		if int(idx) >= len(state.locals) {
			return nil, false, fmt.Errorf("local index %d out of range", idx)
		}
		v, err := state.pop()
		if err != nil {
			return nil, false, err
		}
		switch insn.Op {
		case rt.OpIstore, 0x3b, 0x3c, 0x3d, 0x3e:
			if v.kind != vtInt {
				return nil, false, fmt.Errorf("istore of non-int value %s", v)
			}
		default:
			if v.kind != vtRef && v.kind != vtNull && v.kind != vtUninit && v.kind != vtUninitThis {
				return nil, false, fmt.Errorf("astore of non-reference value %s", v)
			}
		}
		state.locals[idx] = v
	case rt.OpPop:
		if _, err := state.pop(); err != nil {
			return nil, false, err
		}
	case rt.OpDup:
		if len(state.stack) == 0 {
			return nil, false, fmt.Errorf("dup with an empty stack")
		}
		state.push(state.stack[len(state.stack)-1])
	case rt.OpIadd, rt.OpIrem:
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		state.push(vtype{kind: vtInt})
	case rt.OpNew:
		ref := m.Refs[insn.Ref]
		if !ctx.ClassExists(ref.Class) {
			return nil, false, fmt.Errorf("new of unknown class %s", ref.Class)
		}
		if classIsInterface(ctx, ref.Class) {
			return nil, false, fmt.Errorf("new of interface %s", ref.Class)
		}
		state.push(vtype{kind: vtUninit, class: ref.Class, off: pc})
	case rt.OpGetfield:
		ref := m.Refs[insn.Ref]
		ft, err := fieldValueType(ref)
		if err != nil {
			return nil, false, err
		}
		recv, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
			return nil, false, fmt.Errorf("getfield receiver %s is not assignable to %s", recv, ref.Owner)
		}
		state.push(ft)
	case rt.OpPutfield:
		ref := m.Refs[insn.Ref]
		ft, err := fieldValueType(ref)
		if err != nil {
			return nil, false, err
		}
		value, err := state.pop()
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, value, ft) {
			return nil, false, fmt.Errorf("putfield value %s is not assignable to %s", value, ft)
		}
		recv, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
			return nil, false, fmt.Errorf("putfield receiver %s is not assignable to %s", recv, ref.Owner)
		}
	case rt.OpGetstatic:
		ref := m.Refs[insn.Ref]
		ft, err := fieldValueType(ref)
		if err != nil {
			return nil, false, err
		}
		state.push(ft)
	case rt.OpPutstatic:
		ref := m.Refs[insn.Ref]
		ft, err := fieldValueType(ref)
		if err != nil {
			return nil, false, err
		}
		value, err := state.pop()
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, value, ft) {
			return nil, false, fmt.Errorf("putstatic value %s is not assignable to %s", value, ft)
		}
	case rt.OpAnewarray:
		ref := m.Refs[insn.Ref]
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		if !ctx.ClassExists(ref.Class) {
			return nil, false, fmt.Errorf("anewarray of unknown component %s", ref.Class)
		}
		state.push(vtype{kind: vtRef, class: arrayDescriptor(ref.Class)})
	case rt.OpAaload:
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		arr, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		component, err := arrayComponentType(arr)
		if err != nil {
			return nil, false, err
		}
		state.push(component)
	case rt.OpAastore:
		value, err := state.pop()
		if err != nil {
			return nil, false, err
		}
		if value.kind != vtRef && value.kind != vtNull {
			return nil, false, fmt.Errorf("aastore of non-reference value %s", value)
		}
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		arr, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		component, err := arrayComponentType(arr)
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, value, component) {
			return nil, false, fmt.Errorf("aastore value %s is not assignable to component %s", value, component)
		}
	case rt.OpArraylength:
		if _, err := state.popRef(false); err != nil {
			return nil, false, err
		}
		state.push(vtype{kind: vtInt})
	case rt.OpInvokestatic, rt.OpInvokevirtual, rt.OpInvokespecial, rt.OpInvokeinterface:
		return applyInvoke(ctx, c, m, pc, insn, state)
	case rt.OpReturn:
		if mt.Return.Kind != 'V' {
			return nil, false, fmt.Errorf("return in a method returning %s", mt.Return.Descriptor)
		}
		return state, true, nil
	case rt.OpIreturn:
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		if mt.Return.Kind != 'I' {
			return nil, false, fmt.Errorf("ireturn in a method returning %s", mt.Return.Descriptor)
		}
		return state, true, nil
	case rt.OpAreturn:
		v, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		want, err := vtypeFromType(mt.Return)
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, v, want) {
			return nil, false, fmt.Errorf("areturn value %s is not assignable to %s", v, want)
		}
		return state, true, nil
	case rt.OpAthrow:
		v, err := state.popRef(false)
		if err != nil {
			return nil, false, err
		}
		if v.kind == vtRef && !classAssignable(ctx, v.class, "java/lang/Throwable") {
			return nil, false, fmt.Errorf("athrow of %s which is not a Throwable", v.class)
		}
		return state, true, nil
	case rt.OpIfeq, rt.OpIfne:
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
	case rt.OpIfIcmpeq, rt.OpIfIcmpne:
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
		if err := state.popInt(); err != nil {
			return nil, false, err
		}
	case rt.OpIfAcmpeq, rt.OpIfAcmpne:
		if _, err := state.popRef(true); err != nil {
			return nil, false, err
		}
		if _, err := state.popRef(true); err != nil {
			return nil, false, err
		}
	case rt.OpGoto:
		return state, true, nil
	default:
		return nil, false, fmt.Errorf("unsupported opcode 0x%02x", insn.Op)
	}
	return state, false, nil
}

func fieldValueType(ref rt.Ref) (vtype, error) {
	ft, err := rt.ParseFieldType(ref.Descriptor)
	if err != nil {
		return vtype{}, err
	}
	return vtypeFromType(ft)
}

func applyInvoke(ctx *Context, c *classfile.Class, m *classfile.Method, pc int32, insn rt.Insn, state *frameState) (*frameState, bool, error) {
	ref := m.Refs[insn.Ref]
	mt, err := rt.ParseMethodType(ref.Descriptor)
	if err != nil {
		return nil, false, err
	}
	for i := len(mt.Params) - 1; i >= 0; i-- {
		want, err := vtypeFromType(mt.Params[i])
		if err != nil {
			return nil, false, err
		}
		got, err := state.pop()
		if err != nil {
			return nil, false, err
		}
		if !assignable(ctx, got, want) {
			return nil, false, fmt.Errorf("argument %d %s is not assignable to %s", i, got, want)
		}
	}
	if insn.Op == rt.OpInvokestatic {
		if ref.Name == "<init>" || ref.Name == "<clinit>" {
			return nil, false, fmt.Errorf("static invocation of %s", ref.Name)
		}
		if err := pushReturn(state, mt); err != nil {
			return nil, false, err
		}
		return state, false, nil
	}
	recv, err := state.popRef(true)
	if err != nil {
		return nil, false, err
	}
	if ref.Name == "<init>" {
		if insn.Op != rt.OpInvokespecial {
			return nil, false, fmt.Errorf("constructor invocation requires invokespecial")
		}
		if mt.Return.Kind != 'V' {
			return nil, false, fmt.Errorf("constructor descriptor must return void")
		}
		var initialized string
		switch recv.kind {
		case vtUninitThis:
			if ref.Owner != c.Name && ref.Owner != c.Super {
				return nil, false, fmt.Errorf("uninitializedThis may only call the current or direct superclass constructor, not %s", ref.Owner)
			}
			initialized = c.Name
		case vtUninit:
			if ref.Owner != recv.class {
				return nil, false, fmt.Errorf("constructor %s does not match allocated class %s", ref.Owner, recv.class)
			}
			initialized = recv.class
		default:
			return nil, false, fmt.Errorf("constructor invocation requires an uninitialized receiver, found %s", recv)
		}
		replaceUninitialized(state, recv, vtype{kind: vtRef, class: initialized})
		return state, false, nil
	}
	if recv.kind != vtRef && recv.kind != vtNull {
		return nil, false, fmt.Errorf("method receiver is not a reference: %s", recv)
	}
	if !assignable(ctx, recv, vtype{kind: vtRef, class: ref.Owner}) {
		return nil, false, fmt.Errorf("method receiver %s is not assignable to %s", recv, ref.Owner)
	}
	if err := pushReturn(state, mt); err != nil {
		return nil, false, err
	}
	return state, false, nil
}

func pushReturn(state *frameState, mt *rt.MethodType) error {
	if mt.Return.Kind == 'V' {
		return nil
	}
	v, err := vtypeFromType(mt.Return)
	if err != nil {
		return err
	}
	state.push(v)
	return nil
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
