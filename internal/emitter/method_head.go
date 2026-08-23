package emitter

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"catty/internal/classfile"
	"catty/internal/verify"
)

// emitMethodBody translates one method's bytecode into Go statements.
//
// Emission model (v1):
//   - every bytecode local and stack slot is a kernel.Value variable,
//     declared up front (`l%d`, `s%d`);
//   - branches map 1:1 onto Go labels; all labels precede any use;
//   - fallible operations route through genrt helpers that return
//     *kernel.Thrown, followed by an inline handler dispatch block whose
//     handler ranges are compile-time constants.
func emitMethodBody(cf *classfile.ClassFile, m *classfile.MethodInfo) (string, error) {
	code := m.Code
	e := &methodEmitter{
		cf:        cf,
		m:         m,
		code:      code.Code,
		maxLocals: int(code.MaxLocals),
		handlers:  code.Handlers,
	}

	targets, err := verify.BranchTargets(code.Code)
	if err != nil {
		return "", err
	}
	for h := range code.Handlers {
		targets[int(code.Handlers[h].HandlerPc)] = true
	}
	e.targets = targets
	if os.Getenv("CATTY_ETRACE") != "" {
		keys := []int{}
		for k := range targets {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		fmt.Fprintf(os.Stderr, "[targets] %s%s %v\n", e.cf.ThisClass, e.m.Name, keys)
	}

	argDescs, _, err := splitMethodDesc(e.m.Desc)
	if err != nil {
		return "", err
	}

	// Prologue: exception carrier + locals + stack slots (blank-used).
	e.p("var exc *kernel.Thrown")
	if e.maxLocals > 0 {
		e.p("var %s kernel.Value", slotList("l", e.maxLocals))
		e.p("_ = []kernel.Value{%s}", slotList("l", e.maxLocals))
	}
	e.p("var %s kernel.Value", slotList("s", int(code.MaxStack)+1))
	e.p("_ = []kernel.Value{%s}", slotList("s", int(code.MaxStack)+1))

	// Argument copying from args[] into locals.
	slot := 0
	if e.m.AccessFlags&classfile.AccStatic == 0 {
		e.p("l0 = recv")
		slot = 1
	}
	for i, a := range argDescs {
		e.p("l%d = args[%d]", slot, i)
		slot += descSlots(a)
	}

	if err := e.body(); err != nil {
		return "", err
	}

	return e.w.String(), nil
}

// localName renders the variable for bytecode local slot idx.
func (e *methodEmitter) localName(idx int) string { return fmt.Sprintf("l%d", idx) }

// localsCount returns the declared local-slot count.
func (e *methodEmitter) localsCount() int { return len(e.localsNames()) }

// localsNames exposes the declared count as a slice length stand-in.
func (e *methodEmitter) localsNames() []struct{} { return make([]struct{}, e.maxLocals) }

// prologue declares the pending-exception carrier used by F-scheme dispatch.
func (e *methodEmitter) prologue() string {
	return "\tvar exc *kernel.Thrown\n"
}

func slotList(prefix string, n int) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return strings.Join(parts, ", ")
}

