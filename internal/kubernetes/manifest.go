package kubernetes

import (
	"fmt"

	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateBuilderPodManifest(config BuildImageDTO) *apicorev1.Pod {
	args := []string{
		fmt.Sprintf("--dockerfile=%s", config.Dockerfile),
		fmt.Sprintf("--context=%s", config.Url),
	}

	for _, dest := range config.Destinations {
		args = append(args, fmt.Sprintf("--destination=%s", dest))
	}

	return &apicorev1.Pod{
		ObjectMeta: apimetav1.ObjectMeta{
			Name: "worker-" + config.Id,
		},
		Spec: apicorev1.PodSpec{
			Containers: []apicorev1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:latest",
					Args:  args,
					VolumeMounts: []apicorev1.VolumeMount{
						{
							Name:      "aws-secret",
							MountPath: "/root/.aws",
						},
					},
				},
			},
			RestartPolicy: apicorev1.RestartPolicyNever,
			Volumes: []apicorev1.Volume{
				{
					Name: "aws-secret",
					VolumeSource: apicorev1.VolumeSource{
						Secret: &apicorev1.SecretVolumeSource{
							SecretName: "aws-secret",
						},
					},
				},
			},
		},
	}
}
