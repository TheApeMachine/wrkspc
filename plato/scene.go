package plato

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

var bcknd bool
var metrics bool

type Scene interface {
	initialize(string, string, map[string]string) Scene
	Action(spdg.Datagram) spdg.Datagram
}

func NewScene(
	sceneType Scene, name string, message string, args map[string]string,
) Scene {
	return sceneType.initialize(name, message, args)
}

type ProtoScene struct {
	name    string
	msg     string
	args    map[string]string
	handler func(spdg.Datagram, map[string]string) spdg.Datagram
}

var protoActions = map[string]func(spdg.Datagram, map[string]string) spdg.Datagram{
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

func (scene ProtoScene) Action(data spdg.Datagram) spdg.Datagram {
	return scene.handler(data, scene.args)
}

func postHttp(artifact spdg.Datagram, args map[string]string) spdg.Datagram {
	return artifact
}

func randomizeValues(data spdg.Datagram, args map[string]string) spdg.Datagram {
	return data
}

func failureRate(data spdg.Datagram, args map[string]string) spdg.Datagram {
	return data
}

func instanceBcknd(data spdg.Datagram, args map[string]string) spdg.Datagram {
	if args["scope"] == "global" && !bcknd {
		bcknd = true

		go func() {
			srv := metric.NewServer()

			if !errnie.Logs(srv.Up()).With(errnie.WARNING).OK {
				errnie.Logs("error instanciating bcknd").With(errnie.WARNING)
			}
		}()
	}

	return data
}

/*
exportMetrics to Prometheus.
*/
func exportMetrics(data spdg.Datagram, args map[string]string) spdg.Datagram {
	if args["scope"] == "global" && !metrics {
		metrics = true
		// Get a new ops exporter and start recording metrics on it.
		_ = metric.NewExporter("sim", promCache)
	}

	return data
}
