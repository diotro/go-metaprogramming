package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func main() {
	src, err := ioutil.ReadFile("src/example/example_test.go")
	if err != nil {
		panic("Could not read example.go")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}
