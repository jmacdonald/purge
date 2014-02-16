package test

import (
	"github.com/jmacdonald/liberator/input"
	"github.com/jmacdonald/liberator/test/double"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Input", func() {
	Describe("Read", func() {
		var data double.Reader
		var navigator *double.Navigator

		JustBeforeEach(func() {
			navigator = new(double.Navigator)
			input.Read(data, navigator)
		})

		Context("input is a 'j'", func() {
			BeforeEach(func() {
				data = "j"
			})

			It("calls SelectNextEntry() on navigator", func() {
				Expect(navigator.SelectNextEntryCalled).To(BeTrue())
			})
		})

		Context("input is a 'k'", func() {
			BeforeEach(func() {
				data = "k"
			})

			It("calls SelectPreviousEntry() on navigator", func() {
				Expect(navigator.SelectPreviousEntryCalled).To(BeTrue())
			})
		})

		Context("input is a carriage return", func() {
			BeforeEach(func() {
				data = "\n"
			})

			It("calls IntoSelectedEntry() on navigator", func() {
				Expect(navigator.IntoSelectedEntryCalled).To(BeTrue())
			})
		})

		Context("input is an 'h'", func() {
			BeforeEach(func() {
				data = "h"
			})

			It("calls ToParentDirectory() on navigator", func() {
				Expect(navigator.ToParentDirectoryCalled).To(BeTrue())
			})
		})
	})
})
