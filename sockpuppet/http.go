package sockpuppet

import "github.com/valyala/fasthttp"

type HTTP struct {
	manager Manager
}

func NewHTTP(manager Manager) Conn {
	return &HTTP{
		manager: manager,
	}
}

func (conn *HTTP) Up(port string) error {
	return fasthttp.ListenAndServe(
		":"+port, func(ctx *fasthttp.RequestCtx) {
		},
	)
}
