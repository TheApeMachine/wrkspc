package plato

import (
	"github.com/theapemachine/wrkspc/hefner"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Simulation ...
*/
type Simulation struct {
	scenario Scenario
	dataset  hefner.Pipe
	events   chan spdg.Datagram
}

/*
NewSimulation ...
*/
func NewSimulation(scenario Scenario, dataset hefner.Pipe) Simulation {
	return Simulation{
		scenario: scenario,
		dataset:  dataset,
		events:   make(chan spdg.Datagram),
	}
}

/*
Run ...
*/
func (sim Simulation) Run() error {
	sim.scenario.Run(sim.events)
	return nil
}
