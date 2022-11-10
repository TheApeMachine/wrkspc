package kube

import (
	apiv1 "k8s.io/api/core/v1"
)

func NewContainer(name, tag, cmd string, args ...string) apiv1.Container {
	return apiv1.Container{
		Name:    name,
		Image:   tag,
		Command: append([]string{"./wrkspc", cmd}, args...),
		Ports: []apiv1.ContainerPort{
			{
				Name:          "http",
				Protocol:      apiv1.ProtocolTCP,
				ContainerPort: 80,
			},
		},
	}
}
