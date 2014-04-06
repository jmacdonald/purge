package test

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Directory", func() {
	Describe("Size", func() {
		It("properly calculates the size of the sample directory", func() {
			// Set the expectedSize to the actual size
			// of the sample directory's contents.
			const expectedSize int64 = 512020

			dir, _ := os.Getwd()
			Expect(directory.Size(dir + "/sample")).To(Equal(expectedSize))
		})
	})

	Describe("Entries", func() {
		It("returns the correct number of entries", func() {
			dir, _ := os.Getwd()
			Expect(len(directory.Entries(dir + "/sample"))).To(Equal(3))
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
					Expect(entry.Size).To(Equal(directory.Size(dir + "/sample/" + entry.Name)))
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
