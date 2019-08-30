package examples

import (
	"go/ast"
	"go/token"
)

func test() {

	_ = 1 + (2 + 3)

	of, args := 3, 4

	functionCall(of, 3, args)

	_ = 3

	_ = "hello"


	threeAst := &ast.BasicLit{Value: "3"}

	helloAst := &ast.BasicLit{Value: "\"hello\""}


	expr1 := &ast.BinaryExpr{
		X:  &ast.BasicLit{Value: "2"},
		Op: token.ADD,
		Y:  &ast.BasicLit{Value: "3"},
	}

	expr2 := &ast.BinaryExpr{
		X:  &ast.BasicLit{Value: "1"},
		Op: token.ADD,
		Y:  expr1,
	}



	hello := "hello"

	six := 1 + (2 + 3)


	helloDec :=

	&ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "hello"}},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.BasicLit{Value: "\"hello\""}},
	}

	sixDec :=

	&ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "six"}},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{expr2},
	}



	// Prevent errors so screenshots are pretty
	_ = six
	_ = sixDec
	_ = hello
	_ = threeAst
	_ = helloAst
	_ = helloDec
	_ = expr2
}

func functionCall(of interface{}, i int, args interface{}) {

}
