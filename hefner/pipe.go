package hefner

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

type Pipe struct {
	origin io.ReadWriteCloser
}

func NewPipe(origin io.ReadWriteCloser) *Pipe {
	return &Pipe{origin}
}

func (pipe *Pipe) Read(p []byte) (n int, err error) {
	errnie.Trace()
	return pipe.origin.Read(p)
}

func (pipe *Pipe) Write(p []byte) (n int, err error) {
	errnie.Trace()
	return pipe.origin.Write(p)
}

func (pipe *Pipe) Close() error {
	errnie.Trace()
	return pipe.origin.Close()
}
