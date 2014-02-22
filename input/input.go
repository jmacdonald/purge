// The input package is responsible for reading input data
// and invoking the corresponding navigator actions.
package input

import "io"
import "unicode/utf8"

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

// Reads input data from source as a sequence of runes,
// invoking a corresponding navigator action for certain values.
func Process(source io.Reader, navigator Navigator) {
	data := make([]byte, 5, 5)
	bytesRead, error := source.Read(data)

	// If there's valid data to be read,
	// read it one rune at a time.
	if bytesRead > 0 && error == nil {
		for _, runeValue := range data {
			switch runeValue {
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
	}
}
