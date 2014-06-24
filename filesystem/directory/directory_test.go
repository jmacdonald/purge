package directory

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestDirectory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Directory Suite")
}

var _ = Describe("Directory", func() {
	Describe("Size", func() {
		var result chan *EntrySize
		var index int

		Context("when passed a directory path and an index", func() {
			BeforeEach(func() {
				result = make(chan *EntrySize)
				dir, _ := os.Getwd()
				index = 4

				go Size(dir+"/sample", index, result)
			})

			It("calculates the size of the directory", func(done Done) {
				// Set the expectedSize to the actual size
				// of the sample directory's contents.
				const expectedSize int64 = 512026

				entrySize := <-result
				Expect(entrySize.Size).To(Equal(expectedSize))
				close(done)
			})

			It("returns the index", func(done Done) {
				entrySize := <-result
				Expect(entrySize.Index).To(Equal(index))
				close(done)
			})
		})
	})
})

func contains(entries []*Entry, value string) bool {
	for _, entry := range entries {
		if entry.Name == value {
			return true
		}
	}
	return false
}
