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
