Playing around with `go/ast`, with the aim of automating the transformation of `testing` tests to `ginkgo` tests.

E.x. to transform

```go
func TestAdd(t *testing.T) {
	t.Run("2 + 2 should be 4", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(add(2, 2)).To(Equal(4))
	})

	t.Run("0 + 0 should be 0", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(add(0, 0)).To(Equal(0))
	})
}
```
into

```go
var _ = Describe("Add", func() {
    It("2 + 2 should be 4", func() {
        Expect(add(2, 2)).To(Equal(4))
    })
    
    It("0 + 0 should be 0", func() {
        Expect(add(2, 2)).To(Equal(4))
    })
})
```
