package drknow

import (
	"io"

	"github.com/theapemachine/wrkspc/twoface"
)

/*
Abstract is a composite object that acts as the single type
flowing through the system.

It is composed from the following sub types:

- twoface.Context
- spd.Datagram
*/
type Abstract struct {
	ctx *twoface.Context
	rwc io.ReadWriteCloser
}

func NewAbstract(rwc io.ReadWriteCloser) *Abstract {
	return &Abstract{twoface.NewContext(), rwc}
}

func (abstract *Abstract) Read(p []byte) (n int, err error) {
	return abstract.rwc.Read(p)
}

func (abstract *Abstract) Write(p []byte) (n int, err error) {
	return abstract.rwc.Write(p)
}

func (abstract *Abstract) Close() error {
	return abstract.rwc.Close()
}
