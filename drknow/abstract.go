package drknow

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Abstract is an item that holds both data and behavior,
which we can group as `knowledge` and `skill`.
*/
type Abstract interface {
	io.ReadWriteCloser
}

func NewAbstract(abstractType Abstract) Abstract {
	errnie.Trace()

	if abstractType == nil {
		return &ProtoAbstract{}
	}

	return abstractType
}

/*
ProtoAbstract is composed of Knowledge and Skill, which can be seen
as data and behavior.
*/
type ProtoAbstract struct {
	*Knowledge
	*Skill
}

func NewProtoAbstract(
	knowledge *Knowledge, skill *Skill,
) *ProtoAbstract {
	return &ProtoAbstract{knowledge, skill}
}

func (proto *ProtoAbstract) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return
}

func (proto *ProtoAbstract) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return
}

func (proto *ProtoAbstract) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")

	return errnie.NewError(nil)
}
