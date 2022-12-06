package hefner

import (
	"io"

	"github.com/google/uuid"
	ipc "github.com/theapemachine/golang-ipc"
	"github.com/theapemachine/wrkspc/errnie"
)

type PipeType uint

const (
	IPC PipeType = iota
	LAN
	WAN
)

func (t PipeType) BuildPipe() *Pipe {
	var server *ipc.Server

	switch t {
	case IPC:
		server, err := ipc.StartServer(uuid.NewString(), nil)
		if errnie.Handles(err) != nil {
			return nil
		}
	}

	return server
}

type Pipe struct {
	r   *io.PipeReader
	w   *io.PipeWriter
	err error
}

func NewPipe() *Pipe {
	r, w := io.Pipe()
	return &Pipe{r, w, nil}
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
