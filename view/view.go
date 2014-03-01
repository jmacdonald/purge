/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

import "fmt"
import "github.com/nsf/termbox-go"

/*
Viewer is an interface used by Render to standardize data
from a data source such that it can be displayed properly.
*/
type Viewer interface {
	View(maxRows uint16) []Row
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
