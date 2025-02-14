package kubernetes

import (
	"fmt"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	acappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	accorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	acmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"strconv"
)

func CreateBuilderPodManifest(config BuildImageDTO) *apicorev1.Pod {
	args := []string{
		fmt.Sprintf("--dockerfile=%s", config.Dockerfile),
		fmt.Sprintf("--context=%s", config.RepoUrl),
		fmt.Sprintf("--context-sub-path=%s", config.RepoRoot),
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

func ApplyBuilderConsumerDeploymentManifest(dto ConfigureMaxWorkerDTO, config *config.Config) *acappsv1.DeploymentApplyConfiguration {
	return acappsv1.Deployment("builder-consumer-deployment", config.KubeNamespace).
		WithNamespace("default").
		WithLabels(map[string]string{"app": "builder-consumer"}).
		WithSpec(acappsv1.DeploymentSpec().
			WithReplicas(dto.WorkerNumber).
			WithSelector(acmetav1.LabelSelector().
				WithMatchLabels(map[string]string{"app": "builder-consumer"})).
			WithTemplate(accorev1.PodTemplateSpec().
				WithLabels(map[string]string{"app": "builder-consumer"}).
				WithSpec(accorev1.PodSpec().
					WithServiceAccountName("system").
					WithContainers(
						accorev1.Container().
							WithName("builder-consumer").
							WithImage(dto.WorkerImageUri).
							WithEnv(
								accorev1.EnvVar().WithName("DEVELOPMENT").WithValue(strconv.FormatBool(config.Development)),
								accorev1.EnvVar().WithName("IN_CLUSTER").WithValue("true"),
								accorev1.EnvVar().WithName("AGENT_PORT").WithValue(strconv.FormatInt(int64(config.Agent.Port), 10)),
								accorev1.EnvVar().WithName("BACKEND_PORT").WithValue(strconv.FormatInt(int64(config.Backend.Port), 10)),
								accorev1.EnvVar().WithName("MONGO_URI").WithValue("mongodb+srv://kappusutonyontebaru:K%40pPu%24Ut0N_yOnte6arUu@capstone.5t8hk.mongodb.net/?retryWrites=true&w=majority&appName=Capstone"),
								accorev1.EnvVar().WithName("BUILDER_QUEUE_URI").WithValue("amqps://ctvdziew:ORU88dANQWlvLo6joDUDz13hxXxHNvka@armadillo.rmq.cloudamqp.com/ctvdziew"),
								accorev1.EnvVar().WithName("BUILDER_QUEUE_NAME").WithValue("org.builder"),
								accorev1.EnvVar().WithName("KUBE_NAMESPACE").WithValue("default"),
							),
					),
				),
			),
		)
}
