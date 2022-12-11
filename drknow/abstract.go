package drknow

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
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
	ctx      *twoface.Context
	datagram *spd.Datagram
}

func NewAbstract() *Abstract {
	return &Abstract{
		twoface.NewContext(),
		&spd.Empty,
	}
}

func (abstract *Abstract) Read(p []byte) (n int, err error) {
	if err = abstract.datagram.Encode(p); err != nil {
		return n, errnie.Handles(err)
	}

	return len(p), io.EOF
}

func (abstract *Abstract) Write(p []byte) (n int, err error) {
	if err = abstract.datagram.Decode(p); err != nil {
		return n, errnie.Handles(err)
	}

	return len(p), io.EOF
}

func (abstract *Abstract) Close() error {
	return nil
}
