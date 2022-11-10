package kube

import (
	apiv1 "k8s.io/api/core/v1"
)

func NewPodSpec(name, tag, cmd string, args ...string) apiv1.PodSpec {
	return apiv1.PodSpec{
		Containers: []apiv1.Container{
			NewContainer(name, tag, cmd, args...),
		},
	}
}
