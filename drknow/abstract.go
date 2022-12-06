package drknow

import (
	"github.com/theapemachine/wrkspc/hefner"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Abstract is a composite object that acts as the single type
flowing through the system.

It is composed from the following sub types:

- twoface.Context
- spd.Datagram
- hefner.Pipe
*/
type Abstract struct {
	*twoface.Context
	*spd.Datagram
	*hefner.Pipe
}

func NewAbstract() *Abstract {
	return &Abstract{twoface.NewContext(), spd.Empty, hefner.NewPipe()}
}
