package test

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
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

	Describe("Entries", func() {
		It("returns the correct number of entries", func() {
			dir, _ := os.Getwd()
			Expect(len(directory.Entries(dir + "/sample"))).To(Equal(4))
		})

		It("returns the proper names", func() {
			dir, _ := os.Getwd()
			entries := directory.Entries(dir + "/sample")
			Expect(contains(entries, "directory")).To(BeTrue())
			Expect(contains(entries, "file")).To(BeTrue())
		})

		It("returns the proper sizes", func() {
			dir, _ := os.Getwd()
			entries := directory.Entries(dir + "/sample")
			for _, entry := range entries {
				entryInfo, _ := os.Stat(dir + "/sample/" + entry.Name)

				if entryInfo.IsDir() {
					result := make(chan *directory.EntrySize)
					go directory.Size(dir+"/sample/"+entry.Name, 0, result)
					Expect(entry.Size).To(Equal((<-result).Size))
				} else {
					Expect(entry.Size).To(Equal(entryInfo.Size()))
				}
			}
		})

		It("returns the proper directory statuses", func() {
			dir, _ := os.Getwd()
			entries := directory.Entries(dir + "/sample")
			for _, entry := range entries {
				fileInfo, _ := os.Stat(dir + "/sample/" + entry.Name)
				Expect(entry.IsDirectory).To(Equal(fileInfo.IsDir()))
			}
		})

		It("returns the entries sorted by size (largest to smallest)", func() {
			dir, _ := os.Getwd()
			entries := directory.Entries(dir + "/sample")

			// Collect the entry names.
			entry_names := make([]string, len(entries), len(entries))
			for index, entry := range entries {
				entry_names[index] = entry.Name
			}

			Expect(entry_names).To(Equal([]string{"directory", "file", "small_file", "empty_file"}))
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
