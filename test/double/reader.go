package double

type Reader string

func (s Reader) Read(target []byte) (int, error) {
	return copy(target, s), nil
}
