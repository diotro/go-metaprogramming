package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	// Read the file as an AST, keeping track of line locations, etc
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "pkg/example/example_test.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if !strings.HasPrefix(fn.Name.Name, "Test") {
			return true
		}
		fmt.Printf("Found test function %s on line %d\n", fn.Name.Name, fset.Position(fn.Pos()).Line)

		testCases := make([]ast.Stmt, len(fn.Body.List))
		for i, statement := range fn.Body.List {
			expr, ok := statement.(*ast.ExprStmt)
			if !ok {
				return true
			}
			callExpr, ok := expr.X.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			exprIdentifier, ok := selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			if !(exprIdentifier.Name == "t" && selector.Sel.Name == "Run") {
				return true
			}
			fmt.Printf("Found t.Run on line %d\n", fset.Position(callExpr.Pos()).Line)

			testCases[i] = ast.Stmt(expr)
		}
		testCasesBlock := ast.BlockStmt{
			Lbrace: 0,
			List:   testCases,
			Rbrace: 0,
		}
		fn.Body = &testCasesBlock
		return true
	})
	// write new ast to file
	f, err := os.Create("pkg/out/example_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := printer.Fprint(f, fset, file); err != nil {
		log.Fatal(err)
	}
}
