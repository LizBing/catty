package emitter

import (
	"fmt"
	"os"
	"strings"

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
	if os.Getenv("CATTY_BODYDBG") != "" && e.cf.ThisClass == "com/eclipsesource/json/Json" &&
		e.m.Desc == "(F)Lcom/eclipsesource/json/JsonValue;" {
		fmt.Fprintf(os.Stderr, "[body-map] len=%d v1=%v ptr=%p\n", len(canonDepths), canonDepths[1], canonDepths)
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
		if d, ok := canonDepths[pc]; ok {
			e.depth = d
			if os.Getenv("CATTY_BODYDBG") != "" && e.cf.ThisClass == "com/eclipsesource/json/Json" &&
				e.m.Desc == "(F)Lcom/eclipsesource/json/JsonValue;" {
				fmt.Fprintf(os.Stderr, "[bd] pc=%d canon=%d depth_now=%d\n", pc, d, e.depth)
			}
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
		if os.Getenv("CATTY_NEG") != "" {
			emitted := e.w.String()[before:]
			if strings.Contains(emitted, "s-") {
				fmt.Fprintf(os.Stderr, "[neg] %s.%s pc=%d depth=%d out=%q\n",
					e.cf.ThisClass, e.m.Name, pc, e.depth, emitted)
			}
		}
		if os.Getenv("CATTY_PC_TRACE") != "" {
			emitted := e.w.String()[before:]
			if true {
				fmt.Fprintf(os.Stderr, "[pc] %d → %q\n", pc, emitted)
			}
		}
	}
	e.p("return nil, nil // unreachable terminal")
	return nil
}
