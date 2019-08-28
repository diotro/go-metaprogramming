package example

import (
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Add", func() {
	Specify("2 + 2 = 4", func() {
		Expect(Add(2, 2)).To(Equal(4))
	})
	Specify("0 + 0 = 0", func() {
		Expect(Add(0, 0)).To(Equal(0))
	})
})

func Test(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "TODO FILL THIS IN Suite")
}

