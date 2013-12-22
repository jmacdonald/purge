package format

import "fmt"

func Size(size int64) (formattedSize string) {
	formattedSize = fmt.Sprintf("%v bytes", size)
	return
}
