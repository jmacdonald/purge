// The double package provides interface implementations
// used for integration testing message expectations.
package double

// Build a struct to track messages/method calls
type Navigator struct {
	SelectNextEntryCalled     bool
	SelectPreviousEntryCalled bool
	IntoSelectedEntryCalled   bool
	ToParentDirectoryCalled   bool
}

// Implement the Navigator interface methods as defined in input.go
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
