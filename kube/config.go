package kube

import (
	"path/filepath"

	"github.com/theapemachine/wrkspc/errnie"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	restCfg *rest.Config
}

func NewConfig() Config {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	errnie.Handles(err).With(errnie.KILL)

	return Config{
		restCfg: config,
	}
}
