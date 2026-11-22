// Package verify performs Gate 1 structural verification and stack-height
// analysis over decoded method bodies. Both engines consume the same verified
// representation.
package verify

import (
	"fmt"
	"sort"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// Context provides the class set needed to resolve verification dependencies.
type Context struct {
	Classes map[string]*classfile.Class
}

// ClassExists reports whether name is an admitted application class, a
// supported bootstrap class, or a valid reference array descriptor.
func (c *Context) ClassExists(name string) bool {
	if strings.HasPrefix(name, "[") {
		_, err := rt.ParseFieldType(name)
		return err == nil
	}
	if rt.BootstrapHasClass(name) {
		return true
	}
	_, ok := c.Classes[name]
	return ok
}

// VerifyClass checks one class's member descriptors, methods, exception tables
// and stack maps. It does not resolve inherited members; the linker does that.
func VerifyClass(c *classfile.Class, ctx *Context) error {
	if c.Flags&rt.AccInterface != 0 {
		// Interface method bodies are outside the profile (only abstract
		// declarations are admitted).
		if c.Super != "" && c.Super != "java/lang/Object" {
			return fmt.Errorf("%s: interface superclass %s is unsupported", c.Name, c.Super)
		}
	}
	for _, f := range c.Fields {
		if _, err := rt.ParseFieldType(f.Descriptor); err != nil {
			return fmt.Errorf("%s.%s: %w", c.Name, f.Name, err)
		}
	}
	for _, m := range c.Methods {
		if err := verifyMethod(c, m, ctx); err != nil {
			return err
		}
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
		return fmt.Errorf("%s: synchronized methods are unsupported", display)
	}
	if m.Flags&rt.AccNative != 0 {
		return fmt.Errorf("%s: native methods are unsupported", display)
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
	// Instructions must cover the code array exactly and contiguously.
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

	// Exception table structure.
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
		if h.CatchType != "" && !ctx.ClassExists(h.CatchType) {
			return fmt.Errorf("%s: exception handler catch type %s is not available", display, h.CatchType)
		}
	}

	// Stack map frames must point at instruction boundaries and be ordered.
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

	// Stack-height data-flow. Gate 1 has category-one values only.
	heights := map[int32]int{}
	queue := []int32{0}
	setState := func(pc int32, h int) error {
		if pc < 0 || pc >= int32(m.CodeLength) {
			return fmt.Errorf("%s: control transfers to invalid offset %d", display, pc)
		}
		if h < 0 || h > m.MaxStack {
			return fmt.Errorf("%s: pc %d: operand stack depth %d exceeds max_stack %d", display, pc, h, m.MaxStack)
		}
		if old, ok := heights[pc]; ok {
			if old != h {
				return fmt.Errorf("%s: pc %d: inconsistent operand stack depth %d versus %d", display, pc, h, old)
			}
			return nil
		}
		heights[pc] = h
		queue = append(queue, pc)
		return nil
	}
	if err := setState(0, 0); err != nil {
		return err
	}
	// Exception handlers are reachable with exactly the throwable on the stack.
	for _, h := range m.Handlers {
		if err := setState(h.HandlerPC, 1); err != nil {
			return err
		}
	}

	for len(queue) > 0 {
		pc := queue[0]
		queue = queue[1:]
		i, ok := index[pc]
		if !ok {
			continue
		}
		insn := m.Insns[i]
		h := heights[pc]
		pops, pushes, terminal, err := stackEffect(insn, m)
		if err != nil {
			return fmt.Errorf("%s: pc %d: %w", display, pc, err)
		}
		if h < pops {
			return fmt.Errorf("%s: pc %d: operand stack underflow", display, pc)
		}
		nh := h - pops + pushes
		if nh < 0 || nh > m.MaxStack {
			return fmt.Errorf("%s: pc %d: operand stack depth %d exceeds max_stack %d", display, pc, nh, m.MaxStack)
		}
		if insn.Op == rt.OpInvokedynamic {
			return fmt.Errorf("%s: pc %d: invokedynamic/bootstrap is unsupported", display, pc)
		}
		target, hasTarget := classfile.BranchTarget(insn)
		switch insn.Op {
		case rt.OpGoto:
			if _, ok := index[target]; !ok {
				return fmt.Errorf("%s: pc %d: branch target %d is not an instruction boundary", display, pc, target)
			}
			if err := setState(target, nh); err != nil {
				return err
			}
			continue
		case rt.OpIfeq, rt.OpIfne, rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfAcmpeq, rt.OpIfAcmpne:
			if _, ok := index[target]; !ok {
				return fmt.Errorf("%s: pc %d: branch target %d is not an instruction boundary", display, pc, target)
			}
			if err := setState(target, nh); err != nil {
				return err
			}
		}
		_ = hasTarget
		if terminal {
			continue
		}
		nextPC := pc + insn.Len
		if int(nextPC) >= m.CodeLength {
			return fmt.Errorf("%s: pc %d: control falls off the end of the method", display, pc)
		}
		if err := setState(nextPC, nh); err != nil {
			return err
		}
	}

	// Validate declared frames against computed heights where reachable, and
	// require frames at every branch target and handler when a table exists.
	if m.HasStackMap {
		for _, f := range m.Frames {
			if got, ok := heights[f.Offset]; ok && got != len(f.Stack) {
				return fmt.Errorf("%s: stack map frame at offset %d declares stack size %d but control reaches it with %d", display, f.Offset, len(f.Stack), got)
			}
		}
		for _, h := range m.Handlers {
			if _, ok := frameAt[h.HandlerPC]; !ok {
				return fmt.Errorf("%s: stack map has no frame at exception handler %d", display, h.HandlerPC)
			}
		}
		for _, insn := range m.Insns {
			if t, ok := classfile.BranchTarget(insn); ok {
				if t == 0 {
					continue
				}
				if t < int32(m.CodeLength) {
					if _, ok := frameAt[t]; !ok {
						return fmt.Errorf("%s: pc %d: stack map has no frame at branch target %d", display, insn.PC, t)
					}
				}
			}
		}
	}

	// Return instructions must agree with the method descriptor.
	for _, insn := range m.Insns {
		switch insn.Op {
		case rt.OpReturn:
			if mt.Return.Kind != 'V' {
				return fmt.Errorf("%s: pc %d: return in a non-void method", display, insn.PC)
			}
		case rt.OpIreturn:
			if mt.Return.Kind != 'I' {
				return fmt.Errorf("%s: pc %d: ireturn in a method returning %s", display, insn.PC, mt.Return.Descriptor)
			}
		case rt.OpAreturn:
			if mt.Return.Kind != 'L' && mt.Return.Kind != 'A' {
				return fmt.Errorf("%s: pc %d: areturn in a method returning %s", display, insn.PC, mt.Return.Descriptor)
			}
		}
	}
	if err := verifyTypes(c, m, mt, ctx, heights, frameAt); err != nil {
		return fmt.Errorf("%s: %w", display, err)
	}
	return nil
}

// Heights recomputes verified operand-stack heights for AOT lowering.
func Heights(m *classfile.Method) (map[int32]int, error) {
	index := map[int32]int{}
	for i, insn := range m.Insns {
		index[insn.PC] = i
	}
	heights := map[int32]int{0: 0}
	queue := []int32{0}
	set := func(pc int32, h int) error {
		if h < 0 || h > m.MaxStack {
			return fmt.Errorf("pc %d: operand stack depth %d exceeds max_stack %d", pc, h, m.MaxStack)
		}
		if old, ok := heights[pc]; ok {
			if old != h {
				return fmt.Errorf("pc %d: inconsistent stack depth %d versus %d", pc, h, old)
			}
			return nil
		}
		heights[pc] = h
		queue = append(queue, pc)
		return nil
	}
	for _, h := range m.Handlers {
		if err := set(h.HandlerPC, 1); err != nil {
			return nil, err
		}
	}
	for len(queue) > 0 {
		pc := queue[0]
		queue = queue[1:]
		i, ok := index[pc]
		if !ok {
			continue
		}
		insn := m.Insns[i]
		pops, pushes, terminal, err := stackEffect(insn, m)
		if err != nil {
			return nil, fmt.Errorf("pc %d: %w", pc, err)
		}
		h := heights[pc]
		if h < pops {
			return nil, fmt.Errorf("pc %d: operand stack underflow", pc)
		}
		nh := h - pops + pushes
		if target, ok := classfile.BranchTarget(insn); ok {
			if _, exists := index[target]; !exists {
				return nil, fmt.Errorf("pc %d: branch target %d is not an instruction boundary", pc, target)
			}
			if err := set(target, nh); err != nil {
				return nil, err
			}
			if insn.Op == rt.OpGoto {
				continue
			}
		}
		if terminal {
			continue
		}
		next := pc + insn.Len
		if int(next) >= m.CodeLength {
			return nil, fmt.Errorf("pc %d: control falls off the end of the method", pc)
		}
		if err := set(next, nh); err != nil {
			return nil, err
		}
	}
	return heights, nil
}

func stackEffect(insn rt.Insn, m *classfile.Method) (pops, pushes int, terminal bool, err error) {
	switch insn.Op {
	case rt.OpAconstNull, rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2,
		rt.OpIconst3, rt.OpIconst4, rt.OpIconst5, rt.OpBipush, rt.OpSipush,
		rt.OpLdc, rt.OpLdcW:
		return 0, 1, false, nil
	case 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d, rt.OpIload, rt.OpAload:
		idx, _ := localIndex(insn.Op, insn.A)
		if int(idx) >= m.MaxLocals {
			return 0, 0, false, fmt.Errorf("local index %d exceeds max_locals %d", idx, m.MaxLocals)
		}
		return 0, 1, false, nil
	case 0x3b, 0x3c, 0x3d, 0x3e, 0x4b, 0x4c, 0x4d, 0x4e, rt.OpIstore, rt.OpAstore:
		idx, _ := localIndex(insn.Op, insn.A)
		if int(idx) >= m.MaxLocals {
			return 0, 0, false, fmt.Errorf("local index %d exceeds max_locals %d", idx, m.MaxLocals)
		}
		return 1, 0, false, nil
	case rt.OpPop:
		return 1, 0, false, nil
	case rt.OpDup:
		return 1, 2, false, nil
	case rt.OpIadd, rt.OpIrem:
		return 2, 1, false, nil
	case rt.OpNew:
		return 0, 1, false, nil
	case rt.OpGetfield:
		return 1, 1, false, nil
	case rt.OpPutfield:
		return 2, 0, false, nil
	case rt.OpGetstatic:
		return 0, 1, false, nil
	case rt.OpPutstatic:
		return 1, 0, false, nil
	case rt.OpAnewarray:
		return 1, 1, false, nil
	case rt.OpAaload:
		return 2, 1, false, nil
	case rt.OpAastore:
		return 3, 0, false, nil
	case rt.OpArraylength:
		return 1, 1, false, nil
	case rt.OpAthrow:
		return 1, 0, true, nil
	case rt.OpReturn:
		return 0, 0, true, nil
	case rt.OpIreturn, rt.OpAreturn:
		return 1, 0, true, nil
	case rt.OpIfeq, rt.OpIfne:
		return 1, 0, false, nil
	case rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfAcmpeq, rt.OpIfAcmpne:
		return 2, 0, false, nil
	case rt.OpGoto:
		return 0, 0, true, nil
	case rt.OpInvokestatic, rt.OpInvokevirtual, rt.OpInvokespecial, rt.OpInvokeinterface:
		ref, err := methodRefOf(m, insn)
		if err != nil {
			return 0, 0, false, err
		}
		mt, err := rt.ParseMethodType(ref.Descriptor)
		if err != nil {
			return 0, 0, false, err
		}
		pops = mt.ArgCount
		if insn.Op != rt.OpInvokestatic {
			pops++
		}
		if mt.Return.Kind != 'V' {
			pushes = 1
		}
		return pops, pushes, false, nil
	case rt.OpInvokedynamic:
		return 0, 0, false, nil
	default:
		return 0, 0, false, fmt.Errorf("unsupported opcode 0x%02x", insn.Op)
	}
}

func localIndex(op byte, operand int32) (int32, bool) {
	switch op {
	case rt.OpIload, rt.OpAload, rt.OpIstore, rt.OpAstore:
		return operand, true
	}
	switch {
	case op >= 0x1a && op <= 0x1d:
		return int32(op - 0x1a), true
	case op >= 0x2a && op <= 0x2d:
		return int32(op - 0x2a), true
	case op >= 0x3b && op <= 0x3e:
		return int32(op - 0x3b), true
	case op >= 0x4b && op <= 0x4e:
		return int32(op - 0x4b), true
	}
	return 0, false
}

func methodRefOf(m *classfile.Method, insn rt.Insn) (rt.Ref, error) {
	if insn.Ref < 0 || int(insn.Ref) >= len(m.Refs) {
		return rt.Ref{}, fmt.Errorf("missing method reference")
	}
	ref := m.Refs[insn.Ref]
	if ref.Kind != rt.RefMethod {
		return rt.Ref{}, fmt.Errorf("instruction does not reference a method")
	}
	return ref, nil
}

// SortedFrames is a helper for deterministic diagnostics.
func SortedFrames(frames []classfile.Frame) []classfile.Frame {
	out := append([]classfile.Frame(nil), frames...)
	sort.Slice(out, func(i, j int) bool { return out[i].Offset < out[j].Offset })
	return out
}
