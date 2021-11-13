package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	schedulingv1 "k8s.io/api/scheduling/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
PriorityClass ...
*/
type PriorityClass struct {
	base *Base
	spec *schedulingv1.PriorityClass
}

/*
NewPriorityClass ...
*/
func NewPriorityClass(client RestClient, file []byte) MigratableKind {
	errnie.Traces()
	kind := &PriorityClass{}
	kind.base = NewBase(file, kind, client)
	return kind
}

/*
Up ...
*/
func (kind *PriorityClass) Up() {
	errnie.Handles(yaml.Unmarshal(kind.base.file, &kind.spec)).With(errnie.KILL)

	kind.spec, kind.base.err = kind.base.client.conn.SchedulingV1().PriorityClasses().Create(
		context.Background(), kind.spec, v1.CreateOptions{},
	)

	kind.base.waiter()
}

/*
Check ...
*/
func (kind *PriorityClass) Check() bool {
	check, err := kind.base.client.conn.SchedulingV1().PriorityClasses().Get(
		context.Background(), kind.spec.Name, v1.GetOptions{},
	)
	return check != nil && err == nil
}

/*
Delete ...
*/
func (kind *PriorityClass) Delete() error {
	return kind.base.client.conn.SchedulingV1().PriorityClasses().Delete(
		context.TODO(), kind.spec.Name, v1.DeleteOptions{},
	)
}
