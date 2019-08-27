package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "pkg/example/example_ginkgo.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}
	for _, decl := range file.Decls {
		fmt.Printf("%#v\n", decl)

		genDecl, ok := decl.(*ast.GenDecl)
		if ok {
			for _, spec := range genDecl.Specs {
				fmt.Printf("%#v\n", spec)
				valueSpec, ok := spec.(*ast.ValueSpec)
				if ok {
					fmt.Printf("%#v\n", valueSpec.Names[0].Obj)
				}
			}
		}
	}

	err = printer.Fprint(os.Stdout, fset, file)
	if err != nil {
		log.Fatal(err)
	}

}
