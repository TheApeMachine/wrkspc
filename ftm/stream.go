package ftm

/*
Stream is a wrapper that mirrors the FollowTheMoney Stream object.
*/
type Stream struct {
	things []Thing
}

/*
NewStream creates a new Stream object.
*/
func NewStream(things []Thing) *Stream {
	return &Stream{things}
}

/*
Read implements the io.Reader interface.
*/
func (s *Stream) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface.
*/
func (s *Stream) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (s *Stream) Close() error {
	return nil
}
