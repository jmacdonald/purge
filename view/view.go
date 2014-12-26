/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

import "fmt"
import "strings"
import "github.com/nsf/termbox-go"
import "unicode/utf8"

// Buffer encapsulates all of the data required to render the view.
type Buffer struct {
	Rows   []Row
	Status [2]string
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

// Initialize prepares the screen for rendering, and should
// only be run once, before constructing a new view.
func Initialize() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

// Close is used to relinquish the screen so that
// it can be used after the program exits.
func Close() {
	termbox.Close()
}

/*
Construct a view that will listen for and render buffers sent to it.
Initialize and Close must be called before and after this function,
respectively.
*/
func New(buffers <-chan *Buffer) {
	for {
		// Wait for a buffer.
		buffer := <-buffers

		// Refresh the contents of the screen.
		err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		if err != nil {
			return
		}

		// Render the source one row at a time.
		for index, row := range buffer.Rows {
			renderRow(row, index)
		}

		// Render the source's status string.
		renderStatus(buffer.Status)

		// Draw the contents to the screen.
		termbox.Flush()
	}
}

// Render a single row of data to the screen.
func renderRow(row Row, rowNumber int) {
	width, _ := termbox.Size()

	// Format the row such that it fills the screen,
	// and properly aligns the left/right columns.
	formattedRow, err := FormatRow(row, width)

	if err == nil {
		// Step through the formatted row one rune at a time,
		// printing the rune to the screen at the correct coordinates.
		for column, character := range formattedRow {
			fgColour, bgColour := termbox.ColorWhite, termbox.ColorBlack

			if row.Highlight {
				fgColour, bgColour = bgColour, fgColour
			}
			if row.Colour {
				fgColour = termbox.ColorYellow
			}

			termbox.SetCell(column, rowNumber, character, fgColour, bgColour)
		}
	}
}

// Render a status message to the bottom of the screen.
func renderStatus(status [2]string) {
	width, height := termbox.Size()

	// The status line components may be too long to fit on-screen. If that's the case,
	// we'll trim the left side of the path, since it's the least important
	// piece of information of the bunch.
	maximumLeftSideWidth := width - len(status[1]) - 1
	if len(status[0]) > maximumLeftSideWidth {
		// Figure out how much of a character surplus we have.
		excess := len(status[0]) - maximumLeftSideWidth

		// Trim the leading part of the path, adding an elipsis.
		status[0] = "..." + status[0][excess+3:]
	}

	// Build a string representing the status line contents, padding with spaces.
	padding := strings.Repeat(" ", (maximumLeftSideWidth-len(status[0]))) + " "
	statusLineContents := status[0] + padding + status[1]

	// Print the status to the bottom of the screen by stepping
	// through the bottom row one cell at a time and printing
	// a character from the status message, or a blank space,
	// until all of the row has been filled.
	for column, offset := 0, 0; column < width; column++ {
		var character rune
		var size int

		// Decode the next rune and advance the offset by its length,
		// or if we've already read the entire string, use a space instead.
		character, size = utf8.DecodeRune([]byte(statusLineContents)[offset:])
		offset += size

		// Print the character to the screen in a highlighted colour.
		termbox.SetCell(column, height-1, character, termbox.ColorBlack, termbox.ColorWhite)
	}
}

func Height() int {
	// Return a height one row smaller than the screen
	// height, so that we have room to render a status bar.
	_, height := termbox.Size()

	// If for some reason the height is zero or less,
	// just return zero to prevent runtime panics.
	if height-1 <= 0 {
		return 0
	} else {
		return height - 1
	}
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
