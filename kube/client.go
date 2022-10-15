package kube

import (
	"context"
	"os"
	"strings"

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
		"type":  "kubectl",
		"multi": "false",
		"sub":   "",
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
		"type":  "kubectl",
		"multi": "false",
		"sub":   "",
	},
	"prometheus": {
		"type":  "kuebctl",
		"multi": "true",
		"sub":   "setup",
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
	errnie.Informs("applying", name)

	if manifests[name]["type"] == "helm" {
		// It is a helm chart, so hand it off to helm.
		client.helm(name, vendor, namespace)
		return
	}

	// We have a standard Kubernetes manifest and can just apply it.
	applyOpts := apply.NewApplyOptions(
		client.dynamicClient, client.discoveryClient,
	)

	n := []string{}

	if manifests[name]["sub"] != "" {
		// There is a sub directory which needs to be handled first.
		// Likely used for some initial setup needed to deploy the main
		// manifest files.
		n = append(n, name+"/"+manifests[name]["sub"])
	}

	// Append the main manifest files to the directory list.
	n = append(n, name)

	// Loop over the path listing.
	for _, s := range n {
		// Start a generator that yields filenames.
		for fh := range brazil.GeneratePath(
			brazil.BuildPath(brazil.Workdir(), ".kubernetes", s),
		) {
			go func(f string) {
				for {
					// Read the file and store as a byte slice.
					data, err := os.ReadFile(f)

					errnie.Handles(err)

					if len(data) > 0 {
						err := applyOpts.Apply(context.TODO(), data)
						errnie.Handles(err)

						if err == nil {
							break
						}
					}
				}
			}(fh)
		}
	}
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
		CleanupOnFail:   true,
	}

	_, err = client.HelmClient.InstallOrUpgradeChart(
		context.Background(), &chartSpec, nil,
	)

	errnie.Handles(err)
}
