package kubernetes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	acappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	accorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	acmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func CreateBuilderPodManifest(config BuildImageDTO) (*apicorev1.Pod, error) {
	var setupArgs []string
	var mountPath string

	if config.ECRCredential != nil {
		mountPath = "/root/.aws"
		setupArgs = []string{
			"mkdir -p /root/.aws && tee /root/.aws/credentials << EOF",
			"[default]",
			fmt.Sprintf("aws_access_key_id=%s", config.ECRCredential.AccessKey),
			fmt.Sprintf("aws_secret_access_key=%s", config.ECRCredential.SecretAccessKey),
			"EOF",
		}
	} else {
		return nil, errors.New("invalid credential")
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
	}, nil
}

func ApplyBuilderConsumerDeploymentManifest(dto ConfigureMaxWorkerDTO, config *config.Config) *acappsv1.DeploymentApplyConfiguration {
	return acappsv1.Deployment("consumer-deployment", config.KubeNamespace).
		WithNamespace(config.KubeNamespace).
		WithLabels(map[string]string{"app": "consumer"}).
		WithSpec(acappsv1.DeploymentSpec().
			WithReplicas(dto.WorkerNumber).
			WithSelector(acmetav1.LabelSelector().
				WithMatchLabels(map[string]string{"app": "consumer"})).
			WithTemplate(accorev1.PodTemplateSpec().
				WithLabels(map[string]string{"app": "consumer"}).
				WithSpec(accorev1.PodSpec().
					WithServiceAccountName("system").
					WithContainers(
						accorev1.Container().
							WithName("consumer").
							WithImage(dto.WorkerImageUri).
							WithEnv(
								accorev1.EnvVar().WithName("IN_CLUSTER").WithValue(strconv.FormatBool(config.InCluster)),
								accorev1.EnvVar().WithName("DEVELOPMENT").WithValue(strconv.FormatBool(config.Development)),
								accorev1.EnvVar().WithName("KUBE_NAMESPACE").WithValue(config.KubeNamespace),
								accorev1.EnvVar().WithName("MONGO_URI").WithValue(config.MongoUri),
								accorev1.EnvVar().WithName("MONGO_DATABASE_NAME").WithValue(config.MongoDatabaseName),
								accorev1.EnvVar().WithName("AGENT_PORT").WithValue(strconv.FormatInt(int64(config.Agent.Port), 10)),
								accorev1.EnvVar().WithName("WORKER_IMAGE_URI").WithValue(config.Agent.WorkerImageUri),
								accorev1.EnvVar().WithName("BACKEND_PORT").WithValue(strconv.FormatInt(int64(config.Backend.Port), 10)),
								accorev1.EnvVar().WithName("BACKEND_AGENT_ENDPOINT").WithValue(config.Backend.AgentEndpoint),
								accorev1.EnvVar().WithName("CONSUMER_QUEUE_URI").WithValue(config.ConsumerConfig.QueueUri),
								accorev1.EnvVar().WithName("CONSUMER_ORGANIZATION_NAME").WithValue(config.ConsumerConfig.OrganizationName),
							),
					),
				),
			),
		)
}

func ApplyDeploymentManifest(dto DeployDTO) *acappsv1.DeploymentApplyConfiguration {
	port := accorev1.ContainerPort()
	if dto.Port != nil {
		port = port.WithContainerPort(int32(*dto.Port))
	}

	envs := []*accorev1.EnvVarApplyConfiguration{}
	for key, val := range dto.Environments {
		envs = append(envs, accorev1.EnvVar().WithName(key).WithValue(val))
	}

	return acappsv1.Deployment(fmt.Sprintf("%s-deployment", dto.ServiceName), dto.Namespace).
		WithLabels(map[string]string{
			"app":            dto.ServiceName,
			"project_id":     dto.ProjectId,
			"service_name":   dto.ServiceName,
			"deployment_env": dto.DeploymentEnv,
		}).
		WithSpec(acappsv1.DeploymentSpec().
			WithReplicas(1).
			WithSelector(acmetav1.LabelSelector().
				WithMatchLabels(map[string]string{"app": dto.ServiceName})).
			WithRevisionHistoryLimit(0).
			WithTemplate(accorev1.PodTemplateSpec().
				WithLabels(map[string]string{"app": dto.ServiceName}).
				WithSpec(accorev1.PodSpec().
					WithServiceAccountName("system").
					WithContainers(
						accorev1.Container().
							WithName(dto.ServiceName).
							WithPorts(port).
							WithImage(dto.ImageUri).
							WithEnv(envs...).
							WithReadinessProbe(accorev1.Probe().
								WithHTTPGet(&accorev1.HTTPGetActionApplyConfiguration{
									Path: utils.Pointer("/hc"),
									Port: utils.Pointer(intstr.FromInt32(8080)),
								}).
								WithInitialDelaySeconds(5).
								WithPeriodSeconds(5).
								WithTimeoutSeconds(1)),
					),
				),
			),
		)
}

func ApplyLoadBalancerService(dto DeployDTO) *accorev1.ServiceApplyConfiguration {
	return accorev1.Service(fmt.Sprintf("%s-service", dto.ServiceName), dto.Namespace).
		WithSpec(&accorev1.ServiceSpecApplyConfiguration{
			Selector: map[string]string{
				"app": dto.ServiceName,
			},
			Ports: []accorev1.ServicePortApplyConfiguration{
				{
					Protocol:   utils.Pointer(apicorev1.ProtocolTCP),
					Port:       dto.Port,
					TargetPort: utils.Pointer(intstr.FromInt32(*dto.Port)),
				},
			},
		}).WithLabels(map[string]string{
		"project_id":     dto.ProjectId,
		"service_name":   dto.ServiceName,
		"deployment_env": dto.DeploymentEnv,
	})
}
