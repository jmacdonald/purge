/*
Package directory implements functionality for navigating
and listing directories, including size calculations.

Directory paths are always returned without a trailing slash.
*/
package directory

import (
	"io/ioutil"
	"os"
	"sort"
)

// Structure representing a directory entry.
type Entry struct {
	Name        string
	Size        int64
	IsDirectory bool
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

// Returns a list of entries (and their sizes) for the given
// path. The current and parent (./..) entries are not included.
func Entries(path string) (entries []*Entry) {
	var size int64
	var directorySizeChannel chan *EntrySize
	var asyncSizeCount int

	// Read the directory entries.
	dirEntries, _ := ioutil.ReadDir(path)
	entries = make([]*Entry, len(dirEntries))

	// Allocate a buffered channel on which we'll receive
	// directory sizes from size-calculating goroutines.
	directorySizeChannel = make(chan *EntrySize, len(dirEntries))

	for index, entry := range dirEntries {
		entryInfo, _ := os.Stat(path + "/" + entry.Name())

		// Figure out the entry's size differently
		// depending on whether or not it's a directory.
		if entryInfo.IsDir() {
			// Calculate the directory's size asynchronously, passing the current
			// index so that we know where to put the result when we receive it later on.
			go Size(path+"/"+entry.Name(), index, directorySizeChannel)
			asyncSizeCount++
		} else {
			size = entryInfo.Size()
		}

		// Store the entry details.
		entries[index] = &Entry{Name: entry.Name(), Size: size, IsDirectory: entryInfo.IsDir()}
	}

	// Listen for the results of the async size calculations.
	for i := 0; i < asyncSizeCount; i++ {
		// Read a directory size from the channel.
		directorySize := <-directorySizeChannel

		// Update the stored entry size.
		entries[directorySize.Index].Size = directorySize.Size
	}

	// Sort the entries, casting them to their sortable equivalent.
	sort.Sort(sortableEntries(entries))

	return
}
