package directory

import (
	"github.com/jmacdonald/purge/view"
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

var _ = Describe("Navigator", func() {
	var (
		navigator    *Navigator
		path         string
		error        error
		originalPath string
		viewBuffer   chan<- *view.Buffer
	)

	BeforeEach(func() {
		originalPath, _ = os.Getwd()
		navigator = new(Navigator)
		viewBuffer = make(chan<- *view.Buffer, 10)
		navigator.view = viewBuffer
		navigator.SetWorkingDirectory(originalPath)
	})

	Describe("SetWorkingDirectory", func() {
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
				Expect(len(navigator.Entries())).To(Equal(4))
			})

			It("resets selected index to zero", func() {
				navigator.SelectNextEntry()
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))

				navigator.SetWorkingDirectory(path)
				Expect(navigator.SelectedIndex()).To(BeZero())
			})

			It("resets previous view data indices", func() {
				_ = navigator.View(1)

				navigator.SetWorkingDirectory(path)
				Expect(navigator.ViewDataIndices()).To(Equal([2]int{0, 0}))
			})
		})

		Context("path is a file", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/sample/file"

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(originalPath))
			})

			It("does not update entries", func() {
				Expect(len(navigator.Entries())).To(Equal(4))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})

		Context("path is invalid", func() {
			BeforeEach(func() {
				path = "/asdf"

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(originalPath))
			})

			It("does not update entries", func() {
				Expect(len(navigator.Entries())).To(Equal(4))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})

		Context("path has a trailing slash", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/"
			})

			It("strips the trailing slash", func() {
				Expect(navigator.CurrentPath()).To(Equal(path[:len(path)-1]))
			})
		})
	})

	Describe("SelectedEntry", func() {
		Context("the second entry is selected", func() {
			BeforeEach(func() {
				navigator.SelectNextEntry()
			})

			It("returns the entry at the currently selected index", func() {
				Expect(navigator.SelectedEntry()).To(Equal(navigator.Entries()[navigator.SelectedIndex()]))
			})
		})

		Context("current directory is empty", func() {
			BeforeEach(func() {
				directory_name := "new_directory"
				os.Mkdir(directory_name, 0700)

				// Navigate into the empty directory
				navigator.SetWorkingDirectory(originalPath + "/" + directory_name)
			})

			It("returns nil", func() {
				Expect(navigator.SelectedEntry()).To(BeNil())
			})
		})
	})

	Describe("SelectFirstEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectFirstEntry()
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
				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for len(navigator.Entries())-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}
				})

				It("resets the selected index to zero", func() {
					Expect(navigator.SelectedIndex()).To(BeZero())
				})
			})
		})
	})

	Describe("SelectNextEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectNextEntry()
		})

		Context("directory has never been set", func() {
			BeforeEach(func() {
				navigator = new(Navigator)
			})

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
				var selectedIndex int

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for len(navigator.Entries())-navigator.SelectedIndex() > 1 {
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
				var selectedIndex int

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for len(navigator.Entries())-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}

					// Keep a reference to the last index.
					selectedIndex = navigator.SelectedIndex()
				})

				It("decrements the selected index by one", func() {
					Expect(navigator.SelectedIndex()).To(BeEquivalentTo(selectedIndex - 1))
				})
			})
		})
	})

	Describe("IntoSelectedEntry", func() {
		JustBeforeEach(func() {
			error = navigator.IntoSelectedEntry()
		})

		BeforeEach(func() {
			path, _ = os.Getwd()
			path += "/sample"
			navigator.SetWorkingDirectory(path)
		})

		Context("a directory is selected", func() {
			BeforeEach(func() {
				for navigator.SelectedEntry().Name != "directory" {
					navigator.SelectNextEntry()
				}
			})

			It("navigates into the selected entry", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(path + "/directory"))
			})

			It("does not return an error", func() {
				Expect(error).To(BeNil())
			})
		})

		Context("a file is selected", func() {
			BeforeEach(func() {
				for navigator.SelectedEntry().Name != "file" {
					navigator.SelectNextEntry()
				}
			})

			It("does not change the working directory", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(path))
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})
		})
	})

	Describe("RemoveSelectedEntry", func() {
		var file_name, directory_name string
		var entryCount int

		JustBeforeEach(func() {
			error = navigator.RemoveSelectedEntry()
		})

		Context("selected entry is a file", func() {
			BeforeEach(func() {
				file_name = "new_file"
				os.Create(file_name)

				// Update the navigator's cached entries.
				navigator.SetWorkingDirectory(originalPath)

				// Keep a reference to the size of the original entry set.
				entryCount = len(navigator.Entries())

				// Select the newly created file.
				for navigator.SelectedEntry().Name != file_name {
					navigator.SelectNextEntry()
				}
			})

			It("deletes the file", func() {
				_, err := os.Stat(file_name)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			It("removes the file from the navigator's entries", func() {
				file_entry := &Entry{Name: file_name, Size: 0, IsDirectory: false}
				Expect(navigator.Entries()).ToNot(ContainElement(file_entry))
			})

			It("leaves the navigator with the correct number of entries", func() {
				Expect(len(navigator.Entries())).To(Equal(entryCount - 1))
			})
		})

		Context("selected entry is a directory with files", func() {
			BeforeEach(func() {
				file_name = "new_file"
				directory_name = "new_directory"
				os.Mkdir(directory_name, 0700)
				os.Create(directory_name + "/" + file_name)

				// Update the navigator's cached entries.
				navigator.SetWorkingDirectory(originalPath)

				for navigator.SelectedEntry().Name != directory_name {
					navigator.SelectNextEntry()
				}
			})

			It("deletes the directory", func() {
				_, err := os.Stat(directory_name)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})

		Describe("selected entry after removal", func() {
			var first_file_name, second_file_name, last_file_name string
			BeforeEach(func() {
				// Create a directory.
				directory_name = "new_directory"
				os.Mkdir(directory_name, 0700)

				// Create three files in that directory, using numbers to guarantee sorting.
				first_file_name = "1"
				second_file_name = "2"
				last_file_name = "3"
				os.Create(directory_name + "/" + first_file_name)
				os.Create(directory_name + "/" + second_file_name)
				os.Create(directory_name + "/" + last_file_name)

				// Navigate into the new directory
				navigator.SetWorkingDirectory(originalPath + "/" + directory_name)
			})

			AfterEach(func() {
				os.RemoveAll(directory_name)
			})

			Context("selected entry is the first entry", func() {
				It("selects the second entry", func() {
					Expect(navigator.SelectedEntry().Name).To(Equal(second_file_name))
				})
			})

			Context("selected entry is the second entry", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
				})

				It("selects the second entry", func() {
					Expect(navigator.SelectedEntry().Name).To(Equal(last_file_name))
				})
			})

			Context("selected entry is the last entry", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
				})

				It("selects the second entry", func() {
					Expect(navigator.SelectedEntry().Name).To(Equal(second_file_name))
				})
			})
		})
	})

	Describe("ToParentDirectory", func() {
		var parent_path string

		JustBeforeEach(func() {
			error = navigator.ToParentDirectory()
		})

		Context("directory has a parent directory", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				parent_path = path + "/sample"
				path += "/sample/directory"
				navigator.SetWorkingDirectory(path)
			})

			It("navigates to the parent directory", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(parent_path))
			})
		})
	})

	Describe("View", func() {
		var buffer *view.Buffer
		var maxRows int

		JustBeforeEach(func() {
			buffer = navigator.View(maxRows)
		})

		It("returns the current directory path as its status", func() {
			Expect(buffer.Status).To(Equal(navigator.CurrentPath()))
		})

		Context("maxRows is set to 1", func() {
			BeforeEach(func() {
				maxRows = 1
			})

			It("returns a buffer with the right number of rows", func() {
				Expect(len(buffer.Rows)).To(BeEquivalentTo(maxRows))
			})

			It("stores the proper view data indices", func() {
				Expect(navigator.ViewDataIndices()).To(Equal([2]int{0, 1}))
			})

			Describe("returned row", func() {
				It("has its left value set to the first entry's name", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[0].Name))
				})

				It("has its right value set to the first entry's formatted size", func() {
					formattedSize := view.Size(navigator.Entries()[0].Size)
					Expect(buffer.Rows[0].Right).To(Equal(formattedSize))
				})

				It("has its highlight value set to the first entry's highlighted status", func() {
					Expect(buffer.Rows[0].Highlight).To(BeTrue())
				})

				Context("selected entry is a directory", func() {
					BeforeEach(func() {
						entry := navigator.SelectedEntry()

						for !entry.IsDirectory {
							navigator.SelectNextEntry()
							entry = navigator.SelectedEntry()
						}
					})

					It("has its colour value set to true", func() {
						Expect(buffer.Rows[0].Colour).To(BeTrue())
					})

					It("has a forward slash appended to its name", func() {
						Expect(buffer.Rows[0].Left).To(Equal(navigator.SelectedEntry().Name + "/"))
					})
				})

				Context("selected entry is not a directory", func() {
					BeforeEach(func() {
						entry := navigator.SelectedEntry()

						for entry.IsDirectory {
							navigator.SelectNextEntry()
							entry = navigator.SelectedEntry()
						}
					})

					It("has its colour value set to false", func() {
						Expect(buffer.Rows[0].Colour).To(BeFalse())
					})
				})
			})
		})

		Context("maxRows is set to 2", func() {
			BeforeEach(func() {
				maxRows = 2
			})

			It("returns a buffer with the right number of rows", func() {
				Expect(len(buffer.Rows)).To(BeEquivalentTo(maxRows))
			})

			It("stores the proper view data indices", func() {
				Expect(navigator.ViewDataIndices()).To(Equal([2]int{0, 2}))
			})

			Context("selected entry has never been changed", func() {
				It("returns the first and second rows", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[0].Name))
					Expect(buffer.Rows[1].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("sets the first row as highlighted", func() {
					Expect(buffer.Rows[0].Highlight).To(BeTrue())
				})
			})

			Context("the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
				})

				It("returns the first row", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[0].Name))
				})

				It("returns the second row", func() {
					Expect(buffer.Rows[1].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(buffer.Rows[1].Highlight).To(BeTrue())
				})
			})

			Context("the second entry is selected, the view is rendered, and then the third entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					_ = navigator.View(maxRows)
					navigator.SelectNextEntry()
				})

				It("returns the second row", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(buffer.Rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the third row as highlighted", func() {
					Expect(buffer.Rows[1].Highlight).To(BeTrue())
				})
			})

			Context("the third entry is selected, the view is rendered, and then the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					_ = navigator.View(maxRows)
					navigator.SelectPreviousEntry()
				})

				It("returns the second row", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(buffer.Rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(buffer.Rows[0].Highlight).To(BeTrue())
				})
			})

			Context("the fourth entry is selected, the view is rendered, and then the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					_ = navigator.View(maxRows)
					navigator.SelectPreviousEntry()
					navigator.SelectPreviousEntry()
				})

				It("returns the second row", func() {
					Expect(buffer.Rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(buffer.Rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(buffer.Rows[0].Highlight).To(BeTrue())
				})
			})
		})

		Context("in an empty directory", func() {
			var emptyDirectoryPath string

			BeforeEach(func() {
				// Create an empty directory.
				path, _ = os.Getwd()
				fileInfo, _ := os.Stat(path + "/sample/")
				emptyDirectoryPath = path + "/sample/empty"
				_ = os.Mkdir(emptyDirectoryPath, fileInfo.Mode())

				// Navigate into the empty directory.
				navigator.SetWorkingDirectory(emptyDirectoryPath)
			})

			AfterEach(func() {
				os.Remove(emptyDirectoryPath)
			})

			It("does not panic", func() {
				Expect(func() { navigator.View(1) }).ToNot(Panic())
			})

			It("returns an empty buffer", func() {
				buffer := navigator.View(1)
				Expect(len(buffer.Rows)).To(BeZero())
			})
		})

		Context("after a file has been removed", func() {
			BeforeEach(func() {
				maxRows = 1

				// Create a new file.
				file_name := "new_file"
				os.Create(file_name)

				// Update the navigator's cached entries.
				navigator.SetWorkingDirectory(originalPath)

				// Select the newly created file.
				for navigator.SelectedEntry().Name != file_name {
					navigator.SelectNextEntry()
				}

				// Remove it.
				navigator.RemoveSelectedEntry()
			})

			It("does not panic", func() {
				Expect(func() { navigator.View(maxRows) }).ToNot(Panic())
			})
		})
	})
})
