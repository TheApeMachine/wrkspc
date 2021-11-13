package bcknd

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Store is anything that can hold data and be queried for it.
*/
type Store interface {
	Peek(*spdg.Datagram) chan *spdg.Datagram
	Poke(*spdg.Datagram)
}

/*
NewStore constructs a Store of the type that is passed in.
*/
func NewStore(storeType Store) Store {
	errnie.Traces()
	return storeType
}
