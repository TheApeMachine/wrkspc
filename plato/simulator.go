package plato

import (
	"errors"
)

/*
Simulator...
*/
type Simulator struct {
	simulation Simulation
}

/*
NewSimulator...
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
	if bad := errnie.Ambient().Log(errnie.FATAL, sim.simulation.Run()); bad {
		return errors.New("simulation died")
	}

	return nil
}
