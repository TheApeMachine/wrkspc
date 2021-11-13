package kube

import (
	"path/filepath"

	"github.com/theapemachine/wrkspc/errnie"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

/*
Config represents the rest API configuration needed to interact with the Kubernetes API.
*/
type Config struct {
	restCfg *rest.Config
}

/*
NewConfig generates a new Config for us to pass to a Rest client.
*/
func NewConfig() Config {
	config, err := clientcmd.BuildConfigFromFlags(
		"", filepath.Join(homedir.HomeDir(), ".kube", "config"),
	)
	errnie.Handles(err).With(errnie.KILL)

	return Config{
		restCfg: config,
	}
}
