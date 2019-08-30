package main

import (
	"fmt"
	"go/ast"
	"go/parser"
)

func main() {
	exprStr := `
func () {
	3 + (2 + 3)
}
`
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		fmt.Print(err)
	}

	inner := expr.(*ast.FuncLit).Body.List[0]

	fmt.Printf("\n\n%#v\n", inner.(*ast.ExprStmt).X.(*ast.BinaryExpr))
}
