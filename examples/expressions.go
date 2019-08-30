package examples

import (
	"go/ast"
	"go/token"
)

func test() {

	_ := 1 + (2 + 3)

	of,args := 3, 4

	functionCall(of, 3, args)


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

	_ = expr2
}

func functionCall(of interface{}, i int, args interface{}) {

}
