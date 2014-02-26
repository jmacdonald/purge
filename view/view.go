/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

import "github.com/nsf/termbox-go"

type Viewer interface {
	View(maxRows uint16) []Row
}

/*
Renderer defines the expectations for a display endpoint,
such that view.Render can properly the data passed to it.
*/
type Renderer interface {
	Flush()
	Clear()
	Size() (int, int)
	SetCell(int, int, rune, termbox.Attribute, termbox.Attribute)
}

/*
Encapsulates information require to draw a row of information.

Left and right represent two columns with matching alignment.
Highlight inverts the row's colours, useful for "selecting" a row.
*/
type Row struct {
	Left      string
	Right     string
	Highlight bool
}

func Render(data Viewer, output Renderer) {
	output.Clear()
	output.Flush()
}
