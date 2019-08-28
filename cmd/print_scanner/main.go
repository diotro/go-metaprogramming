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
	fs := token.NewFileSet()

	for _, arg := range os.Args[1:] {
		b, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		f := fs.AddFile(arg, fs.Base(), len(b))
		var s scanner.Scanner
		s.Init(f, b, nil, scanner.ScanComments)

		tokens := make([]tokAndContent, 0)
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.SEMICOLON {
				fmt.Println()
				continue
			}
			tokens = append(tokens, tokAndContent{tok, lit})
			fmt.Printf("[%s %s] ", tok, lit)
		}

	}
}

type tokAndContent struct {
	tok token.Token
	content string
}

func (tc *tokAndContent) String() string {
	return fmt.Sprintf("(%s, %s)", tc.tok.String(), tc.content)
}
