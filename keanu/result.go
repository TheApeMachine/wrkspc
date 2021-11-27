package keanu

import (
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Result ...
*/
type Result struct{}

/*
NewResult ...
*/
func NewResult(datagram *spdg.Datagram) chan *spdg.Datagram {
	return NewTree().Peek(datagram)
}
