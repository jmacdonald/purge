/*
Package view implements display-related functionality, such as
data formatting and display updates (using termbox).
*/
package view

type Viewer interface {
	View(maxRows uint16) [][2]string
}
