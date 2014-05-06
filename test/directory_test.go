package test

import (
	"github.com/jmacdonald/purge/filesystem/directory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Directory", func() {
	Describe("Size", func() {
		var result chan *directory.EntrySize
		var index int

		Context("when passed a directory path and an index", func() {
			BeforeEach(func() {
				result = make(chan *directory.EntrySize)
				dir, _ := os.Getwd()
				index = 4

				go directory.Size(dir+"/sample", index, result)
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

func contains(entries []*directory.Entry, value string) bool {
	for _, entry := range entries {
		if entry.Name == value {
			return true
		}
	}
	return false
}
