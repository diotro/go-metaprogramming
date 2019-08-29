package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fs := token.NewFileSet()
	f := fs.AddFile(filename, fs.Base(), len(b))
	var s scanner.Scanner
	s.Init(f, b, nil, scanner.ScanComments)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.SEMICOLON  {
			fmt.Println()
			continue
		}
		if lit == "" {
			fmt.Printf("%s ", tok)
		} else {
			fmt.Printf("[%s %s] ", tok, lit)
		}

		if tok == token.LBRACE {
			fmt.Println()
		}
	}
}
