package contempt

import "github.com/theapemachine/wrkspc/machine"

/*
Connection is an established uplink with a networked device.
*/
type Connection interface {
	Dial(machine.Credential) Connection
	Hangup()
}

/*
NewConnection constructs a Connection of the type that is passed in.
*/
func NewConnection(connectionType Connection) Connection {
	return connectionType
}
