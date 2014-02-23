package directory

import "os"
import "errors"
import "path/filepath"

// Structure used to keep state when
// navigating directories and their entries.
type Navigator struct {
	currentPath   string
	selectedIndex uint16
	entries       []*Entry
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

// Generates a two-dimensional slice with all
// of the data required for display.
func (navigator *Navigator) View(maxRows uint16) (viewData [][2]string) {
	entryCount := len(navigator.Entries())
	viewData = make([][2]string, entryCount, entryCount)
	return
}
