package plato

import (
	"github.com/theapemachine/wrkspc/spdg"
)

var bcknd bool
var metrics bool
var promCache balena.Client

type Scene interface {
	initialize(string, string, map[string]string) Scene
	Action(data.Artifact) data.Artifact
}

func NewScene(
	sceneType Scene, name string, message string, args map[string]string,
) Scene {
	errnie.Ambient().Log(errnie.DEBUG, "plato.Scene.NewScene", sceneType, name, message, args)
	return sceneType.initialize(name, message, args)
}

type ProtoScene struct {
	name    string
	msg     string
	args    map[string]string
	handler func(data.Artifact, map[string]string) data.Artifact
}

var protoActions = map[string]func(data.Artifact, map[string]string) data.Artifact{
	"post-http":        postHttp,
	"randomize-values": randomizeValues,
	"failure-rate":     failureRate,
	"instance-bcknd":   instanceBcknd,
	"export-metrics":   exportMetrics,
}

func (scene ProtoScene) initialize(name string, msg string, args map[string]string) Scene {
	scene.name = name
	scene.msg = msg
	scene.args = args
	scene.handler = protoActions[scene.name]

	return scene
}

func (scene ProtoScene) Action(data data.Artifact) data.Artifact {
	return scene.handler(data, scene.args)
}

func postHttp(artifact data.Artifact, args map[string]string) data.Artifact {
	datagram := spdg.NewDatagram(spdg.Payload{})
	request := please.NewRequest(please.Rest{}, please.POST, datagram)
	results, err := request.Do()

	errnie.Ambient().Log(errnie.ERROR, err)
	errnie.Ambient().Log(errnie.INFO, results)

	return artifact
}

func pluginAction(data data.Artifact, args map[string]string) data.Artifact {
	return data
}

func randomizeValues(data data.Artifact, args map[string]string) data.Artifact {
	return data
}

func failureRate(data data.Artifact, args map[string]string) data.Artifact {
	return data
}

func instanceBcknd(data data.Artifact, args map[string]string) data.Artifact {
	if args["scope"] == "global" && !bcknd {
		bcknd = true

		go func() {
			srv := metric.NewServer()

			if ok := errnie.Ambient().Log(errnie.WARNING, srv.Up()); !ok {
				errnie.Ambient().Log(errnie.ERROR, "error instanciating bcknd")
			}
		}()
	}

	return data
}

func exportMetrics(data data.Artifact, args map[string]string) data.Artifact {
	if args["scope"] == "global" && !metrics {
		metrics = true
		// Get a new ops exporter and start recording metrics on it.
		_ = metric.NewExporter("simbots", promCache)
	}

	return data
}
