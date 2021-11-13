package hello

import "github.com/theapemachine/wrkspc/spdg"

/*
Operator dynamically builds connections between objects that want to call each other.
It does not matter where these objects are, could be they don't even exist on the same machine.
It does not matter what these objects are, could be the don't even have the same type.
*/
type Operator struct {
	sender   *spdg.Datagram
	receiver *spdg.Datagram
}

/*
NewOperator returns a reference to a new instance of Operator.
*/
func NewOperator() *Operator {
	return &Operator{}
}

/*
Connect two objects to each other, such that they can exchange messages.
Both objects should be traveling inside a Datagram, so we can pass in any (inner) type
we want. By anything I mean you could indeed also pass in one Operator with already objects
connected, and also pass in another object so you get a threeway conversation. Or pass in two
Operators with connected objects and make it 4 participants. Etc.
Consider this as a fanout kind of setup.
*/
func (operator *Operator) Connect(sender, receiver *spdg.Datagram) *Operator {
	return operator
}
