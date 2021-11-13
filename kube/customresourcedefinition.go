package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CustomResourceDefinition struct {
	base   Base
	client ExtendedClient
	spec   *apiextensionsv1.CustomResourceDefinition
	file   []byte
	handle *apiextensionsv1.CustomResourceDefinition
}

func NewCustomResourceDefinition(client ExtendedClient, file []byte) MigratableKind {
	return &CustomResourceDefinition{
		client: client,
		file:   file,
	}
}

func (kind *CustomResourceDefinition) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *CustomResourceDefinition) Check() bool {
	check, _ := kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil
}

func (kind *CustomResourceDefinition) Down() error {
	return kind.base.teardown()
}

func (kind *CustomResourceDefinition) Delete() error {
	return kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Delete(
		context.TODO(), kind.handle.Name, v1.DeleteOptions{},
	)
}

func (kind *CustomResourceDefinition) Name() string {
	return kind.handle.Name
}
