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

	Describe("SetWorkingDirectory", func() {
		var error error

		// Change the working directory right before every test.
		JustBeforeEach(func() {
			error = navigator.SetWorkingDirectory(path)
		})

		Context("path is a valid directory", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/sample"
			})

			It("returns a nil error", func() {
				Expect(error).To(BeNil())
			})

			It("updates current path with its path argument", func() {
				Expect(navigator.CurrentPath()).To(Equal(path))
			})

			It("updates entries using path argument", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(path)))
			})

			It("resets selected index to zero", func() {
				navigator.SelectNextEntry()
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))

				navigator.SetWorkingDirectory(path)
				Expect(navigator.SelectedIndex()).To(BeZero())
			})
		})

		Context("path is a file", func() {
			original_path, _ := os.Getwd()

			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/sample/file"

				// Set the working directory to something valid
				// so that current path and entries are set.
				navigator.SetWorkingDirectory(original_path)

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(original_path))
			})

			It("does not update entries", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(original_path)))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})

		Context("path is invalid", func() {
			original_path, _ := os.Getwd()

			BeforeEach(func() {
				path = "/asdf"

				// Set the working directory to something valid
				// so that current path and entries are set.
				navigator.SetWorkingDirectory(original_path)

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(original_path))
			})

			It("does not update entries", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(original_path)))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})
	})

	Describe("SelectNextEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectNextEntry()
		})

		Context("directory has never been set", func() {
			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})
		})

		Context("directory has been set and has entries", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				navigator.SetWorkingDirectory(path)
			})

			It("increments the selected index by one", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})

			Context("last entry is selected", func() {
				var selectedIndex uint16

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for uint16(len(navigator.Entries()))-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}

					// Keep a reference to the last index.
					selectedIndex = navigator.SelectedIndex()
				})

				It("does not change the selected index", func() {
					Expect(navigator.SelectedIndex()).To(Equal(selectedIndex))
				})
			})
		})
	})

	Describe("SelectPreviousEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectPreviousEntry()
		})

		Context("directory has never been set", func() {
			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})
		})

		Context("directory has been set and has entries", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				navigator.SetWorkingDirectory(path)
			})

			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})

			Context("last entry is selected", func() {
				var selectedIndex uint16

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for uint16(len(navigator.Entries()))-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}

					// Keep a reference to the last index.
					selectedIndex = navigator.SelectedIndex()
				})

				It("decrements the selected index by one", func() {
					Expect(navigator.SelectedIndex()).To(BeEquivalentTo(selectedIndex-1))
				})
			})
		})
	})
})
