package matrix

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/kube"
)

/*
Container ...
*/
type Container struct {
	cli  *client.Client
	ctx  context.Context
	name string
}

/*
NewContainer ...
*/
func NewContainer(name string) Container {
	errnie.Traces()

	return Container{
		ctx:  context.Background(),
		name: name,
	}
}

/*
Run ...
*/
func (container Container) Run() error {
	errnie.Traces()

	client := kube.NewRestClient()
	pod := kube.NewPod(container.ctx, client, container.name)

	for {
		if pod.Running() {
			break
		}
	}

	return nil
}
