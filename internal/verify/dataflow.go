package verify

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/classfile"
)

// canonicalFrame is a StackMapTable frame lifted into simulation types,
// with category-2 slots expanded ([value][top]).
type canonicalFrame struct {
	locals []vtype
	stack  []vtype
}

// frame is a simulated operand state.
type frame struct {
	locals []vtype
	stack  []vtype
}

func (f *frame) clone() *frame {
	return &frame{
		locals: append([]vtype(nil), f.locals...),
		stack:  append([]vtype(nil), f.stack...),
	}
}

func (f *frame) push(v vtype) {
	f.stack = append(f.stack, v)
	if isCat2(v) {
		f.stack = append(f.stack, tTop)
	}
}

func (f *frame) pop() (vtype, error) {
	if len(f.stack) == 0 {
		return vtype{}, fmt.Errorf("stack underflow")
	}
	v := f.stack[len(f.stack)-1]
	f.stack = f.stack[:len(f.stack)-1]
	return v, nil
}

func (f *frame) popCat2() (vtype, error) {
	if len(f.stack) < 2 {
		return vtype{}, fmt.Errorf("stack underflow (cat2)")
	}
	top := f.stack[len(f.stack)-1]
	v := f.stack[len(f.stack)-2]
	if top.kind != kTop || !isCat2(v) {
		return vtype{}, fmt.Errorf("expected category-2 pair, got %s/%s", v, top)
	}
	f.stack = f.stack[:len(f.stack)-2]
	return v, nil
}

func kindName(k vkind) string {
	return vtype{kind: k}.String()
}

func (f *frame) popExpect(want vkind, what string) (vtype, error) {
	if want == kUnknown {
		return f.pop()
	}
	v, err := f.pop()
	if err != nil {
		return vtype{}, fmt.Errorf("%s: %w", what, err)
	}
	switch want {
	case kInt, kFloat, kLong, kDouble:
		if v.kind != want && v.kind != kUnknown {
			return vtype{}, fmt.Errorf("%s: want %s, got %s", what, kindName(want), v)
		}
	case kObj:
		if !isRef(v) {
			return vtype{}, fmt.Errorf("%s: want reference, got %s", what, v)
		}
	default:
		if v.kind != want && v.kind != kUnknown {
			return vtype{}, fmt.Errorf("%s: want %s, got %s", what, kindName(want), v)
		}
	}
	return v, nil
}

func (f *frame) clearStack() { f.stack = f.stack[:0] }

