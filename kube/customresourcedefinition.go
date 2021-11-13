package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
CustomResourceDefinition ...
*/
type CustomResourceDefinition struct {
	base   *Base
	client ExtendedClient
	spec   *apiextensionsv1.CustomResourceDefinition
}

/*
NewCustomResourceDefinition ...
*/
func NewCustomResourceDefinition(client ExtendedClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &CustomResourceDefinition{client: client}
	kind.base = NewBase(file, kind, RestClient{})
	return kind
}

/*
Up ...
*/
func (kind *CustomResourceDefinition) Up() {
	errnie.Traces()
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *CustomResourceDefinition) Check() bool {
	check, err := kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *CustomResourceDefinition) Delete() error {
	return kind.client.conn.ApiextensionsV1().CustomResourceDefinitions().Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
