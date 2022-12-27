package drknow

import (
	"bytes"
	"io"

	"github.com/theapemachine/wrkspc/spd"
)

/*
Abstract is a wrapper around any other object that implements the
io.ReadWriteCloser interface.
*/
type Abstract struct {
	concrete *spd.Datagram
	buffer   *bytes.Buffer
}

/*
NewAbstract constructs an Abstract instance and returns a pointer
reference to it.
*/
func NewAbstract(concrete any) *Abstract {
	return &Abstract{concrete, bytes.NewBuffer([]byte{})}
}

/*
Read implements the io.Reader interface.
*/
func (abstract *Abstract) Read(p []byte) (n int, err error) {
	n = copy(p, abstract.buffer.Bytes())
	return n, io.EOF
}

/*
Write implements the io.Writer interface.
*/
func (abstract *Abstract) Write(p []byte) (n int, err error) {
	return abstract.buffer.Write(p)
}

/*
Close implements the io.Closer interface.
*/
func (abstract *Abstract) Close() error {
	return nil
}
