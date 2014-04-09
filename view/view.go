/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

import "fmt"
import "github.com/nsf/termbox-go"
import "unicode/utf8"

/*
Viewer is an interface used by Render to standardize data
from a data source such that it can be displayed properly.
*/
type Viewer interface {
	View(maxRows uint16) ([]Row, string)
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
	Colour    bool
}

/*
Render a data source that implements the
Viewer interface to the terminal using termbox.
*/
func Render(source Viewer) {
	width, height := termbox.Size()

	// Request the view data with a row maximum that's
	// one row smaller than the screen height, so that
	// we can render a status bar.
	rows, status := source.View(uint16(height - 1))

	// Clear the screen so that we can render new content to it.
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		return
	}

	// Step through the data one row at a time.
	for row, rowData := range rows {
		// Format the row such that it fills the screen,
		// and properly aligns the left/right columns.
		formattedRow, err := FormatRow(rowData, width)

		if err == nil {
			// Step through the formatted row one rune at a time,
			// printing the rune to the screen at the correct coordinates.
			column := 0
			for _, character := range formattedRow {
				fgColour := termbox.ColorWhite
				bgColour := termbox.ColorBlack

				if rowData.Highlight {
					fgColour, bgColour = bgColour, fgColour
				}
				if rowData.Colour {
					fgColour = termbox.ColorYellow
				}

				termbox.SetCell(column, row, character, fgColour, bgColour)
				column++
			}
		}
	}

	// Print the status to the bottom of the screen by stepping
	// through the bottom row one cell at a time and printing
	// a character from the status message, or a blank space,
	// until all of the row has been filled.
	for column, offset := 0, 0; column < width; column++ {
		var character rune
		var size int

		// Decode the next rune and advance the offset by its length,
		// or if we've already read the entire string, use a space instead.
		if offset < len(status) {
			character, size = utf8.DecodeRune([]byte(status)[offset:])
			offset += size
		} else {
			character = ' '
		}

		// Print the character to the screen in a highlighted colour.
		termbox.SetCell(column, height-1, character, termbox.ColorBlack, termbox.ColorWhite)
	}

	// Draw new content to the screen.
	termbox.Flush()
}

/*
FormatRow returns a string with the row's left/right
elements placed at the far left/right with spaces in between.
*/
func FormatRow(row Row, size int) (string, error) {
	// Figure out how large the left field needs to be, including
	// padding, to have the right field properly aligned to size.
	leftSize := size - len(row.Right)

	// Don't bother trying to format this row if the left and
	// right columns can't be separated by at least one space.
	if leftSize <= len(row.Left) {
		return "", fmt.Errorf("view: formatting row to a size of %d"+
			" with '%s' and '%s' values is impossible", size, row.Left, row.Right)
	}

	// Generate a format string with the appropriate spacing.
	formatString := fmt.Sprintf("%%-%ds%%s", leftSize)

	// Generate and return the formatted row.
	return fmt.Sprintf(formatString, row.Left, row.Right), nil
}
