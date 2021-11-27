package sockpuppet

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Connection is an interface that abstracts socket-based connections.
*/
type Connection interface {
	Dial(string, string, twoface.Disposer) chan *spdg.Datagram
}

/*
NewConnection constructs a new websocket based client.
*/
func NewConnection(connectionType Connection) Connection {
	errnie.Traces()
	return connectionType
}

/*
ProtoConnection is the standard type that can perform the requirements we usually have for
a websocket based streaming connection between services.
*/
type ProtoConnection struct {
	Disposer *twoface.Disposer
}

/*
Dial outward to another service's egress stream.
*/
func (connection ProtoConnection) Dial(host, path string) chan *spdg.Datagram {
	errnie.Traces()

	u := url.URL{Scheme: "ws", Host: host, Path: path}
	errnie.Logs("connecting to", u.String()).With(errnie.INFO)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if !errnie.Handles(err).With(errnie.NOOP).OK {
		return nil
	}

	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)
		defer c.Close()
		connection.worker(c, out)
	}()

	return out
}

/*
worker is the main loop for our connection.
*/
func (connection ProtoConnection) worker(c *websocket.Conn, out chan *spdg.Datagram) {
	errnie.Traces()

	for {
		select {
		case <-connection.Disposer.Done():
			return
		default:
			_, message, err := c.ReadMessage()
			errnie.Handles(err).With(errnie.NOOP)

			// We just discovered another use-case for NullDatagrams!
			// Make an empty datagram and unmarshal the bytes coming over the websocket
			// connection, which are: a Datagram, just encoded to bytes for sending.
			dg := spdg.NullDatagram()
			out <- dg.Unmarshal(message)
		}
	}
}