func splitMethodDesc(desc string) ([]string, string, error) {
	if len(desc) == 0 || desc[0] != '(' {
		return nil, "", fmt.Errorf("bad descriptor %q", desc)
	}
	i := 1
	var args []string
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

func descSlots(d string) int {
	if d == "J" || d == "D" {
		return 2
	}
	return 1
}

type methodEmitter struct {
	cf         *classfile.ClassFile
	m          *classfile.MethodInfo
	code       []byte
	w          strings.Builder
	depth      int
	maxLocals  int
	handlers   []classfile.ExceptionHandler
	targets    map[int]bool
	mergeDepth map[int]int    // lazily computed stack depth per merge pc
	handlerAt  map[int]string // pc -> catch class name (""=any)
}

var _ = fmt.Sprintf

func (e *methodEmitter) p(format string, args ...any) {
	fmt.Fprintf(&e.w, "\t"+format+"\n", args...)
}

func (e *methodEmitter) push(expr string) {
	e.p("s%d = %s", e.depth, expr)
	e.depth++
}

func (e *methodEmitter) pop() string {
	e.depth--
	return fmt.Sprintf("s%d", e.depth)
}

func (e *methodEmitter) peekTop() string { return fmt.Sprintf("s%d", e.depth-1) }

// label emits a label if pc is a branch/handler target.
func (e *methodEmitter) label(pc int) {
	if e.targets[pc] {
		e.p("L%d:", pc)
	}
}

// excDispatch emits the inline exception dispatch for a faulting site at pc.
// exc must hold the pending *kernel.Thrown. Falls through when unhandled
// after propagating.
func (e *methodEmitter) excDispatch(pc int) {
	// Only propagate when exc is actually non-nil.
	e.p("if exc != nil {")
	for _, h := range e.handlers {
		if pc < int(h.StartPc) || pc >= int(h.EndPc) {
			continue
		}
		if h.CatchType == 0 {
			e.p("\tgoto L%d", h.HandlerPc)
			continue
		}
		cn, err := e.cf.ClassName(h.CatchType)
		if err != nil {
			continue
		}
		e.p("\tif genrt.InstanceOf(exc.Obj, %q) { goto L%d }", cn, h.HandlerPc)
	}
	e.p("\treturn nil, exc // no handler matched")
	e.p("}")
}

func (e *methodEmitter) fallible(pc int, lhs string, call string) {
	if lhs != "" {
		e.p("%s = %s", lhs, call)
	} else {
		e.p("_ = %s", call)
	}
	e.excDispatch(pc)
}

// smStackDepth counts raw stack slots including cat2 second slots.
func smStackDepth(fr *classfile.StackMapFrame) int {
	n := 0
	for _, it := range fr.Stack {
		n++
		if it.Tag == classfile.VItemLong || it.Tag == classfile.VItemDouble {
			n++
		}
	}
	return n
}

// catchName resolves a catch_type to its internal class name (" "=any).
func catchName(cf *classfile.ClassFile, idx uint16) string {
	if idx == 0 {
		return ""
	}
	n, err := cf.ClassName(idx)
	if err != nil {
		return ""
	}
	return n
}

// branchTarget resolves a two-byte signed branch offset at pc.
func branchTarget(pc int, code []byte) int {
	return pc + int(int16(uint16(code[pc+1])<<8|uint16(code[pc+2])))
}

// refAt reads a field/method reference at the given operand index.
func refAt(cf *classfile.ClassFile, code []byte, pc int) (cls, name, desc string, tag classfile.ConstTag, err error) {
	idx := uint16(code[pc+1])<<8 | uint16(code[pc+2])
	return cf.Ref(idx)
}

// returnsValue reports whether a method descriptor returns a value.
func returnsValue(desc string) (bool, error) {
	_, has, err := splitMethodDesc(desc)
	_ = has
	if err != nil {
		return false, err
	}
	i := len(desc) - 1
	return desc[i] != 'V', nil
}

func sortInts(xs []int) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j] < xs[j-1]; j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}

// mergeStackDepth lazily computes the canonical operand-stack depth for a
// merge pc from its StackMapFrame (0 when absent).
func (e *methodEmitter) mergeStackDepth(pc int) int {
	if e.mergeDepth == nil {
		e.mergeDepth = make(map[int]int)
	}
	if d, ok := e.mergeDepth[pc]; ok {
		return d
	}
	n := 0
	for i := range e.m.Code.StackMaps {
		fr := &e.m.Code.StackMaps[i]
		if int(fr.Offset) == pc {
			for _, it := range fr.Stack {
				n++
				if it.Tag == classfile.VItemLong || it.Tag == classfile.VItemDouble {
					n++
				}
			}
			break
		}
	}
	e.mergeDepth[pc] = n
	return n
}

// checkLocalIdx validates a local-slot index.
func (e *methodEmitter) checkLocalIdx(idx int) error {
	if idx < 0 || idx >= e.maxLocals {
		return fmt.Errorf("local[%d] out of range (%d locals)", idx, e.maxLocals)
	}
	return nil
}

// LocalName renders the variable for a local slot.
func (e *methodEmitter) LocalName(idx int) string { return e.localName(idx) }

// atypeName maps newarray component codes to descriptor letters.
func atypeName(at int) (string, bool) {
	switch at {
	case 4:
		return "Z", true
	case 5:
		return "C", true
	case 6:
		return "F", true
	case 7:
		return "D", true
	case 8:
		return "B", true
	case 9:
		return "S", true
	case 10:
		return "I", true
	case 11:
		return "J", true
	}
	return "", false
}

// joinStrings joins parts with sep.
func joinStrings(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

// smFrameAt finds the StackMapFrame at a given pc, or nil.
func (e *methodEmitter) smFrameAt(pc int) *classfile.StackMapFrame {
	for i := range e.m.Code.StackMaps {
		fr := &e.m.Code.StackMaps[i]
		if int(fr.Offset) == pc {
			return fr
		}
	}
	return nil
}
