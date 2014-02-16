package double

type Navigator struct {
	SelectNextEntryCalled     bool
	SelectPreviousEntryCalled bool
	IntoSelectedEntryCalled   bool
	ToParentDirectoryCalled   bool
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

func (n *Navigator) ToParentDirectory() error {
	n.ToParentDirectoryCalled = true
	return nil
}
