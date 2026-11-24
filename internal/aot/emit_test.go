package aot

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"catty/internal/classfile"
	"catty/rt"
)

func insn(op uint8, pc, length int32) rt.Insn {
	return rt.Insn{Op: op, PC: pc, Len: length, Ref: -1}
}

// TestGenerateMonitorLowering checks that both engines' shared monitor helpers
// are emitted with the original instruction PC and that the throw transfer goes
// through the ordinary handler selection.
func TestGenerateMonitorLowering(t *testing.T) {
	m := &classfile.Method{
		Owner: "Fixture", Name: "m", Descriptor: "()V",
		Flags: rt.AccPublic | rt.AccStatic, MaxStack: 1, MaxLocals: 0, CodeLength: 5,
		Insns: []rt.Insn{
			insn(rt.OpAconstNull, 0, 1),
			insn(rt.OpMonitorenter, 1, 1),
			insn(rt.OpAconstNull, 2, 1),
			insn(rt.OpMonitorexit, 3, 1),
			insn(rt.OpReturn, 4, 1),
		},
		Handlers: []rt.Handler{{StartPC: 0, EndPC: 5, HandlerPC: 4, CatchType: ""}},
	}
	heights := map[int32]int{0: 0, 1: 1, 2: 0, 3: 1, 4: 0}
	shapes := map[int32][]int{0: nil, 1: []int{1}, 2: nil, 3: []int{1}, 4: nil}
	source, err := Generate("cattyMethod0", m, heights, shapes)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "generated.go", "package main\n\n"+source, 0)
	if err != nil {
		t.Fatalf("generated source does not parse: %v\n%s", err, source)
	}
	var fn *ast.FuncDecl
	for _, decl := range file.Decls {
		if candidate, ok := decl.(*ast.FuncDecl); ok && candidate.Name.Name == "cattyMethod0" {
			fn = candidate
		}
	}
	if fn == nil {
		t.Fatal("generated function missing")
	}

	type monitorCall struct {
		name string
		pc   int64
	}
	var calls []monitorCall
	handleThrows := 0
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch sel.Sel.Name {
		case "MonitorEnter", "MonitorExit":
			if len(call.Args) < 2 {
				t.Fatalf("%s call has %d arguments", sel.Sel.Name, len(call.Args))
			}
			literal, ok := call.Args[1].(*ast.BasicLit)
			if !ok || literal.Kind != token.INT {
				t.Fatalf("%s call does not carry a literal PC: %v", sel.Sel.Name, call.Args[1])
			}
			pc, err := strconv.ParseInt(literal.Value, 0, 64)
			if err != nil {
				t.Fatalf("bad PC literal: %v", err)
			}
			calls = append(calls, monitorCall{name: sel.Sel.Name, pc: pc})
		case "HandleThrow":
			handleThrows++
		}
		return true
	})
	if len(calls) != 2 {
		t.Fatalf("generated %d monitor helper calls, want 2", len(calls))
	}
	if calls[0].name != "MonitorEnter" || calls[0].pc != 1 {
		t.Fatalf("first monitor lowering = %+v, want MonitorEnter at pc 1", calls[0])
	}
	if calls[1].name != "MonitorExit" || calls[1].pc != 3 {
		t.Fatalf("second monitor lowering = %+v, want MonitorExit at pc 3", calls[1])
	}
	if handleThrows == 0 {
		t.Fatal("monitor lowering did not use the shared handler-transfer path")
	}
}
