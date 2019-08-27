Playing around with `go/ast`, with the aim of automating the transformation of `testing` tests to `ginkgo` tests.

E.x. to transform

```go
package example

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAdd(t *testing.T) {
	t.Run("2 + 2 = 4", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(add(2, 2)).To(Equal(4))
	})

	t.Run("0 + 0 = 0", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(add(0, 0)).To(Equal(0))
	})
}
```
into

```go
package example

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add", func() {
	It("2 + 2 should be 4", func() {
		Expect(add(2, 2)).To(Equal(4))
	})

	It("0 + 0 should be 0", func() {
		Expect(add(2, 2)).To(Equal(4))
	})
})
```
package main

import (
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

// Mostly just for testing

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, os.Args[0], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	err = printer.Fprint(os.Stdout, fset, file)
	if err != nil {
		log.Fatal(err)
	}
}
