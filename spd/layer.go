package spd

import "github.com/theapemachine/wrkspc/errnie"

/*
Payload returns the first level of the Layer stack in the Datagram.
*/
func (dg Datagram) Payload() []byte {
	list, err := dg.Layers()
	errnie.Handles(err)

	data, err := list.At(0)
	errnie.Handles(err)

	return data
}
