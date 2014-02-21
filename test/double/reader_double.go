package double

// Alias the string type to create a fake data source
type Reader string

// Implement the Reader interface, as required by input.go
func (s Reader) Read(target []byte) (int, error) {
	return copy(target, s), nil
}
