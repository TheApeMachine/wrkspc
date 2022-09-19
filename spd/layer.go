package spd

import "github.com/theapemachine/wrkspc/errnie"

func (dg Datagram) Payload() []byte {
	list, err := dg.Layers()
	errnie.Handles(err)

	data, err := list.At(0)
	errnie.Handles(err)

	return data
}
