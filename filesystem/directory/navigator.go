package directory

import (
	"errors"
	"github.com/jmacdonald/liberator/view"
	"os"
	"path/filepath"
)

// Structure used to keep state when
// navigating directories and their entries.
type Navigator struct {
	currentPath     string
	selectedIndex   uint16
	entries         []*Entry
	viewDataIndices [2]uint16
}

// NewNavigator constructs a new navigator object.
func NewNavigator(path string) (navigator *Navigator) {
	navigator = new(Navigator)
	navigator.SetWorkingDirectory(path)
	return
}

// Returns the navigator's current directory path.
func (navigator *Navigator) CurrentPath() string {
	return navigator.currentPath
}

// Returns the navigator's currently selected index.
func (navigator *Navigator) SelectedIndex() uint16 {
	return navigator.selectedIndex
}

// Returns the navigator's current directory entries. This method does
// not read from disk and may not accurately reflect filesystem contents.
func (navigator *Navigator) Entries() []*Entry {
	return navigator.entries
}

// Returns the last slice indices used by View(). This is only used internally, with the
// exception of tests, to provide view updates that take previous context into account.
func (navigator *Navigator) ViewDataIndices() [2]uint16 {
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
		navigator.entries = Entries(path)
		navigator.selectedIndex = 0
		navigator.viewDataIndices = [2]uint16{0, 0}
	} else if error == nil {
		error = errors.New("path is not a directory")
	}

	return
}

// Moves the selectedIndex to the next entry in the
// list, if the current selection isn't already at the end.
func (navigator *Navigator) SelectNextEntry() {
	if uint16(len(navigator.entries))-navigator.selectedIndex > 1 {
		navigator.selectedIndex++
	}
}

// Moves the selectedIndex to the previous entry in the
// list, if the current selection isn't already at the beginning.
func (navigator *Navigator) SelectPreviousEntry() {
	if navigator.selectedIndex > 0 {
		navigator.selectedIndex--
	}
}

// Navigates into the selected entry, if it is a directory.
func (navigator *Navigator) IntoSelectedEntry() error {
	entry := navigator.Entries()[navigator.SelectedIndex()]
	return navigator.SetWorkingDirectory(navigator.CurrentPath() + "/" + entry.Name)
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
func (navigator *Navigator) View(maxRows uint16) (viewData []view.Row) {
	var start, end, size uint16

	// Create a slice with a size that is the lesser of the entry count and maxRows.
	entryCount := len(navigator.Entries())
	if maxRows > uint16(entryCount) {
		size = uint16(entryCount)
	} else {
		size = maxRows
	}
	viewData = make([]view.Row, size, size)

	// Determine the range of entries to return.
	if navigator.viewDataIndices[1] != 0 && navigator.viewDataIndices[0] <= navigator.SelectedIndex() &&
		navigator.SelectedIndex() <= navigator.viewDataIndices[1] {

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
		highlight := i == int(navigator.SelectedIndex())

		// Add a trailing slash to the name
		// if the entry is a directory.
		var name string
		if entry.IsDirectory {
			name = entry.Name + "/"
		} else {
			name = entry.Name
		}

		viewData[i] = view.Row{name, view.Size(entry.Size), highlight, entry.IsDirectory}
	}

	// Store the indices used to generate the view data.
	navigator.viewDataIndices = [2]uint16{start, end}

	return
}
