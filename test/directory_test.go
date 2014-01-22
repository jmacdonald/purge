package liberator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/jmacdonald/liberator/filesystem/directory"
	"os"
)

var _ = Describe("Directory", func() {
	Describe("Size", func() {
		It("properly calculates the size of the sample directory", func() {
			// Set the expectedSize to the actual size
			// of the sample directory's contents.
			const expectedSize int64 = 512020

			dir, _ := os.Getwd()
			Expect(Size(dir + "/sample")).To(Equal(expectedSize));
		})
	})

	Describe("Entries", func() {
		It("returns the correct number of entries", func() {
			dir, _ := os.Getwd()
			Expect(len(Entries(dir + "/sample"))).To(Equal(2));
		})

		It("returns the proper names", func() {
			dir, _ := os.Getwd()
			entries := Entries(dir + "/sample")
			Expect(contains(entries, "directory")).To(BeTrue())
			Expect(contains(entries, "file")).To(BeTrue())
		})

		It("returns the proper sizes", func() {
			dir, _ := os.Getwd()
			entries := Entries(dir + "/sample")
			for _, entry := range entries {
				Expect(entry.Size).To(Equal(Size(dir + "/sample/" + entry.Name)))
			}
		})
	})
})

func contains(entries []*Entry, value string) bool {
	for _, entry := range entries {
		if entry.Name == value { return true }
	}
	return false
}
