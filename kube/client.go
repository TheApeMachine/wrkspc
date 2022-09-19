package kube

import (
	"context"
	"os"

	clientset "github.com/minio/operator/pkg/client/clientset/versioned"
	promclientset "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/pytimer/k8sutil/apply"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*
Client wraps the various Kubernetes clients that are needed to manipulate both
Kubernetes native, as well as extended or custom resources.
*/
type Client struct {
	KubeClient       *kubernetes.Clientset
	dynamicClient    dynamic.Interface
	discoveryClient  *discovery.DiscoveryClient
	ControllerClient *clientset.Clientset
	ExtClient        *apiextension.Clientset
	PromClient       *promclientset.Clientset
}

/*
NewClient returns a handle on the various clients that we will need access to.
*/
func NewClient() *Client {
	config, err := clientcmd.BuildConfigFromFlags("", brazil.BuildPath(
		brazil.HomePath(), "/.kube/config",
	))
	errnie.Handles(err)

	kubeClient, err := kubernetes.NewForConfig(config)
	errnie.Handles(err)

	dynamicClient, err := dynamic.NewForConfig(config)
	errnie.Handles(err)

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	errnie.Handles(err)

	controllerClient, err := clientset.NewForConfig(config)
	errnie.Handles(err)

	extClient, err := apiextension.NewForConfig(config)
	errnie.Handles(err)

	promClient, err := promclientset.NewForConfig(config)
	errnie.Handles(err)

	return &Client{
		KubeClient:       kubeClient,
		dynamicClient:    dynamicClient,
		discoveryClient:  discoveryClient,
		ControllerClient: controllerClient,
		ExtClient:        extClient,
		PromClient:       promClient,
	}
}

func (client *Client) Apply(fname string) {
	applyOpts := apply.NewApplyOptions(client.dynamicClient, client.discoveryClient)
	errnie.Handles(applyOpts.Apply(context.TODO(), os.ReadFile(fname)))
}
