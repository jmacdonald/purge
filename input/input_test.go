package input

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInput(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Input Suite")
}

// Alias the string type to create a fake data source
type Reader string

// Implement the Reader interface, as required by input.go
func (s Reader) Read(target []byte) (int, error) {
	return copy(target, s), nil
}

var _ = Describe("Input", func() {
	Describe("Read", func() {
		var data Reader
		var result rune

		JustBeforeEach(func() {
			result = Read(data)
		})

		Context("data is a single byte character", func() {
			BeforeEach(func() {
				data = "j"
			})

			It("returns the complete data", func() {
				Expect(result).To(BeEquivalentTo(data))
			})
		})
	})
})
