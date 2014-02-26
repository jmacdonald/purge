package test

import (
	"github.com/jmacdonald/liberator/test/double"
	"github.com/jmacdonald/liberator/view"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("View", func() {
	Describe("Render", func() {
		var viewer *double.Viewer
		var renderer *double.Renderer

		BeforeEach(func() {
			viewer = &double.Viewer{Left: "left", Right: "right", Highlight: false}
			renderer = new(double.Renderer)
		})

		JustBeforeEach(func() {
			view.Render(viewer, renderer)
		})

		It("calls Clear() on output", func() {
			Expect(renderer.ClearCalled).To(BeTrue())
		})

		It("calls Flush() on output", func() {
			Expect(renderer.FlushCalled).To(BeTrue())
		})
	})
})
