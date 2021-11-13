package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
ServiceAccount ...
*/
type ServiceAccount struct {
	base *Base
	spec *apiv1.ServiceAccount
}

/*
NewServiceAccount ...
*/
func NewServiceAccount(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &ServiceAccount{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *ServiceAccount) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.CoreV1().ServiceAccounts(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *ServiceAccount) Check() bool {
	check, err := kind.base.client.conn.CoreV1().ServiceAccounts(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *ServiceAccount) Delete() error {
	return kind.base.client.conn.CoreV1().ServiceAccounts(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
