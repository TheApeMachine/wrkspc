package sockpuppet

import (
	context "context"

	"github.com/theapemachine/wrkspc/errnie"
)

type IngressServer struct {
	director Director
}

func NewIngressServer(director Director) Ingress {
	return IngressServer{director: director}
}

func (Ingress) Handle(ctx context.Context, call Ingress_handler) error {
	res, err := call.AllocResults()
	errnie.Handles(err)
	return nil
}
