package sockpuppet

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocket struct {
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

	c, err := websocket.Accept(w, r, nil)
	errnie.Handles(err)

	defer c.Close(websocket.StatusInternalError, err.Error())

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	errnie.Handles(wsjson.Read(ctx, c, &v))

	log.Printf("received: %v", v)

	c.Close(websocket.StatusNormalClosure, "")
}
