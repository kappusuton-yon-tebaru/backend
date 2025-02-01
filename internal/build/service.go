package build

import (
	"context"
	"fmt"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	apicorev1 "k8s.io/api/core/v1"
)

type Service struct {
	namespace string
	kube      *kubernetes.Kubernetes
}

func NewService(cfg *config.Config, kube *kubernetes.Kubernetes) *Service {
	return &Service{
		cfg.KubeNamespace,
		kube,
	}
}

func (s *Service) BuildImage(ctx context.Context, dto kubernetes.BuildImageDTO) error {
	podClient := s.kube.NewPodClient(s.namespace)
	manifest := kubernetes.CreateBuilderPodManifest(dto)
	builderPod, err := podClient.Create(ctx, manifest)
	if err != nil {
		return err
	}

	for {
		pod, err := podClient.Get(ctx, builderPod.Name)
		if err != nil {
			return err
		}

		if pod.Status.Phase == apicorev1.PodFailed || pod.Status.Phase == apicorev1.PodUnknown {
			fmt.Println("pod failed")

			podClient.GetLogString(ctx, builderPod.Name)

			break
		} else if pod.Status.Phase == apicorev1.PodSucceeded {
			fmt.Println("pod succeeded")
			break
		}

		time.Sleep(time.Second)
	}

	err = podClient.Delete(ctx, builderPod.Name)
	if err != nil {
		return err
	}

	return nil
}
