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
	default:
		fail("unknown source check %q", check)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
