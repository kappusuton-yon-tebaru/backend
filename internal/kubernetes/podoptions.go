package kubernetes

import apicorev1 "k8s.io/api/core/v1"

type PodLogOptions func(opt *apicorev1.PodLogOptions)

func WithContainer(container string) PodLogOptions {
	return func(opt *apicorev1.PodLogOptions) {
		opt.Container = container
	}
}
