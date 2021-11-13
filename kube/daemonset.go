package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
DaemonSet ...
*/
type DaemonSet struct {
	base *Base
	spec *appsv1.DaemonSet
}

/*
NewDaemonSet ....
*/
func NewDaemonSet(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &DaemonSet{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *DaemonSet) Up() {
	errnie.Traces()
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.AppsV1().DaemonSets(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *DaemonSet) Check() bool {
	check, err := kind.base.client.conn.AppsV1().DaemonSets(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil || err == nil
}

/*
Delete ...
*/
func (kind *DaemonSet) Delete() error {
	return kind.base.client.conn.AppsV1().DaemonSets(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
