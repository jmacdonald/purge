package double

type Navigator struct {
	SelectNextEntryCalled     bool
	SelectPreviousEntryCalled bool
	IntoSelectedEntryCalled   bool
	ToParentEntryCalled       bool
}

func (n *Navigator) SelectNextEntry() {
	n.SelectNextEntryCalled = true
}

func (n *Navigator) SelectPreviousEntry() {
	n.SelectPreviousEntryCalled = true
}

func (n *Navigator) IntoSelectedEntry() error {
	n.IntoSelectedEntryCalled = true
	return nil
}

func (n *Navigator) ToParentEntry() error {
	n.ToParentEntryCalled = true
	return nil
}
