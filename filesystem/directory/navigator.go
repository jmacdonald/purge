package directory

type Navigator struct {
	CurrentPath   string
	SelectedIndex uint16
	Entries       []*Entry
}

func (navigator *Navigator) ChangeDirectory(path string) {
	navigator.CurrentPath = path
	navigator.Entries = Entries(path)
}

func (navigator *Navigator) SelectNextEntry() {
	if uint16(len(navigator.Entries))-navigator.SelectedIndex > 1 {
		navigator.SelectedIndex++
	}
}
