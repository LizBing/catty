package emitter

import (
	"fmt"
	"os"

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
	e.reach = reach

	jumpTargets, err := verify.BranchTargets(e.code)
	if err != nil {
		return err
	}
	for _, hp := range handlerPCs {
		jumpTargets[hp] = true
	}
	e.targets = jumpTargets

	canonDepths := e.canonDepths

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

	debugOn := os.Getenv("CATTY_DEPTH") != ""

	for _, pc := range pcs {
		if debugOn && true && e.m.Name == "main" {
			fmt.Fprintf(os.Stderr, "[dep] pc=%d canon=%d before\n", pc, func() int {
				if d, ok := canonDepths[pc]; ok {
					return d
				}
				return -999
			}())
		}
		isJumpTarget := jumpTargets[pc]
		_, isHandler := e.handlerAt[pc]
		if d, ok := canonDepths[pc]; ok {
			e.depth = d
		}
		if isHandler {
			e.p("L%d:", pc)
			e.p("s%d = exc.Obj", e.depth-1)
		} else if isJumpTarget {
			e.p("L%d:", pc)
		}
		before := e.w.Len()
		if os.Getenv("CATTY_PC_TRACE") != "" && pc == 87 {
			fmt.Fprintf(os.Stderr, "[pc87] depth=%d op=%#x\n", e.depth, e.code[pc])
		}
		if err := e.emitOne(pc); err != nil {
			return err
		}
		if os.Getenv("CATTY_PC_TRACE") != "" {
			emitted := e.w.String()[before:]
			if true {
				fmt.Fprintf(os.Stderr, "[pc] %d → %q\n", pc, emitted)
			}
		}
	}
	return nil
}
