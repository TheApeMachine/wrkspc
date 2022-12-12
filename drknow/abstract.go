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

func NewAbstract(datagram *spd.Datagram) *Abstract {
	return &Abstract{
		twoface.NewContext(),
		datagram,
	}
}

func (abstract *Abstract) Read(p []byte) (n int, err error) {
	return abstract.datagram.Read(p)
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
