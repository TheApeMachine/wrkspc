package hefner

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Pipe defines the interface for Hefner Pipes, an object with only one method
that receives an object like itself and returns an object channel like itself.
This makes it very easy to turn an object into a Hefner Pipe and while initially
designed to abstract away the transportation layer in concurrent processing, they
can do other things as well.
*/
type Pipe interface {
	IO(chan *spdg.Datagram) chan *spdg.Datagram
}

/*
NewPipe constructs a Hefner Pipe of the type that is passed in.
*/
func NewPipe(pipeType Pipe) Pipe {
	errnie.Traces()
	return pipeType
}

/*
ProtoPipe is the canonical implementation of Hefner Pipes and is currently used to
make using Go channels and WebSockets transparent to work with. This means you can
arbitrarily define concurrent workflows that have local as well as remote processes.
*/
type ProtoPipe struct {
	Disposer *twoface.Disposer
	Operator func(chan *spdg.Datagram) chan *spdg.Datagram
}

/*
IO performs some operation on the incoming Pipe channel and sends some data out over
the outgoing Pipe channel. The two do not need to be related.
*/
func (pipe ProtoPipe) IO(in chan *spdg.Datagram) chan *spdg.Datagram {
	errnie.Traces()
	return twoface.NewGenerator(twoface.GoGenerator{
		Operator: pipe.Operator,
		Disposer: pipe.Disposer,
	}).Yield(in)
}
