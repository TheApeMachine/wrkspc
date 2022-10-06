package sockpuppet

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
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
			data := ctx.PostBody()
			dg := spd.NewCached(
				fastjson.GetString(data, "role"),
				fastjson.GetString(data, "scope"),
				fastjson.GetString(data, "identity"),
				fastjson.GetString(data, "payload"),
			)

			switch fastjson.GetString(data, "role") {
			case "question":
				conn.manager.Read(dg)
			case "datapoint":
				conn.manager.Write(dg)
			}
		},
	)
}
