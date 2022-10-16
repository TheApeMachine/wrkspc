package kube

import (
	apiv1 "k8s.io/api/core/v1"
)

func NewContainer(name string) apiv1.Container {
	return apiv1.Container{
		Name:  name,
		Image: "theapemachine/name:release",
		Ports: []apiv1.ContainerPort{
			{
				Name:          "http",
				Protocol:      apiv1.ProtocolTCP,
				ContainerPort: 80,
			},
		},
	}
}
