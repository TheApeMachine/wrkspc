package please

import (
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Request is an interface that can be implemented by any object that
wants to be able to act as a network request.
*/
type Request interface {
	Do(*spdg.Datagram) chan *spdg.Datagram
}

/*
NewRequest constructs a Request of the type that is passed in.
*/
func NewRequest(requestType Request) Request {
	return requestType
}
