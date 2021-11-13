package contempt

import (
	"github.com/sfreiberg/simplessh"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/machine"
)

/*
SSH represents a Connector and Connection over secure shell.
*/
type SSH struct {
	Conn *simplessh.Client
	IP   string
	Auth machine.Credential
	OK   bool
}

/*
Dial the network entity and attempt to establish a Connection.
*/
func (connection SSH) Dial(auth machine.Credential) Connection {
	var err error

	// Remove the explicit username to use the `current` user.
	connection.Conn, err = simplessh.ConnectWithKeyFile(
		connection.IP, auth.Identifier(), brazil.HomePath()+".ssh/id_rsa",
	)

	// We are connected using the keyfile, so return the Connection in its current state.
	if errnie.Handles(err).With(errnie.NOOP).OK {
		return connection
	}

	// Remove the explicit username to use the `current` user.
	connection.Conn, err = simplessh.ConnectWithPassword(
		connection.IP, auth.Identifier(), auth.Secret(),
	)

	// We are connected using username ans password, so return the Connection in its current state.
	if errnie.Handles(err).With(errnie.NOOP).OK {
		return connection
	}

	return connection
}

/*
Hangup on the Connection and perform the needed cleanup.
*/
func (connection SSH) Hangup() {}
