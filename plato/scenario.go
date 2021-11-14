package plato

import "github.com/theapemachine/wrkspc/spdg"

/*
Scenario...
*/
type Scenario interface {
	initialize(string) Scenario
	Run(chan spdg.Datagram)
}

/*
NewScenario...
*/
func NewScenario(scenarioType Scenario, namespace string) Scenario {
	return scenarioType.initialize(namespace)
}
