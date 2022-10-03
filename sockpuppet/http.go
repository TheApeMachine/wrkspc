package sockpuppet

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/valyala/fasthttp"
)

type HTTP struct {
	manager Manager
}

func NewHTTP(manager Manager) Conn {
	errnie.Traces()
	return &HTTP{
		manager: manager,
	}
}

func (conn *HTTP) Up(port string) error {
	errnie.Debugs("HTTP going for up on", port)
	return fasthttp.ListenAndServe(
		":"+port, func(ctx *fasthttp.RequestCtx) {
			spew.Dump(ctx.PostBody())

			dg := spd.Unmarshal(ctx.PostBody())
			role, err := dg.Role()
			errnie.Handles(err)

			switch role {
			case "question":
				conn.manager.Read(ctx.PostBody())
			case "datapoint":
				conn.manager.Write(ctx.PostBody())
			}
		},
	)
}
