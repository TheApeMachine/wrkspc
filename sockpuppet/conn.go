package sockpuppet

import "github.com/theapemachine/wrkspc/errnie"

/*
Connection is a link between objects, agnostic about the protocol
these objects use, whether or not they use the same protocol, or where
the objects may be located within distributed systems.
*/
type Connection struct {
}

func NewConnection() *Connection {
	return &Connection{}
}

func (conn *Connection) Read(p []byte) (n int, err error) {
	return
}

func (conn *Connection) Write(p []byte) (n int, err error) {
	return
}

func (conn *Connection) Close() error {
	return errnie.Handles(nil)
}
