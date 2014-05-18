package view

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestView(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "View Suite")
}

var _ = Describe("View", func() {
	Describe("FormatRow", func() {
		var result, errorMessage string
		var err error
		var row Row
		var size int

		JustBeforeEach(func() {
			errorMessage = fmt.Sprintf("view: formatting row to a size of %d"+
				" with '%s' and '%s' values is impossible", size, row.Left, row.Right)

			result, err = FormatRow(row, size)
		})

		Context("row values are set and size is larger than their sum", func() {
			BeforeEach(func() {
				row = Row{Left: "left", Right: "right"}
				size = 10
			})

			It("formats the row properly", func() {
				Expect(result).To(Equal("left right"))
			})
		})

		Context("row values are set and size is equal to their sum", func() {
			BeforeEach(func() {
				row = Row{Left: "left", Right: "right"}
				size = 9
			})

			It("returns an empty string", func() {
				Expect(result).To(Equal(""))
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})

			It("returns the proper error message", func() {
				Expect(err.Error()).To(Equal(errorMessage))
			})
		})

		Context("row values are set and size is smaller than their sum", func() {
			BeforeEach(func() {
				row = Row{Left: "left", Right: "right"}
				size = 5
			})

			It("returns an empty string", func() {
				Expect(result).To(Equal(""))
			})

			It("returns an error", func() {
				Expect(err.Error()).ToNot(BeNil())
			})

			It("returns the proper error message", func() {
				Expect(err.Error()).To(Equal(errorMessage))
			})
		})

		Context("one of the row values isn't set", func() {
			BeforeEach(func() {
				row = Row{Right: "right"}
				size = 10
			})

			It("formats the row properly", func() {
				Expect(result).To(Equal("     right"))
			})
		})
	})
})

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

		Context("When passed more than a megabyte but less than a gigabyte", func() {
			BeforeEach(func() {
				input = 1290000
			})

			It("returns the size in megabytes", func() {
				Expect(output).To(Equal("1.2 MB"))
			})
		})

		Context("When passed more than a gigabyte but less than a terabyte", func() {
			BeforeEach(func() {
				input = 1290000000
			})

			It("returns the size in gigabytes", func() {
				Expect(output).To(Equal("1.2 GB"))
			})
		})

		Context("When passed more than a terabyte", func() {
			BeforeEach(func() {
				input = 1350000000000
			})

			It("returns the size in terabytes", func() {
				Expect(output).To(Equal("1.2 TB"))
			})
		})

		Context("When passed more than a petabyte", func() {
			BeforeEach(func() {
				input = 1360000000000000
			})

			It("still returns the size in terabytes", func() {
				Expect(output).To(Equal("1236.9 TB"))
			})
		})
	})
})
