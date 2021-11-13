package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
ConfigMap represent a Kubernetes Kind.
*/
type ConfigMap struct {
	base   *Base
	client RestClient
	spec   *apiv1.ConfigMap
	file   []byte
	handle *apiv1.ConfigMap
}

/*
NewConfigMap prepares the Kubernetes Kind.
*/
func NewConfigMap(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &ConfigMap{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up deployes the Kubernetes Kind.
*/
func (kind *ConfigMap) Up() {
	errnie.Traces()
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.handle, kind.base.err = kind.client.conn.CoreV1().ConfigMaps(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check to see if the deployment was a success.
*/
func (kind *ConfigMap) Check() bool {
	errnie.Traces()

	check, _ := kind.client.conn.CoreV1().ConfigMaps(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil
}

/*
Delete removes the Kubernetes Kind.
*/
func (kind *ConfigMap) Delete() error {
	errnie.Traces()

	return kind.client.conn.CoreV1().ConfigMaps(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
