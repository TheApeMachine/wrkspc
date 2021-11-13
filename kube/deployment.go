package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Deployment ...
*/
type Deployment struct {
	base *Base
	spec *appsv1.Deployment
}

/*
NewDeployment ...
*/
func NewDeployment(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &Deployment{}
	kind.base = NewBase(file, kind, RestClient{})
	return kind
}

/*
Up ...
*/
func (kind *Deployment) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.AppsV1().Deployments(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *Deployment) Check() bool {
	check, err := kind.base.client.conn.AppsV1().Deployments(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *Deployment) Delete() error {
	return kind.base.client.conn.AppsV1().Deployments(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
