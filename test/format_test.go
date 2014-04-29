package test

import (
	"github.com/jmacdonald/purge/view"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Format", func() {
	var output string
	var input int64

	Describe("Size", func() {
		JustBeforeEach(func() {
			output = view.Size(input)
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
