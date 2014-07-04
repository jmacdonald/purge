package directory

import (
	"errors"
	"fmt"
	"github.com/jmacdonald/purge/view"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Structure used to keep state when
// navigating directories and their entries.
type Navigator struct {
	currentPath         string
	selectedIndex       int
	entries             []*Entry
	viewDataIndices     [2]int
	view                chan<- *view.Buffer
	DirectorySizes      chan *EntrySize
	pendingCalculations int
}

// NewNavigator constructs a new navigator object and waits indefinitely
// for commands sent to it. It sends an updated buffer whenever the
// navigator changes state.
// This function is meant to be run in a goroutine.
func NewNavigator(path string, commands <-chan string, buffers chan<- *view.Buffer) {
	navigator := new(Navigator)

	// Link the navigator up to the view.
	navigator.view = buffers

	// Set the initial working directory using
	// the path passed in as an argument.
	navigator.SetWorkingDirectory(path)

	for {
		select {
		case command := <-commands: // A command has arrived.
			// Invoke the command on the navigator.
			switch command {
			case "SelectNextEntry":
				navigator.SelectNextEntry()
			case "SelectPreviousEntry":
				navigator.SelectPreviousEntry()
			case "IntoSelectedEntry":
				navigator.IntoSelectedEntry()
			case "ToParentDirectory":
				navigator.ToParentDirectory()
			case "RemoveSelectedEntry":
				navigator.RemoveSelectedEntry()
			}

			// Refresh the view.
			buffers <- navigator.View(view.Height())

		case directorySize := <-navigator.DirectorySizes: // A directory size calculation has completed.
			// Update the stored entry size and flag it as calculated.
			navigator.entries[directorySize.Index].Size = directorySize.Size
			navigator.entries[directorySize.Index].SizeCalculated = true

			// Reduce this count so the view increases the completion percentage.
			navigator.pendingCalculations--

			// Update the view, since we have another directory size.
			navigator.view <- navigator.View(view.Height())
		}
	}
}

// Returns the navigator's current directory path.
func (navigator *Navigator) CurrentPath() string {
	return navigator.currentPath
}

// Returns the navigator's currently selected index.
func (navigator *Navigator) SelectedIndex() int {
	return navigator.selectedIndex
}

// Returns the navigator's current directory entries. This method does
// not read from disk and may not accurately reflect filesystem contents.
func (navigator *Navigator) Entries() []*Entry {
	return navigator.entries
}

// Returns the navigator's currently selected entry.
func (navigator *Navigator) SelectedEntry() *Entry {
	// Prevent an empty directory from accessing an out-of-bounds index.
	if navigator.SelectedIndex() < len(navigator.Entries()) {
		return navigator.Entries()[navigator.SelectedIndex()]
	}

	return nil
}

// Returns the last slice indices used by View(). This is only used internally, with the
// exception of tests, to provide view updates that take previous context into account.
func (navigator *Navigator) ViewDataIndices() [2]int {
	return navigator.viewDataIndices
}

// Sets the navigator's current directory path,
// fetches the entries for the newly changed directory,
// and resets the selected index to zero (if the directory is valid).
func (navigator *Navigator) SetWorkingDirectory(path string) (error error) {
	file, error := os.Stat(path)
	if error == nil && file.IsDir() {
		// Strip trailing slash, if present.
		if path[len(path)-1:] == "/" {
			path = path[:len(path)-1]
		}

		navigator.currentPath = path
		navigator.selectedIndex = 0
		navigator.viewDataIndices = [2]int{0, 0}
		navigator.populateEntries()
	} else if error == nil {
		error = errors.New("path is not a directory")
	}

	return
}

func (navigator *Navigator) populateEntries() {
	var size int64

	// Read the directory entries.
	dirEntries, _ := ioutil.ReadDir(navigator.currentPath)
	navigator.entries = make([]*Entry, len(dirEntries))

	// Allocate a buffered channel on which we'll receive
	// directory sizes from size-calculating goroutines.
	navigator.DirectorySizes = make(chan *EntrySize, len(dirEntries))

	// Reset the number of pending calculations.
	navigator.pendingCalculations = 0

	for index, entry := range dirEntries {
		entryInfo, _ := os.Stat(navigator.currentPath + "/" + entry.Name())

		// Figure out the entry's size differently
		// depending on whether or not it's a directory.
		if entryInfo.IsDir() {
			navigator.pendingCalculations++

			// Calculate the directory's size asynchronously, passing the current
			// index so that we know where to put the result when we receive it later on.
			go Size(navigator.currentPath+"/"+entry.Name(), index, navigator.DirectorySizes)
		} else {
			size = entryInfo.Size()
		}

		// Store the entry details.
		navigator.entries[index] = &Entry{Name: entry.Name(), Size: size, IsDirectory: entryInfo.IsDir(), SizeCalculated: !entryInfo.IsDir()}
	}

	// Update the view, since we have sizes for files.
	navigator.view <- navigator.View(view.Height())

	// Sort the entries, casting them to their sortable equivalent.
	// sort.Sort(sortableEntries(entries))
}

// Moves the selectedIndex to the next entry in the
// list, if the current selection isn't already at the end.
func (navigator *Navigator) SelectNextEntry() {
	if len(navigator.entries)-navigator.selectedIndex > 1 {
		navigator.selectedIndex++
	}
}

// Moves the selectedIndex to the last entry in the list.
func (navigator *Navigator) SelectLastEntry() {
	if len(navigator.entries) > 0 {
		navigator.selectedIndex = len(navigator.entries) - 1
	}
}

// Moves the selectedIndex to the previous entry in the
// list, if the current selection isn't already at the beginning.
func (navigator *Navigator) SelectPreviousEntry() {
	if navigator.selectedIndex > 0 {
		navigator.selectedIndex--
	}
}

// Moves the selectedIndex to the first entry in the list.
func (navigator *Navigator) SelectFirstEntry() {
	navigator.selectedIndex = 0
}

// Navigates into the selected entry, if it is a directory.
func (navigator *Navigator) IntoSelectedEntry() error {
	entry := navigator.SelectedEntry()
	return navigator.SetWorkingDirectory(navigator.CurrentPath() + "/" + entry.Name)
}

// Removes/deletes the selected entry.
func (navigator *Navigator) RemoveSelectedEntry() error {
	err := os.RemoveAll(navigator.CurrentPath() + "/" + navigator.SelectedEntry().Name)
	if err == nil {
		if navigator.selectedIndex == len(navigator.entries)-1 {
			navigator.selectedIndex = len(navigator.entries) - 2

			// Trim the last entry off of the slice
			navigator.entries = navigator.entries[0 : navigator.selectedIndex+1]
		} else {
			// Create a new slice of entries by combining slices surrounding the deleted entry.
			navigator.entries = append(navigator.entries[0:navigator.selectedIndex],
				navigator.entries[navigator.selectedIndex+1:]...)
		}
	}

	return err
}

// Navigates to the parent directory.
func (navigator *Navigator) ToParentDirectory() error {
	parent_path, error := filepath.Abs(navigator.CurrentPath() + "/..")
	if error != nil {
		return error
	}
	return navigator.SetWorkingDirectory(parent_path)
}

// Generates a slice of rows with all of the data required for display.
func (navigator *Navigator) View(maxRows int) *view.Buffer {
	var start, end, size int
	var entrySize string

	// Return the current directory path as the status.
	status := navigator.CurrentPath()

	// Append a percentage to the status line, if
	// we're still calculating directory sizes.
	if navigator.pendingCalculations > 0 {
		entryCount := len(navigator.entries)
		status += fmt.Sprintf(" (%d%%)", (entryCount-navigator.pendingCalculations)*100/entryCount)
	}

	// Create a slice with a size that is the lesser of the entry count and maxRows.
	entryCount := len(navigator.Entries())
	if maxRows > entryCount {
		size = entryCount
	} else {
		size = maxRows
	}
	viewData := make([]view.Row, size, size)

	// Don't bother going any further if there are no entries to work with.
	if size == 0 {
		return &view.Buffer{Rows: viewData, Status: status}
	}

	// Deleting an entry can result in the cached view range indices
	// being out of bounds; correct that if it has occurred.
	if navigator.viewDataIndices[1] > entryCount {
		navigator.viewDataIndices[1] = entryCount
	}

	// Determine the range of entries to return.
	if navigator.viewDataIndices[1] != 0 && navigator.viewDataIndices[0] <= navigator.SelectedIndex() &&
		navigator.SelectedIndex() < navigator.viewDataIndices[1] {

		// The selected entry is still visible in the slice last returned. Return
		// the same range of entries to keep the view as consistent as possible.
		start, end = navigator.viewDataIndices[0], navigator.viewDataIndices[1]

	} else if navigator.viewDataIndices[1] != 0 && navigator.SelectedIndex() < navigator.viewDataIndices[0] {

		// The selected entry is beneath the range of entries previously returned.
		// Shift the range down just enough to include the selected entry.
		start = navigator.SelectedIndex()
		end = navigator.SelectedIndex() + size

	} else if navigator.SelectedIndex() >= size {
		// The selected entry is either above the previously returned range, or
		// this function hasn't been called for the current directory yet. Either way,
		// it would be outside of the returned range if we started at zero, so return
		// a range with it at the top.
		start = navigator.SelectedIndex() + 1 - size
		end = navigator.SelectedIndex() + 1

	} else {
		// Use the range starting at index 0.
		start = 0
		end = size
	}

	// Copy the navigator entries' names and
	// formatted sizes into the slice we'll return.
	for i, entry := range navigator.Entries()[start:end] {
		highlight := i+int(start) == int(navigator.SelectedIndex())

		// Add a trailing slash to the name
		// if the entry is a directory.
		var name string
		if entry.IsDirectory {
			name = entry.Name + "/"
		} else {
			name = entry.Name
		}

		if entry.SizeCalculated {
			entrySize = view.Size(entry.Size)
		} else {
			entrySize = "Calculating..."
		}

		viewData[i] = view.Row{name, entrySize, highlight, entry.IsDirectory}
	}

	// Store the indices used to generate the view data.
	navigator.viewDataIndices = [2]int{start, end}

	return &view.Buffer{Rows: viewData, Status: status}
}
