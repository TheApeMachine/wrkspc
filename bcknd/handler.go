package bcknd

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Handler ...
*/
type Handler struct {
	manager  *Manager
	disposer *twoface.Disposer
}

/*
NewHandler ...
*/
func NewHandler(egress *Egress) *Handler {
	errnie.Traces()
	disposer := twoface.NewDisposer()

	return &Handler{
		manager: NewManager(egress, disposer),
	}
}

/*
Response ...
*/
func (handler *Handler) Response(response http.ResponseWriter, request *http.Request) {
	errnie.Traces()

	// We want to be able to target data to this client, while not having them stick
	// around. For this we can use a key to act as part `security`, part `topic`.
	uid, err := uuid.NewUUID()
	if !errnie.Handles(err).With(errnie.NOOP).OK {
		go handler.disposer.Cleanup()
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.manager.Execute(spdg.ContextDatagram(
		spdg.QUESTION, spdg.NewAnnotation("lookup", "test"),
	))

	// Send the uid back in the response. This acts as a topical key for the request
	// so in a distributed data pipeline we can keep track of who to send the results.
	fmt.Fprintf(response, "%v", string(spdg.QuickDatagram(
		spdg.TOPIC, "json",
		bytes.NewBuffer([]byte(uid.String())),
	).Marshal()))
}
