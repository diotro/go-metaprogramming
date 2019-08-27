package main

import (
	"fmt"
	"go/ast"
	"go/parser"
)

func main() {
	exprStr := `
func () {
	g := NewGomegaWithT(t)
}
`
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		fmt.Print(err)
	}

	inner := expr.(*ast.FuncLit).Body.List[0]

	fmt.Printf("%#v\n", inner)
}
