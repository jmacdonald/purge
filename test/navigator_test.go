package liberator_test

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Navigator", func() {
	var (
		navigator *directory.Navigator
		path      string
	)

	BeforeEach(func() {
		navigator = new(directory.Navigator)
	})

	Describe("ChangeDirectory", func() {
		BeforeEach(func() {
			path, _ = os.Getwd()
			navigator.ChangeDirectory(path)
		})

		It("updates CurrentPath with its path argument", func() {
			Expect(navigator.CurrentPath).To(Equal(path))
		})

		It("updates Entries using path argument", func() {
			Expect(navigator.Entries).To(Equal(directory.Entries(path)))
		})
	})

	Describe("SelectNextEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectNextEntry()
		})

		Context("directory has never been set", func() {
			It("does not increment the selected index", func() {
				Expect(navigator.SelectedIndex).To(BeZero())
			})
		})

		Context("directory has been set and has entries", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				navigator.ChangeDirectory(path)
			})

			It("increments the selected index by one", func() {
				Expect(navigator.SelectedIndex).To(BeEquivalentTo(1))
			})

			Context("last entry is selected", func() {
				var selectedIndex uint16

				BeforeEach(func() {
					for uint16(len(navigator.Entries))-navigator.SelectedIndex > 1 {
						navigator.SelectNextEntry()
					}
					selectedIndex = navigator.SelectedIndex
				})

				It("does not increment the selected index", func() {
					Expect(navigator.SelectedIndex).To(Equal(selectedIndex))
				})
			})
		})
	})
})
