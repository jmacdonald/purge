package input

import (
	"github.com/jmacdonald/purge/test/double"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInput(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Input Suite")
}

var _ = Describe("Input", func() {
	Describe("Read", func() {
		var data double.Reader
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
