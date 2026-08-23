package emitter

import (
	"catty/internal/classfile"
	"catty/internal/verify"
)

type handled int

const (
	handledFallthrough handled = iota
	handledTerminal
)

// body emits the linear translation of every REACHABLE instruction.
func (e *methodEmitter) body() error {
	var handlerPCs []int
	for _, h := range e.handlers {
		handlerPCs = append(handlerPCs, int(h.HandlerPc))
	}

	reach, err := verify.Reachable(e.code, handlerPCs)
	if err != nil {
		return err
	}

	jumpTargets, err := verify.BranchTargets(e.code)
	if err != nil {
		return err
	}
	for _, hp := range handlerPCs {
		jumpTargets[hp] = true
	}
	e.targets = jumpTargets
	e.handlerAt = make(map[int]string)
	for _, h := range e.handlers {
		e.handlerAt[int(h.HandlerPc)] = catchName(e.cf, h.CatchType)
	}
	e.handlerAt = make(map[int]string)
	for _, h := range e.handlers {
		e.handlerAt[int(h.HandlerPc)] = catchName(e.cf, h.CatchType)
	}

	argDescs, _, err := splitMethodDesc(e.m.Desc)
	if err != nil {
		return err
	}
	slot := 0
	if e.m.AccessFlags&classfile.AccStatic == 0 {
		e.p("l0 = recv")
		slot = 1
	}
	for i, a := range argDescs {
		e.p("l%d = args[%d]", slot, i)
		slot += descSlots(a)
	}

	pcs := make([]int, 0, len(reach))
	for pc := range reach {
		pcs = append(pcs, pc)
	}
	sortInts(pcs)

	for _, pc := range pcs {
		isJumpTarget := jumpTargets[pc]
		_, isHandler := e.handlerAt[pc]
		if isHandler {
			// Exception handler entry (JVMS §2.6, §4.10.1.6): the operand
			// stack is cleared and the caught exception pushed, so the
			// canonical depth is always 1 — no StackMapFrame lookup needed.
			// The push must come AFTER the label: every arrival goes through
			// excDispatch's `goto`, so a statement emitted before the label
			// is dead code and the handler would observe a stale slot.
			e.p("L%d:", pc)
			e.p("s0 = exc.Obj")
			e.depth = 1
		} else if isJumpTarget {
			// Merge point: reset operand stack to the StackMapTable
			// canonical depth so subsequent slots are named consistently
			// with the verifier's frame.
			e.depth = e.mergeStackDepth(pc)
			e.p("L%d:", pc)
		}
		if err := e.emitOne(pc); err != nil {
			return err
		}
	}
	return nil
}
