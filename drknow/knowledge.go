package drknow

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

type Knowledge struct {
	io.ReadWriteCloser
}

func NewKnowledge() *Knowledge {
	return &Knowledge{}
}

func (knowledge *Knowledge) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return
}

func (knowledge *Knowledge) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return
}

func (knowledge *Knowledge) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return errnie.NewError(nil)
}
