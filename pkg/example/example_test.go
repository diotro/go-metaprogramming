package example_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/julian-zucker/go-ast-magic"
)

func TestAdd(t *testing.T) {
	t.Run("2 + 2 = 4", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(add(2, 2)).To(Equal(4))
	})
}

