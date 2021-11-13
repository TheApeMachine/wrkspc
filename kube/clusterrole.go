package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"
	rbac "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRole struct {
	base   Base
	client RestClient
	spec   *rbac.ClusterRole
	file   []byte
	handle *rbac.ClusterRole
}

func NewClusterRole(client RestClient, file []byte) MigratableKind {
	return &ClusterRole{
		client: client,
		file:   file,
	}
}

func (kind *ClusterRole) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.RbacV1().ClusterRoles().Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *ClusterRole) Check() bool {
	check, _ := kind.client.conn.RbacV1().ClusterRoles().Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *ClusterRole) Down() error {
	return kind.base.teardown()
}

func (kind *ClusterRole) Delete() error {
	return kind.client.conn.CoreV1().ServiceAccounts(kind.handle.Namespace).Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *ClusterRole) Name() string {
	return kind.handle.Name
}
