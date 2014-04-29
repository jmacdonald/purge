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

	Describe("Map", func() {
		var character rune
		var navigator *double.Navigator
		var result bool

		JustBeforeEach(func() {
			navigator = new(double.Navigator)
			result = input.Map(character, navigator)
		})

		Context("input is a 'j'", func() {
			BeforeEach(func() {
				character = 'j'
			})

			It("calls SelectNextEntry() on navigator", func() {
				Expect(navigator.SelectNextEntryCalled).To(BeTrue())
			})
		})

		Context("input is a 'k'", func() {
			BeforeEach(func() {
				character = 'k'
			})

			It("calls SelectPreviousEntry() on navigator", func() {
				Expect(navigator.SelectPreviousEntryCalled).To(BeTrue())
			})
		})

		Context("input is a carriage return", func() {
			BeforeEach(func() {
				character = '\r'
			})

			It("calls IntoSelectedEntry() on navigator", func() {
				Expect(navigator.IntoSelectedEntryCalled).To(BeTrue())
			})
		})

		Context("input is an 'h'", func() {
			BeforeEach(func() {
				character = 'h'
			})

			It("calls ToParentDirectory() on navigator", func() {
				Expect(navigator.ToParentDirectoryCalled).To(BeTrue())
			})
		})

		Context("input is an 'x'", func() {
			BeforeEach(func() {
				character = 'x'
			})

			It("calls RemoveSelectedEntry() on navigator", func() {
				Expect(navigator.RemoveSelectedEntryCalled).To(BeTrue())
			})
		})

		Context("input is a 'q'", func() {
			BeforeEach(func() {
				character = 'q'
			})

			It("returns true", func() {
				Expect(result).To(BeTrue())
			})
		})
	})
})
