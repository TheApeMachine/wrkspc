package ford

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

type Workspace struct {
	io.ReadWriteCloser
}

func NewWorkspace() *Workspace {
	return &Workspace{}
}

func (wrkspc *Workspace) Read(p []byte) (n int, err error) {
	return
}

func (wrkspc *Workspace) Write(p []byte) (n int, err error) {
	return
}

func (wrkspc *Workspace) Close() error {
	return errnie.NewError(nil)
}
