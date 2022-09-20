package sockpuppet

import "github.com/theapemachine/wrkspc/twoface"

type Director struct {
	managers []Manager
}

func NewDirector(managers ...Manager) *Director {
	return &Director{
		managers: managers,
	}
}

func (director Director) Direct(ctx *twoface.Context) {
	for _, manager := range director.managers {
		manager.Manage(ctx)
	}
}
