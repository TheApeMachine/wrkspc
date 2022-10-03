package sockpuppet

import (
	context "context"

	"github.com/davecgh/go-spew/spew"
	"github.com/theapemachine/wrkspc/errnie"
)

type IngressServer struct {
}

func NewIngressServer() IngressServer {
	return IngressServer{}
}

func (Ingress) Handle(ctx context.Context, call Ingress_handler) error {
	res, err := call.AllocResults()
	errnie.Handles(err)
	spew.Dump(res)
	return nil
}
