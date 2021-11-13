package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"
	storev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
CSIDriver storage something something...
*/
type CSIDriver struct {
	base *Base
	spec *storev1.CSIDriver
}

/*
NewCSIDriver ...
*/
func NewCSIDriver(client RestClient, file []byte) MigratableKind {
	kind := &CSIDriver{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *CSIDriver) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.StorageV1().CSIDrivers().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *CSIDriver) Check() bool {
	check, err := kind.base.client.conn.StorageV1().CSIDrivers().Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *CSIDriver) Delete() error {
	return kind.base.client.conn.StorageV1().CSIDrivers().Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
