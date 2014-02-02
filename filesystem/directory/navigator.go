package directory

type Navigator struct {
	currentPath   string
	selectedIndex uint16
	entries       []*Entry
}

/* Accessor Methods */
func (navigator *Navigator) CurrentPath() string {
	return navigator.currentPath
}

func (navigator *Navigator) SelectedIndex() uint16 {
	return navigator.selectedIndex
}

func (navigator *Navigator) Entries() []*Entry {
	return navigator.entries
}

func (navigator *Navigator) ChangeDirectory(path string) {
	navigator.currentPath = path
	navigator.entries = Entries(path)
}

func (navigator *Navigator) SelectNextEntry() {
	if uint16(len(navigator.entries))-navigator.selectedIndex > 1 {
		navigator.selectedIndex++
	}
}
