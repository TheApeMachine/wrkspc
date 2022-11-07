package sockpuppet

import (
	"fmt"
	"math/big"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/valyala/fasthttp"
)

type Notary interface {
	Validate(string, *fasthttp.RequestCtx) (*big.Int, error)
}

type Contract struct {
	notary Notary
}

func NewContract(notary Notary) *Contract {
	errnie.Traces()

	return &Contract{
		notary: notary,
	}
}

func (contract *Contract) Up(port string) error {
	return fasthttp.ListenAndServe(
		":"+port, func(ctx *fasthttp.RequestCtx) {
			code, err := contract.notary.Validate(string(ctx.Path()), ctx)

			ctx.SetContentType("appplication/json")
			ctx.SetStatusCode(int(code.Int64()))

			fmt.Fprint(ctx, err.Error())
		},
	)
}
