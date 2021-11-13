package kube

import (
	"context"

	"github.com/ghodss/yaml"
	openebsapiv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/theapemachine/wrkspc/errnie"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CStorPoolCluster struct {
	base   Base
	client OpenEBSClient
	spec   *openebsapiv1.CStorPoolCluster
	file   []byte
	handle *openebsapiv1.CStorPoolCluster
}

func NewCStorPoolCluster(client OpenEBSClient, file []byte) MigratableKind {
	return &CStorPoolCluster{
		client: client,
		file:   file,
	}
}

func (kind *CStorPoolCluster) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *CStorPoolCluster) Check() bool {
	check, _ := kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil
}

func (kind *CStorPoolCluster) Down() error {
	return kind.base.teardown()
}

func (kind *CStorPoolCluster) Delete() error {
	return kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Delete(
		context.TODO(), kind.handle.Name, v1.DeleteOptions{},
	)
}

func (kind *CStorPoolCluster) Name() string {
	return kind.handle.Name
}
