package double

import "github.com/nsf/termbox-go"

// Build a struct to track messages/method calls
// that implements the Renderer interface.
type Renderer struct {
	FlushCalled            bool
	ClearCalled            bool
	SetCellCalledCorrectly bool
}

func (r *Renderer) Flush() {
	r.FlushCalled = true
}

func (r *Renderer) Clear() {
	r.ClearCalled = true
}

func (renderer *Renderer) Size() (int, int) {
	return 15, 15
}

func (r *Renderer) SetCell(x, y int, value rune, fg, bg termbox.Attribute) {
	r.SetCellCalledCorrectly = true
}
