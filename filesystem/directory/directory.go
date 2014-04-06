/*
Package directory implements functionality for navigating
and listing directories, including size calculations.

Directory paths are always returned without a trailing slash.
*/
package directory

import (
	"io/ioutil"
	"os"
)

// Structure representing a directory entry.
type Entry struct {
	Name        string
	Size        int64
	IsDirectory bool
}

// Calculates and returns the size (in
// bytes) of the directory for the given path.
func Size(path string) (size int64) {
	// Read the directory entries.
	entries, _ := ioutil.ReadDir(path)

	// Sum the entry sizes, recursing if necessary.
	for _, entry := range entries {
		if os.FileMode.IsDir(entry.Mode()) {
			size += Size(path + "/" + entry.Name())
		} else {
			size += entry.Size()
		}
	}
	return
}

// Returns a list of entries (and their sizes) for the given
// path. The current and parent (./..) entries are not included.
func Entries(path string) (entries []*Entry) {
	// Read the directory entries.
	dirEntries, _ := ioutil.ReadDir(path)
	entries = make([]*Entry, len(dirEntries))

	for index, entry := range dirEntries {
		fileInfo, _ := os.Stat(path + "/" + entry.Name())
		entries[index] = &Entry{Name: entry.Name(), Size: Size(path + "/" + entry.Name()), IsDirectory: fileInfo.IsDir()}
	}
	return
}
