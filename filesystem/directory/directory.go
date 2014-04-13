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

type EntrySize struct {
	Index int
	Size  int64
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
// Directory sizes are calculated asynchronously and are returned
// on the delayedEntrySizes channel, along with their array index.
func Entries(path string) (entries []*Entry, delayedEntrySizes chan *EntrySize, delayedEntryCount int) {
	var size int64

	// Read the directory entries.
	dirEntries, _ := ioutil.ReadDir(path)
	entries = make([]*Entry, len(dirEntries))

	// Allocate a channel to return delayed entry sizes.
	delayedEntrySizes = make(chan *EntrySize, len(dirEntries))

	for index, entry := range dirEntries {
		entryInfo, _ := os.Stat(path + "/" + entry.Name())

		// Figure out the entry's size differently
		// depending on whether or not it's a directory.
		if entryInfo.IsDir() {
			go Size(path+"/"+entry.Name(), index, delayedEntrySizes)
			delayedEntryCount++
		} else {
			size = entryInfo.Size()
		}

		entries[index] = &Entry{Name: entry.Name(), Size: size, IsDirectory: entryInfo.IsDir()}
	}
	return
}
