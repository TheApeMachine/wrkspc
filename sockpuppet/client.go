package sockpuppet

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/theapemachine/errnie/v2"
	"github.com/theapemachine/wrkspc/spdg"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*
WSClient is a wrapper between the websocket connection and the Hub.
This object runs a concurrent read/write pump mechanism for every client.
*/
type WSClient struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

/*
readPump is called as a goroutine for every connected client and reads messages/events
from their channel.
*/
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	errnie.Handles(c.conn.SetReadDeadline(time.Now().Add(pongWait))).With(errnie.RECV)

	c.conn.SetPongHandler(func(string) error {
		return errnie.Handles(
			c.conn.SetReadDeadline(time.Now().Add(pongWait)),
		).With(errnie.RECV).ERR
	})

	for {
		_, message, err := c.conn.ReadMessage()
		errnie.Handles(err).With(errnie.KILL)

		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			errnie.Handles(err).With(errnie.KILL)
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.Broadcast <- *spdg.QuickDatagram(spdg.DATAPOINT, "json", bytes.NewBuffer(message))
	}
}

/*
writePump is started as a goroutine for every connected client and sends messages/events to their
channel when a broadcast happens.
*/
func (c *WSClient) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			errnie.Handles(
				c.conn.SetWriteDeadline(time.Now().Add(writeWait)),
			).With(errnie.KILL)

			if !ok {
				errnie.Handles(
					c.conn.WriteMessage(websocket.CloseMessage, []byte{}),
				).With(errnie.KILL)

				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)

			if err != nil {
				errnie.Handles(
					c.conn.WriteMessage(websocket.CloseMessage, []byte{}),
				).With(errnie.KILL)

				return
			}

			_, err = w.Write([]byte(message))

			if err != nil {
				errnie.Handles(
					c.conn.WriteMessage(websocket.CloseMessage, []byte{}),
				).With(errnie.KILL)

				return
			}

			n := len(c.send)

			for i := 0; i < n; i++ {
				errnie.Handles(w.Write(newline)).With(errnie.KILL)
				errnie.Handles(w.Write(<-c.send)).With(errnie.KILL)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			errnie.Handles(
				c.conn.SetWriteDeadline(time.Now().Add(writeWait)),
			).With(errnie.KILL)

			errnie.Handles(c.conn.WriteMessage(websocket.PingMessage, nil)).With(errnie.KILL)
		}
	}
}

/*
serveWs will be called by the handler defined in the router of the service.
This responds to the incoming requests and starts the process of upgrading the connection
protocol from plain http to ws.
*/
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	errnie.Log("upgrading client connection")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	errnie.Handles(err).With(errnie.KILL)

	client := &WSClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	errnie.Logs.Info("client fully upgraded start read/write pump")
	go client.writePump()
	go client.readPump()
}
