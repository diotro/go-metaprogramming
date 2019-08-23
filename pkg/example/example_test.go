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
