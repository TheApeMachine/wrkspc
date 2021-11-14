package hefner

import (
	"github.com/google/uuid"
)

/*
Initialize was supposed to be for any type that wanted to initialize iteself.
Then I was violently reminded how absolutely satanistic Go's type system is.
This can only work for the Pipe type, obvs.
*/
type Initializer interface {
	Initialize(Pipe) Pipe
}

/*
Pipe is a data channel.
*/
type Pipe interface {
	Initializer
	Generator
}

/*
ProtoPipe is the default, proto-typical Pipe.
*/
type ProtoPipe struct {
	i   Pipe
	o   chan Pipe
	key uuid.UUID
}

/*
NewPipe sets up a new Pipe of the type that is passed in.
*/
func NewPipe(pipeType Pipe, i Pipe) Pipe {
	return pipeType.Initialize(i)
}

/*
Initialize the Pipe.
*/
func (pipe ProtoPipe) Initialize(i Pipe) Pipe {
	if i == nil {
		i = ProtoPipe{}
	}

	pipe.i = i
	pipe.o = make(chan Pipe)
	pipe.key = uuid.New()

	return pipe
}
