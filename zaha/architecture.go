package zaha

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/ford"
	"github.com/theapemachine/wrkspc/passepartout"
)

type Architecture struct {
	name    string
	network *Network
	stores  []passepartout.Store
}

func NewArchitecture(
	name string, network *Network, stores []passepartout.Store,
) *Architecture {
	return &Architecture{
		name:    name,
		network: network,
		stores:  stores,
	}
}

func (architecture *Architecture) Run(
	workspace ford.Workspace,
) errnie.Error {
	return workspace.Add(architecture)
}

func (architecture *Architecture) Do() errnie.Error {
	return errnie.NewError(nil)
}

func (architecture *Architecture) Read(p []byte) (n int, err error) {
	return
}

func (architecture *Architecture) Write(p []byte) (n int, err error) {
	return
}
