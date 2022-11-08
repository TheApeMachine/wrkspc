package sockpuppet

import (
	"fmt"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/valyala/fasthttp"
)

/*
Notary is an interface that manages data routing and validation
from the FastHTTP server.
*/
type Notary interface {
	Validate(string, *fasthttp.RequestCtx) (int, string)
}

/*
Contract is a wrapper around a smart contract function.
*/
type Contract struct {
	notary Notary
}

/*
NewContract instantiates a pointer to a Contract.
*/
func NewContract(notary Notary) *Contract {
	errnie.Traces()

	return &Contract{
		notary: notary,
	}
}

/*
Up brings up the smart contract function.
*/
func (contract *Contract) Up(port string) error {
	return fasthttp.ListenAndServe(
		":"+port, func(ctx *fasthttp.RequestCtx) {
			// Call the Validate method on the Notary object, which
			// manages the incoming request data and routes it to
			// the smart contract logic.
			code, response := contract.notary.Validate(string(ctx.Path()), ctx)

			// Set the response headers on the connection context.
			ctx.SetContentType("appplication/json")
			ctx.SetStatusCode(code)

			// Write the response data back to the calling connection.
			fmt.Fprint(ctx, response)
		},
	)
}
