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
		if _, isHandler := e.handlerAt[pc]; isHandler {
			e.p("L%d:", pc)
			d := e.mergeStackDepth(pc)
			e.depth = d
			e.p("s%d = exc.Obj", d-1)
		} else if jumpTargets[pc] {
			e.p("L%d:", pc)
		}
		if err := e.emitOne(pc); err != nil {
			return err
		}
	}
	return nil
}
