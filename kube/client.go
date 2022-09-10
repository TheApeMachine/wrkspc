package kube

import (
	clientset "github.com/minio/operator/pkg/client/clientset/versioned"
	promclientset "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*
Client wraps the various Kubernetes clients that are needed to manipulate both
Kubernetes native, as well as extended or custom resources.
*/
type Client struct {
	KubeClient       *kubernetes.Clientset
	ControllerClient *clientset.Clientset
	ExtClient        *apiextension.Clientset
	PromClient       *promclientset.Clientset
}

/*
NewClient returns a handle on the various clients that we will need access to.
*/
func NewClient() Client {
	config, err := clientcmd.BuildConfigFromFlags("", brazil.BuildPath(
		brazil.HomePath(), "/.kube/config",
	))
	errnie.Handles(err).With(errnie.NOOP)

	kubeClient, err := kubernetes.NewForConfig(config)
	errnie.Handles(err).With(errnie.NOOP)

	controllerClient, err := clientset.NewForConfig(config)
	errnie.Handles(err).With(errnie.NOOP)

	extClient, err := apiextension.NewForConfig(config)
	errnie.Handles(err).With(errnie.NOOP)

	promClient, err := promclientset.NewForConfig(config)
	errnie.Handles(err).With(errnie.NOOP)

	return Client{
		KubeClient:       kubeClient,
		ControllerClient: controllerClient,
		ExtClient:        extClient,
		PromClient:       promClient,
	}
}
