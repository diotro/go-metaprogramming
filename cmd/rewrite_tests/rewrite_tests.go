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

		// first: extract each test case and transform it to Ginkgo style
		transformedTestCases := make([]ast.Stmt, len(fn.Body.List))
		for i, statement := range fn.Body.List {
			expr, ok := statement.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := expr.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			// Make sure that this has shape `t.Run`
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

			// Extract the test case name
			testCaseName := callExpr.Args[0]
			testCaseFuncStmt := callExpr.Args[1]

			testCaseFunc, ok := testCaseFuncStmt.(*ast.FuncLit)
			if !ok {
				continue
			}
			testCaseFuncStmts := testCaseFunc.Body.List
			transformedStmts := make([]ast.Stmt, 0)
			for _, stmt := range testCaseFuncStmts {
				// Remove `g := NewGomegaWithT(t)`
				if isGEqualNewGomega(stmt) {
					continue
				}
				transformedStmts = append(transformedStmts, stmt)
			}
			transformedTestCases[i] = &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "It"},
					Args: []ast.Expr{
						testCaseName,
						&ast.FuncLit{
							Type: &ast.FuncType{
								Params: &ast.FieldList{
									List: []*ast.Field{},
								},
								Results: &ast.FieldList{
									List: []*ast.Field{},
								},
							},
							Body: &ast.BlockStmt{List: transformedStmts},
						},
					},
				},
			}
		}

		transformedTestCasesBlockStmt := &ast.BlockStmt{
			List: transformedTestCases,
		}

		// Then, rewrite the declaration to use the new Ginkgo tests,
		// and wrap them in a "var _ = Describe(...)'
		var rewrittenDecl ast.Decl = &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{
						Name: "_",
					}},

					Type: nil,
					Values: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "Describe"},
							Args: []ast.Expr{
								&ast.BasicLit{
									Value: "\"Add\"",
								},
								&ast.FuncLit{
									Type: &ast.FuncType{
										Params: &ast.FieldList{
											List: []*ast.Field{},
										},
										Results: &ast.FieldList{
											List: []*ast.Field{},
										},
									},
									Body: transformedTestCasesBlockStmt,
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

func isGEqualNewGomega(stmt ast.Stmt) bool {
	fmt.Printf("%#v\n", stmt)
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok || assignStmt.Tok != token.DEFINE {
		return false
	}

	if len(assignStmt.Lhs) != 1 {
		return false
	}
	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if lhsIdent.Name != "g" {
		return false
	}

	fmt.Println("lhs")

	rhsFunc, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	rhsFuncIdent, ok := rhsFunc.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	fmt.Println("rhs")
	if rhsFuncIdent.Name != "NewGomegaWithT" {
		return false
	}
	fmt.Println("true")
	// Could check for rhsFunc's args being `t`, but what other
	// testing object do you have access to?
	return true
}
