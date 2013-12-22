package directory

import "io/ioutil"
import "os"

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
