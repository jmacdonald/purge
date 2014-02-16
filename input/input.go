package input

import "io"

type Navigator interface {
	SelectNextEntry()
	SelectPreviousEntry()
	IntoSelectedEntry() error
	ToParentDirectory() error
}

func Read(source io.Reader, navigator Navigator) {
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
