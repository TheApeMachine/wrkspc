package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Namespace ...
*/
type Namespace struct {
	base *Base
	spec *apiv1.Namespace
}

/*
NewNamespace ...
*/
func NewNamespace(client RestClient, file []byte) MigratableKind {
	kind := &Namespace{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *Namespace) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.CoreV1().Namespaces().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)
	kind.base.waiter()
}

/*
Check ...
*/
func (kind *Namespace) Check() bool {
	return kind.spec.Status.Phase == apiv1.NamespaceActive
}

/*
Delete ...
*/
func (kind *Namespace) Delete() error {
	return kind.base.client.conn.CoreV1().Namespaces().Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
