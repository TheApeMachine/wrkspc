package http

import (
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/valyala/fasthttp"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Up(port string) error {
	fasthttp.ListenAndServeTLS(
		"0.0.0.0:"+port,
		"/etc/ssl/certs/tls.crt", "/etc/ssl/certs/tls.key",
		func(ctx *fasthttp.RequestCtx) {
			dg := spd.New(
				ctx.Request.Header.ContentType(),
				spd.REQUEST, spd.HTTP,
			)

			// Write the entire request as a Layer in the
			// Datagram instance.
			ctx.Request.WriteTo(dg)
			drknow.NewAbstract(dg)
		},
	)

	return nil
}
