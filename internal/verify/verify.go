// Package verify performs Gate 2 structural verification and typed data-flow
// analysis over decoded method bodies. Both engines consume the same verified
// representation; the analysis also produces the logical operand-stack height
// at every reachable instruction for the generated Go emitter.
package verify

import (
	"fmt"
	"sort"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// Context provides the class set needed to resolve verification dependencies
// and collects the verified stack heights used by code generation.
type Context struct {
	Classes map[string]*classfile.Class
	Profile string
	Work    *Work
	workErr error
	// MemberFlags resolves a member's access flags and declaring class for the
	// protected-receiver check. Admission supplies it; image builds leave it nil.
	MemberFlags func(owner, name, descriptor string) (flags uint16, declaring string, ok bool)
	// RuntimePackage reports same-runtime-package identity including the
	// defining loader when admission can distinguish loaders.
	RuntimePackage func(a, b string) bool
	Heights        map[*classfile.Method]map[int32]int
	Shapes         map[*classfile.Method]map[int32][]int
}

// ClassExists reports whether name is an admitted application class, a
// supported bootstrap class, or a valid array descriptor.
func (c *Context) profile() string {
	if c.Profile == "" {
		return rt.Gate2ProfileID
	}
	return c.Profile
}

func (c *Context) ClassExists(name string) bool {
	if strings.HasPrefix(name, "[") {
		_, err := rt.ParseFieldType(name)
		return err == nil
	}
	if rt.BootstrapHasClassFor(c.profile(), name) {
		return true
	}
	_, ok := c.Classes[name]
	return ok
}

func (c *Context) ensureHeights() {
	if c.Heights == nil {
		c.Heights = map[*classfile.Method]map[int32]int{}
	}
	if c.Shapes == nil {
		c.Shapes = map[*classfile.Method]map[int32][]int{}
	}
}

// Heights returns the verified logical stack height at each reachable
// instruction offset for m.
func (c *Context) HeightsOf(m *classfile.Method) (map[int32]int, bool) {
	h, ok := c.Heights[m]
	return h, ok
}

// ShapesOf returns the verified logical operand widths at instruction entry.
func (c *Context) ShapesOf(m *classfile.Method) (map[int32][]int, bool) {
	s, ok := c.Shapes[m]
	return s, ok
}

func isThrowable(ctx *Context, name string) (bool, error) {
	seen := map[string]bool{}
	for name != "" && !seen[name] {
		if err := ctx.charge(1, "throwable hierarchy query"); err != nil {
			return false, err
		}
		seen[name] = true
		if name == "java/lang/Throwable" {
			return true, nil
		}
		if c, ok := ctx.Classes[name]; ok {
			name = c.Super
			continue
		}
		if info, ok := rt.BootstrapClassDefFor(ctx.profile(), name); ok {
			name = info.Super
			continue
		}
		return false, nil
	}
	return false, nil
}

// VerifyClass checks one class's member descriptors, methods, exception tables
// and stack maps. It does not resolve inherited members; the linker does that.
func VerifyClass(c *classfile.Class, ctx *Context) error {
	ctx.ensureHeights()
	if c.Flags&rt.AccInterface != 0 {
		if c.Super != "" && c.Super != "java/lang/Object" {
			return fmt.Errorf("%s: interface superclass %s is unsupported", c.Name, c.Super)
		}
	}
	for _, f := range c.Fields {
		if _, err := rt.ParseFieldType(f.Descriptor); err != nil {
			return fmt.Errorf("%s.%s: %w", c.Name, f.Name, err)
		}
	}
	for _, f := range c.Fields {
		if f.Constant != nil && f.Flags&rt.AccStatic == 0 {
			return fmt.Errorf("%s.%s: ConstantValue requires a static field", c.Name, f.Name)
		}
		if rt.ProfileHasThreads(ctx.profile()) && f.Flags&rt.AccFinal != 0 && f.Flags&rt.AccVolatile != 0 {
			return fmt.Errorf("%s.%s: a field must not be both final and volatile", c.Name, f.Name)
		}
	}
	for _, m := range c.Methods {
		if err := verifyMethod(c, m, ctx); err != nil {
			return err
		}
	}
	if err := ctx.WorkError(); err != nil {
		return err
	}
	return nil
}

func verifyMethod(c *classfile.Class, m *classfile.Method, ctx *Context) error {
	display := fmt.Sprintf("%s.%s%s", c.Name, m.Name, m.Descriptor)
	mt, err := rt.ParseMethodType(m.Descriptor)
	if err != nil {
		return fmt.Errorf("%s: %w", display, err)
	}
	if m.Flags&rt.AccSynchronized != 0 {
		if !rt.ProfileHasThreads(ctx.profile()) {
			return fmt.Errorf("%s: synchronized methods are unsupported", display)
		}
		switch {
		case m.Flags&rt.AccNative != 0:
			return fmt.Errorf("%s: synchronized native methods are unsupported", display)
		case m.Name == "<init>":
			return fmt.Errorf("%s: an instance constructor must not be synchronized", display)
		case m.Flags&rt.AccAbstract != 0:
			return fmt.Errorf("%s: an abstract method must not be synchronized", display)
		case c.Flags&rt.AccInterface != 0 && m.Name != "<clinit>":
			return fmt.Errorf("%s: an interface method must not be synchronized", display)
		}
	}
	if m.Flags&rt.AccNative != 0 {
		return fmt.Errorf("%s: native methods are unsupported", display)
	}
	if rt.ProfileHasThreads(ctx.profile()) && m.Name == "<clinit>" {
		if m.Flags&rt.AccStatic == 0 || m.Descriptor != "()V" {
			return fmt.Errorf("%s: a class initializer must be static with descriptor ()V", display)
		}
		if m.Flags&rt.AccAbstract != 0 {
			return fmt.Errorf("%s: a class initializer must not be abstract", display)
		}
	}
	abstract := m.Flags&rt.AccAbstract != 0
	if abstract {
		if m.CodeLength != 0 || len(m.Insns) != 0 {
			return fmt.Errorf("%s: abstract method has a Code attribute", display)
		}
		return nil
	}
	if m.CodeLength == 0 || len(m.Insns) == 0 {
		return fmt.Errorf("%s: concrete method has no Code attribute", display)
	}
	needLocals := mt.ParamSlots
	if m.Flags&rt.AccStatic == 0 {
		needLocals++
	}
	if m.MaxLocals < needLocals {
		return fmt.Errorf("%s: max_locals %d is smaller than the %d argument slots", display, m.MaxLocals, needLocals)
	}

	for _, insn := range m.Insns {
		if insn.Op == rt.OpInvokedynamic {
			return fmt.Errorf("%s: pc %d: invokedynamic/bootstrap constants are unsupported", display, insn.PC)
		}
	}
	index := map[int32]int{}
	for i, insn := range m.Insns {
		if insn.PC < 0 {
			return fmt.Errorf("%s: negative instruction offset %d", display, insn.PC)
		}
		if _, ok := index[insn.PC]; ok {
			return fmt.Errorf("%s: duplicate instruction offset %d", display, insn.PC)
		}
		index[insn.PC] = i
	}
	next := int32(0)
	for _, insn := range m.Insns {
		if insn.PC != next {
			return fmt.Errorf("%s: pc %d: instruction does not follow previous instruction at %d", display, insn.PC, next)
		}
		if insn.Len <= 0 {
			return fmt.Errorf("%s: pc %d: invalid instruction length", display, insn.PC)
		}
		next = insn.PC + insn.Len
	}
	if int(next) != m.CodeLength {
		return fmt.Errorf("%s: decoded code length %d differs from Code length %d", display, next, m.CodeLength)
	}
	boundary := func(pc int32) bool {
		if pc < 0 || int(pc) > m.CodeLength {
			return false
		}
		if int(pc) == m.CodeLength {
			return true
		}
		_, ok := index[pc]
		return ok
	}

	for _, h := range m.Handlers {
		if h.StartPC < 0 || h.EndPC > int32(m.CodeLength) || h.HandlerPC < 0 || h.HandlerPC >= int32(m.CodeLength) {
			return fmt.Errorf("%s: exception handler range [%d,%d)->%d is outside the code", display, h.StartPC, h.EndPC, h.HandlerPC)
		}
		if h.StartPC >= h.EndPC {
			return fmt.Errorf("%s: exception handler range [%d,%d) is empty", display, h.StartPC, h.EndPC)
		}
		if _, ok := index[h.StartPC]; !ok {
			return fmt.Errorf("%s: exception handler start %d is not an instruction boundary", display, h.StartPC)
		}
		if !boundary(h.EndPC) {
			return fmt.Errorf("%s: exception handler end %d is not an instruction boundary", display, h.EndPC)
		}
		if _, ok := index[h.HandlerPC]; !ok {
			return fmt.Errorf("%s: exception handler target %d is not an instruction boundary", display, h.HandlerPC)
		}
		if h.CatchType != "" {
			if !ctx.ClassExists(h.CatchType) {
				return fmt.Errorf("%s: exception handler catch type %s is not available", display, h.CatchType)
			}
			isThrowableType, err := isThrowable(ctx, h.CatchType)
			if err != nil {
				return err
			}
			if !isThrowableType {
				return fmt.Errorf("%s: exception handler catch type %s is not a Throwable", display, h.CatchType)
			}
		}
	}

	frameAt := map[int32]*classfile.Frame{}
	for i := range m.Frames {
		f := &m.Frames[i]
		if _, ok := index[f.Offset]; !ok {
			return fmt.Errorf("%s: stack map frame %d offset %d is not an instruction boundary", display, i, f.Offset)
		}
		if _, ok := frameAt[f.Offset]; ok {
			return fmt.Errorf("%s: duplicate stack map frame at offset %d", display, f.Offset)
		}
		frameAt[f.Offset] = f
	}

	if m.HasStackMap {
		for _, h := range m.Handlers {
			if _, ok := frameAt[h.HandlerPC]; !ok {
				return fmt.Errorf("%s: stack map has no frame at exception handler %d", display, h.HandlerPC)
			}
		}
		for _, insn := range m.Insns {
			if target, ok := classfile.BranchTarget(insn); ok {
				if target == 0 {
					continue
				}
				if target >= 0 && target < int32(m.CodeLength) {
					if _, ok := frameAt[target]; !ok {
						return fmt.Errorf("%s: pc %d: stack map has no frame at branch target %d", display, insn.PC, target)
					}
				}
			}
			if _, targets, ok := classfile.SwitchTargets(insn); ok {
				for _, target := range targets {
					if target == 0 {
						continue
					}
					if target >= 0 && target < int32(m.CodeLength) {
						if _, ok := frameAt[target]; !ok {
							return fmt.Errorf("%s: pc %d: stack map has no frame at switch target %d", display, insn.PC, target)
						}
					}
				}
				if def := insn.A; def != 0 && def >= 0 && def < int32(m.CodeLength) {
					if _, ok := frameAt[def]; !ok {
						return fmt.Errorf("%s: pc %d: stack map has no frame at switch default %d", display, insn.PC, def)
					}
				}
			}
		}
	}

	heights, shapes, err := analyzeMethod(c, m, mt, ctx, index, frameAt)
	if err != nil {
		return fmt.Errorf("%s: %w", display, err)
	}
	ctx.Heights[m] = heights
	ctx.Shapes[m] = shapes
	return nil
}

// analyzeMethod runs the typed worklist over reachable instructions and
// returns the logical operand-stack height at each reachable instruction.
func analyzeMethod(c *classfile.Class, m *classfile.Method, mt *rt.MethodType,
	ctx *Context, index map[int32]int, frameAt map[int32]*classfile.Frame) (map[int32]int, map[int32][]int, error) {
	if err := ctx.charge(int64(m.MaxLocals+m.MaxStack), "frame analysis storage"); err != nil {
		return nil, nil, err
	}
	initial, err := initialState(c, m, mt)
	if err != nil {
		return nil, nil, err
	}
	heights := map[int32]int{}
	shapes := map[int32][]int{}
	states := map[int32]*frameState{}
	processed := map[int32]bool{}
	var queue []int32

	seed := func(pc int32, state *frameState) {
		if _, ok := states[pc]; ok {
			return
		}
		states[pc] = state
		queue = append(queue, pc)
	}

	start := initial
	if f := frameAt[0]; f != nil {
		fs, err := frameStateFrom(ctx, f, m.MaxLocals)
		if err != nil {
			return nil, nil, fmt.Errorf("pc 0: %w", err)
		}
		if err := checkAssignableState(ctx, initial, fs); err != nil {
			return nil, nil, fmt.Errorf("pc 0: initial state does not match stack map: %w", err)
		}
		start = fs
	}
	seed(0, start)

	// Exception handlers are reachable with exactly the throwable on the
	// operand stack and a local state described by their stack map frame.
	handlerSeeded := map[int32]bool{}
	for _, h := range m.Handlers {
		if handlerSeeded[h.HandlerPC] {
			continue
		}
		handlerSeeded[h.HandlerPC] = true
		var hs *frameState
		if f := frameAt[h.HandlerPC]; f != nil {
			fs, err := frameStateFrom(ctx, f, m.MaxLocals)
			if err != nil {
				return nil, nil, fmt.Errorf("handler pc %d: %w", h.HandlerPC, err)
			}
			if len(fs.stack) != 1 {
				return nil, nil, fmt.Errorf("handler pc %d: stack map declares %d operand values, want exactly the throwable", h.HandlerPC, len(fs.stack))
			}
			throwable := fs.stack[0].kind == vtNull
			if !throwable && fs.stack[0].kind == vtRef {
				throwable, err = classAssignable(ctx, fs.stack[0].class, "java/lang/Throwable")
				if err != nil {
					return nil, nil, fmt.Errorf("handler pc %d: %w", h.HandlerPC, err)
				}
			}
			if !throwable {
				return nil, nil, fmt.Errorf("handler pc %d: stack map declares %s, not a Throwable", h.HandlerPC, fs.stack[0])
			}
			hs = fs
		} else if m.HasStackMap {
			return nil, nil, fmt.Errorf("handler pc %d: missing stack map frame", h.HandlerPC)
		} else {
			var err error
			hs, err = copyState(ctx, initial)
			if err != nil {
				return nil, nil, fmt.Errorf("handler pc %d: %w", h.HandlerPC, err)
			}
			hs.stack = nil
			hs.stack = append(hs.stack, vtype{kind: vtRef, class: "java/lang/Throwable"})
		}
		seed(h.HandlerPC, hs)
	}

	propagate := func(fromPC int32, target int32, state *frameState) error {
		if target < 0 || int(target) >= m.CodeLength {
			return fmt.Errorf("pc %d: control transfers to invalid offset %d", fromPC, target)
		}
		if _, ok := index[target]; !ok {
			return fmt.Errorf("pc %d: branch target %d is not an instruction boundary", fromPC, target)
		}
		if f := frameAt[target]; f != nil {
			fs, err := frameStateFrom(ctx, f, m.MaxLocals)
			if err != nil {
				return fmt.Errorf("pc %d: target %d: %w", fromPC, target, err)
			}
			if err := checkAssignableState(ctx, state, fs); err != nil {
				return fmt.Errorf("pc %d: branch to %d: %w", fromPC, target, err)
			}
			if !processed[target] {
				if _, ok := states[target]; !ok {
					states[target] = fs
					queue = append(queue, target)
				}
			}
			return nil
		}
		if existing, ok := states[target]; ok {
			equal, err := statesEqual(ctx, existing, state)
			if err != nil {
				return err
			}
			if !equal {
				return fmt.Errorf("pc %d: control-flow join at %d has no stack map frame and the incoming states differ", fromPC, target)
			}
			return nil
		}
		copied, err := copyState(ctx, state)
		if err != nil {
			return err
		}
		states[target] = copied
		queue = append(queue, target)
		return nil
	}

	for len(queue) > 0 {
		pc := queue[0]
		queue = queue[1:]
		if processed[pc] {
			continue
		}
		processed[pc] = true
		state := states[pc]
		// Exception edges use the local state at the throwing instruction and
		// exactly the caught throwable on the operand stack. Every handler whose
		// protected range covers this instruction receives the same edge check,
		// including shared handler targets and catch-all handlers.
		for _, h := range m.Handlers {
			if err := ctx.charge(1, "handler range inspection"); err != nil {
				return nil, nil, err
			}
			if pc < h.StartPC || pc >= h.EndPC {
				continue
			}
			if err := ctx.charge(1, "exception edge"); err != nil {
				return nil, nil, err
			}
			f := frameAt[h.HandlerPC]
			if f == nil {
				continue
			}
			fs, err := frameStateFrom(ctx, f, m.MaxLocals)
			if err != nil {
				return nil, nil, fmt.Errorf("exception edge at pc %d: handler pc %d: %w", pc, h.HandlerPC, err)
			}
			catch := h.CatchType
			if catch == "" {
				catch = "java/lang/Throwable"
			}
			edge, err := copyState(ctx, state)
			if err != nil {
				return nil, nil, fmt.Errorf("exception edge at pc %d: %w", pc, err)
			}
			edge.stack = []vtype{{kind: vtRef, class: catch}}
			if err := checkAssignableState(ctx, edge, fs); err != nil {
				return nil, nil, fmt.Errorf("exception edge at pc %d to handler pc %d: %w", pc, h.HandlerPC, err)
			}
		}
		heights[pc] = len(state.stack)
		widths := make([]int, len(state.stack))
		for i, v := range state.stack {
			widths[i] = v.width()
		}
		shapes[pc] = widths
		if err := ctx.charge(1, "instruction"); err != nil {
			return nil, nil, err
		}
		insn := m.Insns[index[pc]]
		stepState, err := copyState(ctx, state)
		if err != nil {
			return nil, nil, fmt.Errorf("pc %d: %w", pc, err)
		}
		t := &typeChecker{ctx: ctx, class: c, method: m, mt: mt, state: stepState, pc: pc, insn: insn}
		succs, terminal, err := t.applyStep()
		if err != nil {
			return nil, nil, fmt.Errorf("pc %d: %w", pc, err)
		}
		if err := ctx.WorkError(); err != nil {
			return nil, nil, err
		}
		if terminal {
			continue
		}
		for _, succ := range succs {
			if err := ctx.charge(1, "control-flow edge"); err != nil {
				return nil, nil, err
			}
			if err := propagate(pc, succ.pc, succ.state); err != nil {
				return nil, nil, err
			}
		}
	}
	return heights, shapes, nil
}

// SortedFrames is a helper for deterministic diagnostics.
func SortedFrames(frames []classfile.Frame) []classfile.Frame {
	out := append([]classfile.Frame(nil), frames...)
	sort.Slice(out, func(i, j int) bool { return out[i].Offset < out[j].Offset })
	return out
}
