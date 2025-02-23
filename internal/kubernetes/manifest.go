package kubernetes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	acappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	accorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	acmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func CreateBuilderPodManifest(config BuildImageDTO) *apicorev1.Pod {
	var setupArgs []string
	var mountPath string

	switch credential := config.Credential.(type) {
	case models.ECRCredential:
		mountPath = "/root/.aws"
		setupArgs = []string{
			"mkdir -p /root/.aws && tee /root/.aws/credentials << EOF",
			"[default]",
			fmt.Sprintf("aws_access_key_id=%s", credential.AccessKey),
			fmt.Sprintf("aws_secret_access_key=%s", credential.SecretAccessKey),
			"EOF",
		}
	}

	builderArgs := []string{
		fmt.Sprintf("--dockerfile=%s", config.Dockerfile),
		fmt.Sprintf("--context=%s", config.RepoUrl),
		fmt.Sprintf("--context-sub-path=%s", config.RepoRoot),
	}

	for _, dest := range config.Destinations {
		builderArgs = append(builderArgs, fmt.Sprintf("--destination=%s", dest))
	}

	return &apicorev1.Pod{
		ObjectMeta: apimetav1.ObjectMeta{
			Name: "worker-" + config.Id,
		},
		Spec: apicorev1.PodSpec{
			Containers: []apicorev1.Container{
				{
					Name:    "setup",
					Image:   "busybox",
					Command: []string{"/bin/sh", "-c"},
					Args:    []string{strings.Join(setupArgs, "\n")},
					VolumeMounts: []apicorev1.VolumeMount{
						{
							Name:      "credentials",
							MountPath: mountPath,
						},
					},
				},
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:latest",
					Args:  builderArgs,
					VolumeMounts: []apicorev1.VolumeMount{
						{
							Name:      "credentials",
							MountPath: mountPath,
						},
					},
				},
			},
			RestartPolicy: apicorev1.RestartPolicyNever,
			Volumes: []apicorev1.Volume{
				{
					Name: "credentials",
					VolumeSource: apicorev1.VolumeSource{
						EmptyDir: nil,
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
								accorev1.EnvVar().WithName("IN_CLUSTER").WithValue(strconv.FormatBool(config.InCluster)),
								accorev1.EnvVar().WithName("DEVELOPMENT").WithValue(strconv.FormatBool(config.Development)),
								accorev1.EnvVar().WithName("KUBE_NAMESPACE").WithValue(config.KubeNamespace),
								accorev1.EnvVar().WithName("MONGO_URI").WithValue(config.MongoUri),
								accorev1.EnvVar().WithName("MONGO_DATABASE_NAME").WithValue(config.MongoDatabaseName),
								accorev1.EnvVar().WithName("AGENT_PORT").WithValue(strconv.FormatInt(int64(config.Agent.Port), 10)),
								accorev1.EnvVar().WithName("AGENT_WORKER_IMAGE_URI").WithValue(config.Agent.WorkerImageUri),
								accorev1.EnvVar().WithName("BACKEND_PORT").WithValue(strconv.FormatInt(int64(config.Backend.Port), 10)),
								accorev1.EnvVar().WithName("BACKEND_AGENT_ENDPOINT").WithValue(config.Backend.AgentEndpoint),
								accorev1.EnvVar().WithName("BUILDER_QUEUE_URI").WithValue(config.BuilderConfig.QueueUri),
								accorev1.EnvVar().WithName("BUILDER_QUEUE_NAME").WithValue(config.BuilderConfig.QueueName),
							),
					),
				),
			),
		)
}
