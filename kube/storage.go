package kube

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7/pkg/set"
	miniov2 "github.com/minio/operator/pkg/apis/minio.min.io/v2"
	informers "github.com/minio/operator/pkg/client/informers/externalversions"
	"github.com/minio/operator/pkg/controller/cluster"
	"github.com/theapemachine/wrkspc/errnie"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
)

var version = "DEVELOPMENT.GOGET"

/*
Storage wraps the MinIO operator to deploy it as a Kubernetes storage class.
*/
type Storage struct {
	client Client
}

/*
NewStorage returns a handle on the wrapped MinIO operator.
*/
func NewStorage(cluster Cluster) Storage {
	return Storage{
		client: NewClient(),
	}
}

/*
Provision the MinIO storage class onto the Kubernetes cluster.
*/
func (storage Storage) Provision(stop chan struct{}) errnie.Error {
	ctx := context.Background()
	_, err := storage.client.KubeClient.CoreV1().Namespaces().Create(
		ctx,
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "minio-operator"},
		},
		metav1.CreateOptions{},
	)
	namespacesENv, isNamespaced := os.LookupEnv("WATCHED_NAMESPACE")
	var namespaces set.StringSet
	if isNamespaced {
		namespaces = set.NewStringSet()
		rawNamespaces := strings.Split(namespacesENv, ",")
		for _, nsStr := range rawNamespaces {
			if nsStr != "" {
				namespaces.Add(strings.TrimSpace(nsStr))
			}
		}
	}

	caContent := miniov2.GetPodCAFromFile()
	namespace := miniov2.GetNSFromFile()

	operatorTLSCert, err := storage.client.KubeClient.CoreV1().Secrets(string(namespace)).Get(
		ctx, cluster.OperatorTLSSecretName, metav1.GetOptions{},
	)
	if err == nil && operatorTLSCert != nil {
		if val, ok := operatorTLSCert.Data["public.crt"]; ok {
			caContent = append(caContent, val...)
		}
		if val, ok := operatorTLSCert.Data["tls.crt"]; ok {
			caContent = append(caContent, val...)
		}
		if val, ok := operatorTLSCert.Data["ca.crt"]; ok {
			caContent = append(caContent, val...)
		}
	}

	operatorCATLSCert, err := storage.client.KubeClient.CoreV1().Secrets(miniov2.GetNSFromFile()).Get(ctx, cluster.OperatorCATLSSecretName, metav1.GetOptions{})
	if err == nil && operatorCATLSCert != nil {
		if val, ok := operatorCATLSCert.Data["public.crt"]; ok {
			caContent = append(caContent, val...)
		}
		if val, ok := operatorCATLSCert.Data["tls.crt"]; ok {
			caContent = append(caContent, val...)
		}
		if val, ok := operatorCATLSCert.Data["ca.crt"]; ok {
			caContent = append(caContent, val...)
		}
	}
	if len(caContent) > 0 {
		crd, err := storage.client.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), "tenants.minio.min.io", metav1.GetOptions{})
		if err != nil {
		} else {
			crd.Spec.Conversion.Webhook.ClientConfig.CABundle = caContent
			crd.Spec.Conversion.Webhook.ClientConfig.Service.Namespace = miniov2.GetNSFromFile()
			_, err := storage.client.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Update(context.Background(), crd, metav1.UpdateOptions{})
			if err != nil {
			}
		}
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(storage.client.KubeClient, time.Second*30)
	minioInformerFactory := informers.NewSharedInformerFactory(storage.client.ControllerClient, time.Second*30)
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "operator-pod"
	}

	mainController := cluster.NewController(
		podName,
		namespaces,
		storage.client.KubeClient,
		storage.client.ControllerClient,
		storage.client.PromClient,
		kubeInformerFactory.Apps().V1().StatefulSets(),
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Pods(),
		minioInformerFactory.Minio().V2().Tenants(),
		kubeInformerFactory.Core().V1().Services(),
		"",
		version,
	)

	go kubeInformerFactory.Start(stop)
	go minioInformerFactory.Start(stop)

	return errnie.Handles(mainController.Start(2, stop)).With(errnie.NOOP)
}
