package plato

import (
	"fmt"
)

/*
Simulation...
*/
type Simulation struct {
	scenario Scenario
	dataset  hefner.Pipe
	events   chan spdg.Datagram
}

/*
NewSimulation...
*/
func NewSimulation(scenario Scenario, dataset hefner.Pipe) Simulation {
	errnie.Ambient().Log(errnie.DEBUG, fmt.Sprintf("plato.NewSimulation <- %v, %v", scenario, dataset))

	return Simulation{
		scenario: scenario,
		dataset:  dataset,
		events:   make(chan spdg.Datagram),
	}
}

func (sim Simulation) Run() error {
	errnie.Ambient().Log(errnie.INFO, "running simulation")

	go func() {
		defer close(sim.events)

		for datagram := range sim.dataset.Generate() {
			datagram = <-datagram.Unwrap(&datagram)
			sim.events <- datagram.(spdg.Datagram)
		}
	}()

	sim.scenario.Run(sim.events)

	return nil
}
