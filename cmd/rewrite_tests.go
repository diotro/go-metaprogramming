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

	for declIndex, topLevelDecl := range file.Decls {
		fn, ok := topLevelDecl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		fmt.Printf("Found test function %s on line %d\n", fn.Name.Name, fset.Position(fn.Pos()).Line)

		testCases := make([]ast.Stmt, len(fn.Body.List))
		for i, statement := range fn.Body.List {
			expr, ok := statement.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := expr.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selector, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			exprIdentifier, ok := selector.X.(*ast.Ident)
			if !ok {
				continue
			}
			if !(exprIdentifier.Name == "t" && selector.Sel.Name == "Run") {
				continue
			}
			fmt.Printf("Found t.Run on line %d\n", fset.Position(callExpr.Pos()).Line)

			testCases[i] = expr
		}

		var rewrittenDecl ast.Decl = &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{
						Name:    "_",
					}},

					Type: nil,
					Values: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "Describe"},
							Args: []ast.Expr{
								&ast.BasicLit{
									Value: "\"Add\"",
								},
							},
						},
					},
				},
			},
		}
		file.Decls[declIndex] = rewrittenDecl
		continue
	}

	fmt.Printf("%#v\n", file.Decls[1])

	if err := printer.Fprint(os.Stdout, fset, file); err != nil {
		log.Fatal(err)
	}

	// Write new ast to file
	//f, err := os.Create("pkg/out/example_test.go")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//if err := printer.Fprint(os.Stdout, fset, file); err != nil {
	//	log.Fatal(err)
	//}
}
