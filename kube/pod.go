package kube

import (
	"context"

	"github.com/theapemachine/wrkspc/errnie"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	client RestClient
	name   string
	conn   *apiv1.Pod
}

func NewPod(ctx context.Context, client RestClient, name string) Pod {
	errnie.Traces()

	conn, err := client.conn.CoreV1().Pods(v1.NamespaceDefault).Create(ctx, &apiv1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{{
				Name:    name,
				Command: []string{"/bin/zsh", "-c", "--"},
				Args:    []string{"while true; do sleep 30; done;"},
				Image:   "theapemachine/" + name,
			}},
		},
	}, v1.CreateOptions{})

	errnie.Handles(err).With(errnie.KILL)

	if err != nil {
		conn, err = client.conn.CoreV1().Pods(v1.NamespaceDefault).Get(ctx, name, v1.GetOptions{})
		errnie.Handles(err).With(errnie.KILL)
	}

	return Pod{
		client: client,
		name:   name,
		conn:   conn,
	}
}

func (pod Pod) Running() bool {
	errnie.Traces()
	return pod.conn.Status.Phase == apiv1.PodRunning
}
