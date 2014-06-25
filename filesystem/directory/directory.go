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
	Name           string
	Size           int64
	IsDirectory    bool
	SizeCalculated bool
}

type EntrySize struct {
	Index int
	Size  int64
}

// Alias a slice of entries so that
// we can implement sort.Interface.
type sortableEntries []*Entry

// Implement sort.Interface length function.
func (e sortableEntries) Len() int {
	return len(e)
}

// Implement sort.Interface comparison function,
// using the entry size as a comparator.
func (e sortableEntries) Less(i, j int) bool {
	if e[i].Size > e[j].Size {
		return true
	} else {
		return false
	}
}

// Implement sort.Interface swap method,
// used to re-arrange misplaced entries.
func (e sortableEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Calculates and returns the size (in
// bytes) of the directory for the given path.
func Size(path string, index int, entrySizeChannel chan *EntrySize) {
	var size int64

	// Read the directory entries.
	entries, _ := ioutil.ReadDir(path)

	// Sum the entry sizes, recursing if necessary.
	for _, entry := range entries {
		if os.FileMode.IsDir(entry.Mode()) {
			// Allocate a channel to receive the size.
			recursiveResult := make(chan *EntrySize)

			// Recurse with a useless index (we're only summarizing, we don't care about order),
			// blocking until we receive an answer (recursive calls don't need to be async).
			go Size(path+"/"+entry.Name(), 0, recursiveResult)
			size += (<-recursiveResult).Size
		} else {
			size += entry.Size()
		}
	}

	// Send the entry size on to the return channel.
	entrySizeChannel <- &EntrySize{Index: index, Size: size}
}
