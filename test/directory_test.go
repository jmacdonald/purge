package liberator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jmacdonald/liberator/filesystem/directory"
	"os"
)

var _ = Describe("Directory", func() {
	Describe("Size", func() {
		It("properly calculates the size of the sample directory", func() {
			// Set the expectedSize to the actual size
			// of the sample directory's contents.
			const expectedSize int64 = 512020

			dir, _ := os.Getwd()
			Expect(directory.Size(dir + "/sample")).To(Equal(expectedSize));
		})
	})
})
