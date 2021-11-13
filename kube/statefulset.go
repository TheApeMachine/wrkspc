package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSet struct {
	base   Base
	client RestClient
	spec   *appsv1.StatefulSet
	file   []byte
	handle *appsv1.StatefulSet
}

func NewStatefulSet(client RestClient, file []byte) MigratableKind {
	return &StatefulSet{
		client: client,
		file:   file,
	}
}

func (kind *StatefulSet) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.AppsV1().StatefulSets(kind.spec.Namespace).Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *StatefulSet) Check() bool {
	check, _ := kind.client.conn.AppsV1().StatefulSets(kind.spec.Namespace).Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *StatefulSet) Down() error {
	return kind.base.teardown()
}

func (kind *StatefulSet) Delete() error {
	return kind.client.conn.AppsV1().StatefulSets(kind.spec.Namespace).Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *StatefulSet) Name() string {
	return kind.handle.Name
}
