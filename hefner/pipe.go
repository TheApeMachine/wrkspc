package hefner

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

type PipeType uint

const (
	IPC PipeType = iota
	LAN
	WAN
)

type Pipe struct {
	r    *io.PipeReader
	w    *io.PipeWriter
	err  error
	Type PipeType
}

func NewPipe(t PipeType) *Pipe {
	r, w := io.Pipe()
	return &Pipe{r, w, nil, t}
}

func (pipe *Pipe) Read(p []byte) (n int, err error) {
	if n, err = pipe.r.Read(p); errnie.Handles(err) != nil {
		pipe.err = err
	}
	return
}

func (pipe *Pipe) Write(p []byte) (n int, err error) {
	go func() {
		defer pipe.w.Close()
		if n, err = pipe.w.Write(p); errnie.Handles(err) != nil {
			pipe.err = err
		}
	}()

	return
}

func (pipe *Pipe) Close() error {
	return pipe.err
}
