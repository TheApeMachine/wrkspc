package bcknd

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/sockpuppet"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Server exposes and orchestrates an Ingress and Egress.
*/
type Server struct {
	ingress *Ingress
	egress  *Egress
	conn    *http.Server
	handler *Handler
}

/*
NewServer sets up the end-to-end networking for an Ingress and Egress.
*/
func NewServer() *Server {
	errnie.Traces()
	egress := NewEgress()
	hub := sockpuppet.NewHub()

	return &Server{
		ingress: NewIngress(),
		egress:  egress,
		handler: NewHandler(egress, hub),
		conn: &http.Server{
			Addr:         ":" + viper.GetString("wrkspc.bcknd.port"),
			ReadTimeout:  viper.GetDuration("wrkspc.bcknd.read-timeout") * time.Second,
			WriteTimeout: viper.GetDuration("wrkspc.bcknd.write-timeout") * time.Second,
		},
	}
}

/*
Up brings the server... Well, up.
*/
func (server *Server) Up() chan *spdg.Datagram {
	errnie.Traces()

	out := make(chan *spdg.Datagram)

	go func() {
		defer close(out)

		router := mux.NewRouter()
		router.Use(mux.CORSMethodMiddleware(router))

		router.HandleFunc("/v1/secure", server.handler.Response)
		router.HandleFunc("/v1/stream", server.handler.Stream)
		router.HandleFunc("/_status/healthz", server.handler.Health)

		// First the server is started, and blocks the goroutine.
		// Wrapped in errnie Handler, if there is an error, use NOOP (no operation).
		// Then call Encode on the stored error.
		// Then send it out the channel.
		errnie.Logs("bcknd going up on port ", server.conn.Addr).With(errnie.INFO)
		out <- spdg.QuickDatagram(
			spdg.ERROR, "error",
			errnie.Handles(
				server.conn.ListenAndServe(),
			).With(errnie.NOOP).ERR.Encode(),
		)
	}()

	return out
}
