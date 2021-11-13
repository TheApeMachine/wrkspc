package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
PersistentVolumeClaim ...
*/
type PersistentVolumeClaim struct {
	base *Base
	spec *apiv1.PersistentVolumeClaim
}

/*
NewPersistentVolumeClaim ...
*/
func NewPersistentVolumeClaim(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &PersistentVolumeClaim{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *PersistentVolumeClaim) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.CoreV1().PersistentVolumeClaims(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *PersistentVolumeClaim) Check() bool {
	check, err := kind.base.client.conn.CoreV1().PersistentVolumeClaims(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)

	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *PersistentVolumeClaim) Delete() error {
	return kind.base.client.conn.CoreV1().PersistentVolumeClaims(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
