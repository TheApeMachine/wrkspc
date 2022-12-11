package sockpuppet

import (
	"net"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/hefner"
)

/*
Conn is a link between objects, agnostic about the protocol
these objects use, whether or not they use the same protocol, or where
the objects may be located within distributed systems.

It implements the net.Conn interface native to Go, which makes it very
compatible as a generic network abstraction.
*/
type Conn struct {
	pipes map[hefner.PipeType][]*hefner.Pipe
}

func NewConn() *Conn {
	errnie.Trace()

	// Return a Conn instance with one instance of each PipeType.
	// TODO: Should auto scale, and load balance the amount of
	//       Pipes of each type.
	return &Conn{
		map[hefner.PipeType][]*hefner.Pipe{
			hefner.IPC: {hefner.NewPipe(hefner.IPC)},
			hefner.LAN: {hefner.NewPipe(hefner.LAN)},
			hefner.WAN: {hefner.NewPipe(hefner.WAN)},
		},
	}
}

func (conn *Conn) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Warns("not implemented")
	return
}

func (conn *Conn) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Warns("not implemented")
	return
}

func (conn *Conn) Close() error {
	errnie.Trace()
	errnie.Warns("not implemented")
	return errnie.Handles(nil)
}

func (conn *Conn) LocalAddr() net.Addr {
	errnie.Trace()
	errnie.Warns("not implemented")
	return &net.IPAddr{}
}

func (conn *Conn) RemoteAddr() net.Addr {
	errnie.Trace()
	errnie.Warns("not implemented")
	return &net.IPAddr{}
}

func (conn *Conn) SetDeadline(t time.Time) error {
	errnie.Trace()
	errnie.Warns("not implemented")
	return nil
}

func (conn *Conn) SetReadDeadline(t time.Time) error {
	errnie.Trace()
	errnie.Warns("not implemented")
	return nil
}

func (conn *Conn) SetWriteDeadline(t time.Time) error {
	errnie.Trace()
	errnie.Warns("not implemented")
	return nil
}
