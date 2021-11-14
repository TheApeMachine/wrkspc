package keanu

import (
)

type Result struct{}

func NewResult(datagram *spdg.Datagram) chan *spdg.Datagram {
	return memory.NewTree().Peek(datagram)
}
