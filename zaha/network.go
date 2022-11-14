package zaha

import (
	"github.com/theapemachine/wrkspc/sockpuppet"
)

type Network struct {
	direction   *bool
	connections []sockpuppet.Conn
}

func NewNetwork(direction *bool, connections []sockpuppet.Conn) *Network {
	return &Network{
		direction:   direction,
		connections: connections,
	}
}
