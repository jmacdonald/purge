/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

type Viewer interface {
	View(maxRows uint16) []Row
}

// Encapsulates information require to draw a row of information.
//
// Left and right represent two columns with matching alignment.
// Highlight inverts the row's colours, useful for "selecting" a row.
type Row struct {
	Left      string
	Right     string
	Highlight bool
}
