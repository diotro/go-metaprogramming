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
