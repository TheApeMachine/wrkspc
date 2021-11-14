package sockpuppet

import (
	"bytes"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/theapemachine/errnie/v2"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Connection is an interface that abstracts socket-based connections.
At this time it will be used to create the client-side for websocket
based connections between services, such that they can consume each other's
egress streams. At some point we will need to migrate the random bits of websocket
server and hub code from the skynet package as well.
*/
type Connection interface {
	Dial(string, string, twoface.Disposer) chan *spdg.Datagram
}

/*
NewConnection constructs a new websocket based client.
*/
func NewConnection(connectionType Connection) Connection {
	return connectionType
}

/*
ProtoConnection is the standard type that can perform the requirements we usually have for
a websocket based streaming connection between services.
*/
type ProtoConnection struct{}

/*
Dial outward to another service's egress stream.
*/
func (connection ProtoConnection) Dial(
	host, path string, disposer twoface.Disposer,
) chan *spdg.Datagram {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	errnie.Logs.Info("connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if !errnie.Handles(err).With(errnie.RECV).OK {
		return nil
	}

	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)
		defer c.Close()

		for {
			select {
			case <-disposer.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				errnie.Handles(err).With(errnie.RECV)

				// We just discovered another use-case for NullDatagrams!
				// Make an empty datagram and unmarshal the bytes coming over the websocket
				// connection, which are: a Datagram, just encoded to bytes for sending.
				dg := spdg.NullDatagram()
				out <- dg.Unmarshal(*bytes.NewBuffer(message))
			}
		}
	}()

	return out
}
