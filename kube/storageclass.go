package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	storev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StorageClass struct {
	base   Base
	client RestClient
	spec   *storev1.StorageClass
	file   []byte
	handle *storev1.StorageClass
}

func NewStorageClass(client RestClient, file []byte) MigratableKind {
	return &StorageClass{
		client: client,
		file:   file,
	}
}

func (kind *StorageClass) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.StorageV1().StorageClasses().Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *StorageClass) Check() bool {
	check, _ := kind.client.conn.StorageV1().StorageClasses().Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *StorageClass) Down() error {
	return kind.base.teardown()
}

func (kind *StorageClass) Delete() error {
	return kind.client.conn.StorageV1().StorageClasses().Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *StorageClass) Name() string {
	return kind.handle.Name
}
