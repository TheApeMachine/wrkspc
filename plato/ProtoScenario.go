package plato

import (
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/spdg"
)

type ProtoScenario struct {
	namespace string
	scenes    []Scene
}

func (scenario ProtoScenario) initialize(namespace string) Scenario {
	// errnie.Ambient().Log(errnie.DEBUG, "plato.ProtoScenario.initialize <-", namespace)

	scenario.namespace = namespace
	scenario.scenes = scenario.loadScenes()

	return scenario
}

func (scenario ProtoScenario) loadScenes() []Scene {
	sc := viper.GetStringSlice(
		"scenarios." + scenario.namespace + ".scenes",
	)

	var scenes []Scene

	for _, scene := range sc {
		scenes = append(
			scenes,
			NewScene(
				ProtoScene{},
				scene,
				viper.GetString("scenes."+scene+".message"),
				viper.GetStringMapString("scenes."+scene+".arguments"),
			),
		)
	}

	return scenes
}

/*
Run is a contained runtime environment. Inside it loops over a channel
where it receives anonymous functions that receive and provide data.
Each function Run receives is executed, passing in the data from the
previous function.
*/
func (scenario ProtoScenario) Run(datagrams chan spdg.Datagram) {
	// errnie.Ambient().Log(errnie.DEBUG, "plato.ProtoScenario.Run <-", datagrams)

	for dg := range datagrams {
		for _, scene := range scenario.scenes {
			artifact := spdg.Datagram{}
			dg.Unwrap(&artifact) <- <-dg.Generate()
			scene.Action(artifact)
		}
	}
}
