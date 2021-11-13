package kube

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/theapemachine/wrkspc/errnie"

	schedulingv1 "k8s.io/api/scheduling/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PriorityClass struct {
	base   Base
	client RestClient
	spec   *schedulingv1.PriorityClass
	file   []byte
	handle *schedulingv1.PriorityClass
}

func NewPriorityClass(client RestClient, file []byte) MigratableKind {
	return &PriorityClass{
		client: client,
		file:   file,
	}
}

func (kind *PriorityClass) Up() error {
	err := yaml.Unmarshal(kind.file, &kind.spec)
	errnie.Handles(err).With(errnie.KILL)

	kind.handle, err = kind.client.conn.SchedulingV1().PriorityClasses().Create(context.Background(), kind.spec, v1.CreateOptions{})
	errnie.Handles(err).With(errnie.KILL)

	kind.base = NewBase(kind)
	kind.base.waiter(true)

	return err
}

func (kind *PriorityClass) Check() bool {
	check, _ := kind.client.conn.SchedulingV1().PriorityClasses().Get(context.Background(), kind.spec.Name, v1.GetOptions{})
	return check != nil
}

func (kind *PriorityClass) Down() error {
	return kind.base.teardown()
}

func (kind *PriorityClass) Delete() error {
	return kind.client.conn.SchedulingV1().PriorityClasses().Delete(context.TODO(), kind.handle.Name, v1.DeleteOptions{})
}

func (kind *PriorityClass) Name() string {
	return kind.handle.Name
}
