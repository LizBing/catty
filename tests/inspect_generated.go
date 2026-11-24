//go:build ignore

// inspect_generated checks the shape of two functions emitted for the fixture.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fail("usage: inspect_generated SYMBOL mutation|control-flow FILE")
	}
	symbol, check, file := os.Args[1], os.Args[2], os.Args[3]
	positions := token.NewFileSet()
	parsed, err := parser.ParseFile(positions, file, nil, 0)
	if err != nil {
		fail("%v", err)
	}

	for _, decl := range parsed.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != symbol || fn.Body == nil {
			continue
		}
		inspect(fn, check)
		fmt.Print(file, "\n")
		if err := printer.Fprint(os.Stdout, positions, fn); err != nil {
			fail("%v", err)
		}
		fmt.Println()
		return
	}
	fail("%s: missing Go function %s", file, symbol)
}

func inspect(fn *ast.FuncDecl, check string) {
	var calls, mutations, branches int
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CallExpr:
			calls++
		case *ast.AssignStmt, *ast.IncDecStmt:
			mutations++
		case *ast.BinaryExpr:
			if node.Op == token.ADD {
				mutations++
			}
		case *ast.IfStmt, *ast.SwitchStmt:
			branches++
		case *ast.BranchStmt:
			if node.Tok == token.GOTO {
				branches++
			}
		}
		return true
	})
	switch check {
	case "mutation":
		if mutations == 0 && calls < 2 {
			fail("%s: missing the lowered field update", fn.Name)
		}
	case "control-flow":
		if branches == 0 || calls == 0 {
			fail("%s: missing lowered control flow or runtime operations", fn.Name)
		}
	case "monitors":
		checkMonitorLowering(fn)
	default:
		fail("unknown source check %q", check)
	}
}

// checkMonitorLowering verifies that every monitor helper call carries the
// original instruction PC and that the same PC feeds the shared exception
// transfer when the method has handlers.
func checkMonitorLowering(fn *ast.FuncDecl) {
	var monitorPCs []int64
	handled := map[int64]bool{}
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
				fail("%s: %s call has %d arguments", fn.Name.Name, sel.Sel.Name, len(call.Args))
			}
			pc, ok := intLiteral(call.Args[1])
			if !ok {
				fail("%s: %s call does not carry an original PC literal", fn.Name.Name, sel.Sel.Name)
			}
			monitorPCs = append(monitorPCs, pc)
		case "HandleThrow":
			if len(call.Args) >= 2 {
				if pc, ok := intLiteral(call.Args[1]); ok {
					handled[pc] = true
				}
			}
		}
		return true
	})
	if len(monitorPCs) == 0 {
		fail("%s: missing monitor helper lowering", fn.Name.Name)
	}
	for _, pc := range monitorPCs {
		if !handled[pc] {
			fail("%s: monitor helper at pc %d does not feed the shared exception transfer", fn.Name.Name, pc)
		}
	}
}

func intLiteral(expr ast.Expr) (int64, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	value, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
