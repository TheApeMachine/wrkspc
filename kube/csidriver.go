package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"
	storev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CSIDriver struct {
	base   Base
	client RestClient
	spec   *storev1.CSIDriver
	file   []byte
	handle *storev1.CSIDriver
}

func NewCSIDriver(client RestClient, file []byte) MigratableKind {
	return &CSIDriver{
		client: client,
		file:   file,
	}
}

func (kind *CSIDriver) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.StorageV1().CSIDrivers().Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *CSIDriver) Check() bool {
	check, _ := kind.client.conn.StorageV1().CSIDrivers().Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *CSIDriver) Down() error {
	return kind.base.teardown()
}

func (kind *CSIDriver) Delete() error {
	return kind.client.conn.StorageV1().CSIDrivers().Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *CSIDriver) Name() string {
	return kind.handle.Name
}