func equalFrames(a, b *frame) bool {
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

// --- SM frame lifting ----------------------------------------------------------

// canonicalizeFrames walks the frames in offset order and produces the
// absolute local state per target. Relative forms (same/append/chop) are
// applied over the running absolute locals, which start as the method's
// expanded entry locals (JVMS §4.7.4 semantics).
func (c *checker) canonicalizeFrames(m *classfile.MethodInfo) (
	map[int]*canonicalFrame, *frame, error,
) {
	code := m.Code
	entry, err := c.entryState(m)
	if err != nil {
		return nil, nil, err
	}
	abs := append([]vtype(nil), entry.locals...)
	out := make(map[int]*canonicalFrame, len(code.StackMaps))

	for i := range code.StackMaps {
		fr := &code.StackMaps[i]
		ft := fr.Kind
		var cf canonicalFrame

		switch {
		case ft <= 63: // same_frame: abs unchanged
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		case ft <= 127: // same_locals_1_stack_item: locals unchanged
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		case ft == 247:
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		case ft <= 250: // chop (251-ft) locals
			k := 251 - int(ft)
			if k > len(abs) {
				return nil, nil, verr(c.mn, int(fr.Offset),
					"chop %d of %d locals", k, len(abs))
			}
			abs = abs[:len(abs)-k]
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		case ft == 251:
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		case ft <= 254: // append (ft-251) locals
			for _, it := range fr.Locals {
				v, err2 := c.vitemToVType(it)
				if err2 != nil {
					return nil, nil, verr(c.mn, int(fr.Offset), "%v", err2)
				}
				abs = append(abs, v)
				if isCat2(v) {
					abs = append(abs, tTop)
				}
			}
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		default: // full_frame: replaces the absolute locals entirely
			abs = abs[:0]
			for _, it := range fr.Locals {
				v, err2 := c.vitemToVType(it)
				if err2 != nil {
					return nil, nil, verr(c.mn, int(fr.Offset), "%v", err2)
				}
				abs = append(abs, v)
				if isCat2(v) {
					abs = append(abs, tTop)
				}
			}
			cf = canonicalFrame{locals: append([]vtype(nil), abs...)}
		}

		for _, it := range fr.Stack {
			v, err2 := c.vitemToVType(it)
			if err2 != nil {
				return nil, nil, verr(c.mn, int(fr.Offset), "%v", err2)
			}
			cf.stack = append(cf.stack, v)
			if isCat2(v) {
				cf.stack = append(cf.stack, tTop)
			}
		}
		if os.Getenv("CATTY_VTRACE") != "" {
			fmt.Fprintf(os.Stderr, "[vtrace] frame off=%d kind=%d absLen=%d stackLen=%d\n",
				fr.Offset, ft, len(cf.locals), len(cf.stack))
		}
		out[int(fr.Offset)] = &cf
	}
	return out, entry, nil
}

func (c *checker) vitemToVType(it classfile.VItem) (vtype, error) {
	switch it.Tag {
	case classfile.VItemTop:
		return tTop, nil
	case classfile.VItemInteger:
		return tInt, nil
	case classfile.VItemFloat:
		return tFloat, nil
	case classfile.VItemLong:
		return tLong, nil
	case classfile.VItemDouble:
		return tDouble, nil
	case classfile.VItemNull:
		return tNull, nil
	case classfile.VItemUninitTh:
		return tUninitT, nil
	case classfile.VItemObject:
		name, err := c.cf.ClassName(it.CPoolIdx)
		if err != nil {
			return vtype{}, err
		}
		return tObj(name), nil
	case classfile.VItemUninit:
		return tUninit(it.Offset), nil
	}
	return vtype{}, fmt.Errorf("bad verification item tag %d", it.Tag)
}

// --- entry state ------------------------------------------------------------------

func (c *checker) entryState(m *classfile.MethodInfo) (*frame, error) {
	args, _, err := splitDesc(m.Desc)
	if err != nil {
		return nil, verr(c.mn, 0, "%v", err)
	}
	f := &frame{}
	add := func(v vtype) {
		f.locals = append(f.locals, v)
		if isCat2(v) {
			f.locals = append(f.locals, tTop)
		}
	}
	if m.AccessFlags&classfile.AccStatic == 0 {
		if m.Name == "<init>" {
			add(tUninitT)
		} else {
			add(tObj(c.cf.ThisClass))
		}
	}
	for _, a := range args {
		add(vtypeForDesc(a))
	}
	if f.locals == nil {
		f.locals = []vtype{}
	}
	return f, nil
}

func vtypeForDesc(d string) vtype {
	switch d[0] {
	case 'J':
		return tLong
	case 'D':
		return tDouble
	case 'F':
		return tFloat
	case 'L':
		return tObj(d[1 : len(d)-1])
	case '[':
		return tObj(d)
	default:
		return tInt
	}
}

// retVType extracts the return descriptor.
func retVType(desc string) (vtype, bool, error) {
	_, ret, err := splitDesc(desc)
	if err != nil {
		return vtype{}, false, err
	}
	if ret == "V" {
		return vtype{}, false, nil
	}
	return vtypeForDesc(ret), true, nil
}

// --- dataflow driver --------------------------------------------------------------

// splitDesc splits a method descriptor into argument type descriptors and
// return descriptor (shared with the structural tier).
func splitDesc(desc string) (args []string, ret string, err error) {
	if len(desc) == 0 || desc[0] != '(' {
		return nil, "", fmt.Errorf("bad descriptor %q", desc)
	}
	i := 1
	for i < len(desc) && desc[i] != ')' {
		start := i
		for i < len(desc) && desc[i] == '[' {
			i++
		}
		if i >= len(desc) {
			return nil, "", fmt.Errorf("bad descriptor %q", desc)
		}
		if desc[i] == 'L' {
			for i < len(desc) && desc[i] != ';' {
				i++
			}
			if i >= len(desc) {
				return nil, "", fmt.Errorf("bad descriptor %q", desc)
			}
		}
		i++
		args = append(args, desc[start:i])
	}
	if i >= len(desc) || desc[i] != ')' {
		return nil, "", fmt.Errorf("bad descriptor %q", desc)
	}
	return args, desc[i+1:], nil
}

func setLocalPair(st *frame, idx int, v vtype) {
	if os.Getenv("CATTY_VTRACE") != "" && idx == 3 {
		fmt.Fprintf(os.Stderr, "[setLocal3] v=%v lenBefore=%v\n", v, len(st.locals))
	}
	for len(st.locals) <= idx {
		st.locals = append(st.locals, tTop)
	}
	st.locals[idx] = v
	if isCat2(v) && len(st.locals) == idx+1 {
		st.locals = append(st.locals, tTop)
	}
}

// checkMethodDataflow simulates the method against its StackMapTable frames.
func (c *checker) checkMethodDataflow(m *classfile.MethodInfo) error {
	code := m.Code
	smAt := make(map[int]*classfile.StackMapFrame, len(code.StackMaps))
	for i := range code.StackMaps {
		fr := &code.StackMaps[i]
		smAt[int(fr.Offset)] = fr
	}

	entry, err := c.entryState(m)
	if err != nil {
		return err
	}
	canonical, entry, err := c.canonicalizeFrames(m)
	if err != nil {
		return err
	}
	states := map[int]*frame{0: entry}
	if cn, ok := canonical[0]; ok {
		if cerr := c.stateCompatible(entry, *cn); cerr != nil {
			return verr(c.mn, 0, "entry frame mismatch: %v", cerr)
		}
		states[0] = c.frameFromCanonical(cn)
	}

	visited := make(map[int]bool)
	queue := []int{0}
	for len(queue) > 0 {
		pc := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if visited[pc] {
			continue
		}
		if os.Getenv("CATTY_VMETHOD") != "" && strings.Contains(c.mn, os.Getenv("CATTY_VMETHOD")) {
			fmt.Fprintf(os.Stderr, "[sim] %s pc=%d locals=%v stack=%v\n", c.mn, pc, states[pc].locals, states[pc].stack)
		}
		visited[pc] = true
		// Simulate on a working copy: effects mutate their operand frame,
		// and the archived state must stay pre-execution.
		working := states[pc]
		states[pc] = working.clone()
		succs, err := c.simulate(m, pc, working)
		if err != nil {
			return err
		}
		for _, sp := range succs {
			if sp.pc < 0 || sp.pc >= len(code.Code) {
				return verr(c.mn, pc, "successor pc %d out of range", sp.pc)
			}
			ns := sp.st
			if cn, ok := canonical[sp.pc]; ok {
				if err := c.stateCompatible(ns, *cn); err != nil {
					return verr(c.mn, sp.pc,
						"merge from pc=%d: %v\n  simulated locals=%v stack=%v\n  frame     locals=%v stack=%v",
						pc, err, ns.locals, ns.stack, cn.locals, cn.stack)
				}
				ns = c.frameFromCanonical(cn)
			}
			if prev, ok := states[sp.pc]; ok {
				if !equalFrames(prev, ns) {
					if os.Getenv("CATTY_VTRACE") != "" {
						fmt.Fprintf(os.Stderr, "[vtrace] revisit pc=%d from %d\n  prev locals=%v stack=%v\n  new  locals=%v stack=%v\n",
							sp.pc, pc, prev.locals, prev.stack, ns.locals, ns.stack)
					}
					return verr(c.mn, sp.pc, "inconsistent state on revisit (from pc=%d)", pc)
				}
				continue
			}
			if os.Getenv("CATTY_VTRACE") != "" {
				fmt.Fprintf(os.Stderr, "[enq] %s pc=%d from=%d locals=%v\n", c.mn, sp.pc, pc, ns.locals)
			}
			states[sp.pc] = ns
			queue = append(queue, sp.pc)
		}
	}
	return nil
}

func (c *checker) frameFromCanonical(cf *canonicalFrame) *frame {
	return &frame{locals: append([]vtype(nil), cf.locals...), stack: append([]vtype(nil), cf.stack...)}
}

type successor struct {
	pc int
	st *frame
}

// replaceUninit swaps every uninit@off marker for the initialized type after
// a successful <init> invocation (JVMS §4.10.1.6).
func replaceUninit(st *frame, off uint16, with vtype) {
	for i := range st.locals {
		if st.locals[i].kind == kUninit &&
			(off == uint16max || st.locals[i].off == off) {
			st.locals[i] = with
		} else if st.locals[i].kind == kUninitThis && off == uint16max {
			st.locals[i] = with
		}
	}
	for i := range st.stack {
		if st.stack[i].kind == kUninit && st.stack[i].off == off {
			st.stack[i] = with
		}
	}
}

const uint16max uint16 = 0xFFFF
