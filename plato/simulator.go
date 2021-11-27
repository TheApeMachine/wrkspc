package plato

import (
	"errors"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Simulator ...
*/
type Simulator struct {
	simulation Simulation
}

/*
NewSimulator ...
*/
func NewSimulator(simulation Simulation) Simulator {
	return Simulator{
		simulation: simulation,
	}
}

/*
Run the simulation.
*/
func (sim Simulator) Run() error {
	if !errnie.Handles(sim.simulation.Run()).OK {
		return errors.New("simulation died")
	}

	return nil
}
