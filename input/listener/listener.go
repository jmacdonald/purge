package listener

type Listener struct {
	input     io.Reader
	navigator Navigator
	data      []byte
}

func New() listener

func (listener *Listener) Read() {}
