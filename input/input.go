// The input package is responsible for reading input data
// and invoking the corresponding navigator actions.
package input

import (
	"io"
	"unicode/utf8"
)

// Navigator defines the interface expected by the input package,
// so that navigator actions can be called based on input data.
type Navigator interface {
	SelectNextEntry()
	SelectPreviousEntry()
	IntoSelectedEntry() error
	ToParentDirectory() error
}

// Reads and returns a single rune from the provided source.
func Read(source io.Reader) (value rune) {
	data := make([]byte, 4, 4)
	bytesRead, error := source.Read(data)

	// If there's valid data to be read,
	// read it one rune at a time.
	if bytesRead > 0 && error == nil {
		value, _ = utf8.DecodeRune(data)
	}
	return
}

// Maps characters to their corresponding navigator actions.
func Map(character rune, navigator Navigator) {
	switch character {
	case 'j':
		navigator.SelectNextEntry()
	case 'k':
		navigator.SelectPreviousEntry()
	case '\n':
		navigator.IntoSelectedEntry()
	case 'h':
		navigator.ToParentDirectory()
	}
}
