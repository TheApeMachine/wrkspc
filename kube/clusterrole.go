package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"
	rbac "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
ClusterRole ...
*/
type ClusterRole struct {
	base *Base
	spec *rbac.ClusterRole
}

/*
NewClusterRole ...
*/
func NewClusterRole(client RestClient, file []byte) MigratableKind {
	kind := &ClusterRole{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *ClusterRole) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.RbacV1().ClusterRoles().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *ClusterRole) Check() bool {
	check, err := kind.base.client.conn.RbacV1().ClusterRoles().Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *ClusterRole) Delete() error {
	return kind.base.client.conn.CoreV1().ServiceAccounts(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
