package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct {
	base   Base
	client RestClient
	spec   *apiv1.ConfigMap
	file   []byte
	handle *apiv1.ConfigMap
}

func NewConfigMap(client RestClient, file []byte) MigratableKind {
	return &ConfigMap{
		client: client,
		file:   file,
	}
}

func (kind *ConfigMap) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.CoreV1().ConfigMaps(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *ConfigMap) Check() bool {
	check, _ := kind.client.conn.CoreV1().ConfigMaps(kind.spec.Namespace).Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *ConfigMap) Down() error {
	return kind.base.teardown()
}

func (kind *ConfigMap) Delete() error {
	return kind.client.conn.CoreV1().ConfigMaps(kind.handle.Namespace).Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *ConfigMap) Name() string {
	return kind.handle.Name
}
