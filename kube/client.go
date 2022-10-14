package kube

import (
	"context"
	"os"
	"strings"
	"time"

	clientset "github.com/minio/operator/pkg/client/clientset/versioned"
	helmclient "github.com/mittwald/go-helm-client"
	promclientset "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/pytimer/k8sutil/apply"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"helm.sh/helm/v3/pkg/repo"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var manifests = map[string]map[string]string{
	"system": {
		"type": "kubectl",
	},
	"base": {
		"type": "helm",
		"url":  "https://istio-release.storage.googleapis.com/charts",
	},
	"istiod": {
		"type": "helm",
		"url":  "https://istio-release.storage.googleapis.com/charts",
	},
	"vault": {
		"type": "helm",
		"url":  "https://helm.releases.hashicorp.com",
	},
	"harbor": {
		"type": "helm",
		"url":  "https://helm.goharbor.io",
	},
	"minio": {
		"type": "kubectl",
	},
}

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
	HelmClient       helmclient.Client
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

	hc, err := helmclient.New(&helmclient.Options{})
	errnie.Handles(err)

	return &Client{
		KubeClient:       kubeClient,
		dynamicClient:    dynamicClient,
		discoveryClient:  discoveryClient,
		ControllerClient: controllerClient,
		ExtClient:        extClient,
		PromClient:       promClient,
		HelmClient:       hc,
	}
}

func (client *Client) Apply(name, vendor, namespace string) {
	if manifests[name]["type"] == "helm" {
		// It is a helm chart, so hand it off to helm.
		client.helm(name, vendor, namespace)
		return
	}

	// We have a standard Kubernetes manifest and can just apply it.
	applyOpts := apply.NewApplyOptions(
		client.dynamicClient, client.discoveryClient,
	)

	data, err := os.ReadFile(brazil.BuildPath(
		brazil.Workdir(), ".kubernetes", name, "deploy.yml",
	))

	errnie.Handles(err)
	errnie.Handles(applyOpts.Apply(context.TODO(), data))
}

func (client *Client) helm(name, vendor, namespace string) {
	// Add a chart-repository to the client.
	errnie.Handles(client.HelmClient.AddOrUpdateChartRepo(
		repo.Entry{
			Name: vendor,
			URL:  manifests[name]["url"],
		},
	))

	data, err := os.ReadFile(brazil.BuildPath(
		brazil.Workdir(), ".kubernetes", name, "values.yml",
	))
	errnie.Handles(err)

	chartSpec := helmclient.ChartSpec{
		ReleaseName:     name,
		ChartName:       strings.Join([]string{vendor, name}, "/"),
		Namespace:       namespace,
		ValuesYaml:      string(data),
		CreateNamespace: true,
		UpgradeCRDs:     true,
		Wait:            true,
		Timeout:         10 * time.Second,
		CleanupOnFail:   true,
	}

	_, err = client.HelmClient.InstallOrUpgradeChart(
		context.Background(), &chartSpec, nil,
	)

	errnie.Handles(err)
}
