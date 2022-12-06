package kube

import (
	"context"

	"github.com/theapemachine/wrkspc/errnie"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Deployment struct {
	name     string
	manifest *appsv1.Deployment
	err      error
}

/*
NewDeployment constructs the Deployment manifest.
*/
func NewDeployment(name, tag, cmd string, args ...string) *Deployment {
	return &Deployment{
		name: name,
		manifest: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: NewDeploymentSpec(name, tag, cmd, args...),
		},
	}
}

/*
Drop the Deployment onto the Cluster.
*/
func (deployment *Deployment) Drop(
	clientset *kubernetes.Clientset,
) *Deployment {
	client := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	_, deployment.err = client.Create(
		context.Background(), deployment.manifest, metav1.CreateOptions{},
	)

	errnie.Handles(deployment.err)

	deployment.err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, deployment.err = client.Get(
			context.Background(), deployment.name, metav1.GetOptions{},
		)

		errnie.Handles(deployment.err)
		return deployment.err
	})

	errnie.Handles(deployment.err)
	return deployment
}
