package test

import (
	"fmt"
	"github.com/jmacdonald/purge/view"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("View", func() {
	Describe("FormatRow", func() {
		var result, errorMessage string
		var err error
		var row view.Row
		var size int

		JustBeforeEach(func() {
			errorMessage = fmt.Sprintf("view: formatting row to a size of %d"+
				" with '%s' and '%s' values is impossible", size, row.Left, row.Right)

			result, err = view.FormatRow(row, size)
		})

		Context("row values are set and size is larger than their sum", func() {
			BeforeEach(func() {
				row = view.Row{Left: "left", Right: "right"}
				size = 10
			})

			It("formats the row properly", func() {
				Expect(result).To(Equal("left right"))
			})
		})

		Context("row values are set and size is equal to their sum", func() {
			BeforeEach(func() {
				row = view.Row{Left: "left", Right: "right"}
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
				row = view.Row{Left: "left", Right: "right"}
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
				row = view.Row{Right: "right"}
				size = 10
			})

			It("formats the row properly", func() {
				Expect(result).To(Equal("     right"))
			})
		})
	})
})
