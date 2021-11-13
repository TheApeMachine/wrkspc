package kube

import (
	"context"

	"github.com/ghodss/yaml"
	openebsapiv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/theapemachine/wrkspc/errnie"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
CStorPoolCluster ...
*/
type CStorPoolCluster struct {
	base   *Base
	spec   *openebsapiv1.CStorPoolCluster
	client OpenEBSClient
}

/*
NewCStorPoolCluster ...
*/
func NewCStorPoolCluster(client OpenEBSClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &CStorPoolCluster{client: client}
	kind.base = NewBase(file, kind, RestClient{})
	return kind
}

/*
Up ...
*/
func (kind *CStorPoolCluster) Up() {
	errnie.Traces()
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *CStorPoolCluster) Check() bool {
	check, _ := kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil
}

/*
Delete ...
*/
func (kind *CStorPoolCluster) Delete() error {
	return kind.client.conn.CstorV1().CStorPoolClusters(kind.spec.Namespace).Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
