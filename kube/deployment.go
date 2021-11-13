package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	base   Base
	client RestClient
	spec   *appsv1.Deployment
	file   []byte
	handle *appsv1.Deployment
}

func NewDeployment(client RestClient, file []byte) MigratableKind {
	return &Deployment{
		client: client,
		file:   file,
	}
}

func (kind *Deployment) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.AppsV1().Deployments(kind.spec.Namespace).Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *Deployment) Check() bool {
	check, _ := kind.client.conn.AppsV1().Deployments(kind.spec.Namespace).Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *Deployment) Down() error {
	return kind.base.teardown()
}

func (kind *Deployment) Delete() error {
	return kind.client.conn.AppsV1().Deployments(kind.spec.Namespace).Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *Deployment) Name() string {
	return kind.handle.Name
}
