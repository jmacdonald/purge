package double

import "github.com/jmacdonald/liberator/view"

// Viewer double stores information require to display
// a view.Row, with the interface expected by view.Viewer.
type Viewer struct {
	Left      string
	Right     string
	Highlight bool
}

// View returns a single row containing the values of the Viewer struct.
func (viewer *Viewer) View(maxRows uint16) (rows []view.Row) {
	rows = make([]view.Row, 1, 1)
	rows[0] = view.Row{viewer.Left, viewer.Right, viewer.Highlight}

	return
}
