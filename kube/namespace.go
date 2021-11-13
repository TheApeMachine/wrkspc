package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	base   Base
	client RestClient
	spec   *apiv1.Namespace
	file   []byte
	handle *apiv1.Namespace
}

func NewNamespace(client RestClient, file []byte) MigratableKind {
	return &Namespace{
		client: client,
		file:   file,
	}
}

func (kind *Namespace) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.CoreV1().Namespaces().Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *Namespace) Check() bool {
	return kind.handle.Status.Phase == apiv1.NamespaceActive
}

func (kind *Namespace) Down() error {
	return kind.base.teardown()
}

func (kind *Namespace) Delete() error {
	return kind.client.conn.CoreV1().Namespaces().Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *Namespace) Name() string {
	return kind.handle.Name
}
