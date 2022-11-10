package sockpuppet

import (
	"context"
	"net/http"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocket struct {
	err     error
	conn    *websocket.Conn
	manager Manager
}

func NewWebsocket(manager Manager) *WebSocket {
	return &WebSocket{
		manager: manager,
	}
}

func (socket *WebSocket) Up(port string) error {
	return nil
}

func (socket *WebSocket) Handle(
	w http.ResponseWriter, r *http.Request,
) {
	errnie.Traces()

	socket.conn, socket.err = websocket.Accept(w, r, nil)
	errnie.Handles(socket.err)

	defer socket.conn.Close(websocket.StatusInternalError, socket.err.Error())

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	errnie.Handles(wsjson.Read(ctx, socket.conn, &v))
	errnie.Logs("ws rx:", v)

	socket.conn.Close(websocket.StatusNormalClosure, "")
}
