package test

import (
	"github.com/jmacdonald/purge/input"
	"github.com/jmacdonald/purge/test/double"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Input", func() {
	Describe("Read", func() {
		var data double.Reader
		var result rune

		JustBeforeEach(func() {
			result = input.Read(data)
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
