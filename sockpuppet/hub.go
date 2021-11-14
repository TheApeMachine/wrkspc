package sockpuppet

import (
	"github.com/google/uuid"
	"github.com/theapemachine/errnie/v2"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Hub controls the message flow between peers. Technically clients could also be sending
messages over this connection towards the Hub, which the server could be listening for
and this may be useful in the future, but for now we do a listen only model.
*/
type Hub struct {
	secureChannels []string
	clients        map[*WSClient]bool
	Broadcast      chan spdg.Datagram
	register       chan *WSClient
	unregister     chan *WSClient
}

// NewHub returns a reference to a new instance of Hub.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan spdg.Datagram),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		clients:    make(map[*WSClient]bool),
	}
}

/*
NewSecureChannel uses a key system to link an identifier to a client connected to the hub.
While this indeed provides an API key style of security layer, it also makes it possible to
target messages to a particular client or group, which is essential otherwise all data will
go to all clients and that makes no sense.
*/
func (hub *Hub) NewSecureChannel() string {
	hub.secureChannels = append(
		hub.secureChannels, uuid.New().String(),
	)

	return hub.secureChannels[len(hub.secureChannels)-1]
}

/*
Run the hub and expose the three needed functionalities:
1. Accept new listeners.
2. Unregister clients that voluntarily disconnect.
3. Broadcast new loglines to all listeners, or kick them if there is an error.
*/
func (hub *Hub) Run() {
	errnie.Logs.Info("websocket hub going for up")

	for {
		// Perform a complete cycle of the hub.
		hub.cycle()
	}
}

/*
cycle switches over the three possible states to see what events are happening
at that moment and handles them as needed.
*/
func (hub *Hub) cycle() {
	select {
	case client := <-hub.register:
		hub.clients[client] = true
		errnie.Logs.Info("new client connected")
	case client := <-hub.unregister:
		if _, ok := hub.clients[client]; ok {
			delete(hub.clients, client)
			close(client.send)
		}
		errnie.Logs.Warning("client disconnected")
	case message := <-hub.Broadcast:
		hub.handleBroadcast(message)
	}
}

/*
handleBroadcast fans out a message to all connected clients/listeners.
This is useful, because it means we only have to have a single log stream
open, no matter how many services want to ingest it.
*/
func (hub *Hub) handleBroadcast(datagram spdg.Datagram) {
	errnie.TraceIn()

	for client := range hub.clients {
		hub.sendMessage(client, datagram)
	}
}

/*
sendMessage exists because I don't like nesting too many levels deep.
*/
func (hub *Hub) sendMessage(client *WSClient, datagram spdg.Datagram) {
	errnie.TraceIn()

	select {
	case client.send <- []byte(datagram.Data.Body.Payload):
		// It's an empty case because evaluating the case and sending
		// the message is the same thing here.
	default:
		// If we're not sending a message, then we're in error
		// so close the connection.
		close(client.send)
		delete(hub.clients, client)

		errnie.Logs.Warning("client was kicked")
	}
}
