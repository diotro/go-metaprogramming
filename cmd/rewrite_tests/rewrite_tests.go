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
		transformedTestCasesBlockStmt := &ast.BlockStmt{
			List: transformTestCases(fn.Body.List),
		}

		testNameString := &ast.BasicLit{
			// Need to wrap in "" to turn from identifier to string. Also, remove Test.
			Value: "\"" + strings.Replace(fn.Name.Name, "Test", "", 1) + "\"",
		}

		// This is the call to Describe(name, …)
		describeCall := &ast.CallExpr{
			Fun: &ast.Ident{Name: "Describe"},
			Args: []ast.Expr{
				testNameString,
				&ast.FuncLit{
					Type: thunkFunctionType,
					Body: transformedTestCasesBlockStmt,
				},
			},
		}

		// Then, rewrite the declaration to use the new Ginkgo tests,
		// and wrap them in a "var _ = Describe(...)"
		var rewrittenDecl ast.Decl = &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{
						Name: "_",
					}},
					Type: nil,
					Values: []ast.Expr{
						describeCall,
					},
				},
			},
		}
		file.Decls[declIndex] = rewrittenDecl
		continue
	}

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

// Transforms any gotest-style tests in the given set of statements
// to Ginkgo tests.
func transformTestCases(stmts []ast.Stmt) []ast.Stmt {
	transformedTestCases := make([]ast.Stmt, len(stmts))

	for i, statement := range stmts {
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

			removed, ok := removeGDot(stmt)
			if ok {
				// If we could remove the `g.`, add the transformed statement
				transformedStmts = append(transformedStmts, removed)
			} else {
				// The removal didn't work, this must be some other shape
				// of statement, so keep it around
				transformedStmts = append(transformedStmts, stmt)
			}
		}

		transformedTestCases[i] = &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{Name: "It"},
				Args: []ast.Expr{
					testCaseName,
					&ast.FuncLit{
						Type: thunkFunctionType,
						Body: &ast.BlockStmt{List: transformedStmts},
					},
				},
			},
		}
	}

	return transformedTestCases
}


// The type of a function that takes no inputs and produces no outputs.
// Really more of a "unit-returning thunk", but this name is fine.
var thunkFunctionType = &ast.FuncType{
	Params: &ast.FieldList{
		List: []*ast.Field{},
	},
	Results: &ast.FieldList{
		List: []*ast.Field{},
	},
}

// Is this statement `g := NewGomegaWithT(t)`?
func isGEqualNewGomega(stmt ast.Stmt) bool {
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

	rhsFunc, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	rhsFuncIdent, ok := rhsFunc.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if rhsFuncIdent.Name != "NewGomegaWithT" {
		return false
	}
	// Could check for rhsFunc's args being `t`, but what other
	// testing object do you have access to?
	return true
}

// Transforms `g.Expect(foo).To(Equal(bar))
// into `Expect(foo).To(Equal(bar))
func removeGDot(stmt ast.Stmt) (ast.Stmt, bool) {
	expr, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}
	call, ok := expr.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	funSelector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}
	// The identifier at the base of the expression, namely
	// the `g` in `g.Expect`
	funSelectorExprAsCall, ok := funSelector.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	funBaseSelector, ok := funSelectorExprAsCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}
	funBaseIdent, ok := funBaseSelector.X.(*ast.Ident)
	if !ok {
		return nil, false
	}
	if funBaseIdent.Name != "g" {
		return nil, false
	}

	expect := &ast.CallExpr{
		Fun: &ast.Ident{
			Name: "Expect",
		},
		Args: funSelectorExprAsCall.Args,
	}

	// This is the "Expect" in `g.Expect`, so we can construct a new
	// function call with just that selector
	//selName := funSelector.Sel.Name

	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   expect,
				Sel: &ast.Ident{Name: "To"},
			},
			Args: call.Args,
		},
	}, true
}
