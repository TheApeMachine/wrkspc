package kube

import (
	openebsclient "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	"github.com/spf13/viper"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
}

type RestClient struct {
	cfg  *rest.Config
	conn *kube.Clientset
}

type ExtendedClient struct {
	cfg  *rest.Config
	conn *apiextensionsclient.Clientset
}

type OpenEBSClient struct {
	cfg  *rest.Config
	conn *openebsclient.Clientset
}

func NewRestClient() RestClient {
	errnie.Traces()

	cfg, err := clientcmd.RESTConfigFromKubeConfig(
		brazil.NewFile(
			viper.GetString("kubernetes.config"),
		).Data.Bytes(),
	)

	errnie.Handles(err).With(errnie.KILL)

	conn, err := kube.NewForConfig(cfg)
	errnie.Handles(err).With(errnie.KILL)

	return RestClient{
		cfg:  cfg,
		conn: conn,
	}
}

func NewExtendedClient(cfg Config) ExtendedClient {
	conn, err := apiextensionsclient.NewForConfig(cfg.restCfg)
	errnie.Handles(err).With(errnie.KILL)

	return ExtendedClient{
		cfg:  cfg.restCfg,
		conn: conn,
	}
}

func NewOpenEBSClient(cfg Config) OpenEBSClient {
	conn, err := openebsclient.NewForConfig(cfg.restCfg)
	errnie.Handles(err).With(errnie.KILL)

	return OpenEBSClient{
		cfg:  cfg.restCfg,
		conn: conn,
	}
}
