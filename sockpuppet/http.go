package sockpuppet

import "github.com/valyala/fasthttp"

type HTTP struct{}

func NewHTTP() Conn {
	return &HTTP{}
}

func (conn *HTTP) Up(port string) error {
	return fasthttp.ListenAndServe(
		":"+port, func(ctx *fasthttp.RequestCtx) {
		},
	)
}
