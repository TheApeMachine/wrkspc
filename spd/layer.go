package spd

/*
Layer is a wrapper around a version of a Datagram payload.
A layer is Write Once Read Many (WORM) and is used to build
an immutable history of a datagram.
*/
type Layer struct {
	payload []byte
}

/*
NewLayer creates a new Layer.
*/
func NewLayer() *Layer {
	return &Layer{
		payload: nil,
	}
}

/*
Read implements the io.Reader interface for Layer.
*/
func (l *Layer) Read(p []byte) (n int, err error) {
	return copy(p, l.payload), nil
}

/*
Write implements the io.Writer interface for Layer.
*/
func (l *Layer) Write(p []byte) (n int, err error) {
	l.payload = p
	return len(p), nil
}
