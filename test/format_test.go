package liberator_test

import (
	. "github.com/jmacdonald/liberator/format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Format", func() {
	var output string
	var input int64

	Describe("Size", func() {
		JustBeforeEach(func() {
			output = Size(input)
		})

		Context("When passed less than a kilobyte", func() {
			BeforeEach(func() {
				input = 512
			})

			It("returns the size in bytes", func() {
				Expect(output).To(Equal("512 bytes"))
			})
		})

		Context("When passed more than a kilobyte but less than a megabyte", func() {
			BeforeEach(func() {
				input = 2900
			})

			It("returns the size in kilobytes", func() {
				Expect(output).To(Equal("2.8 KB"))
			})
		})
	})
})
