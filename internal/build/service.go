package build

import (
	"context"
	"net/http"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type Service struct {
	namespace string
	kube      *kubernetes.Kubernetes
	logger    *logger.Logger
}

func NewService(cfg *config.Config, kube *kubernetes.Kubernetes, logger *logger.Logger) *Service {
	return &Service{
		cfg.KubeNamespace,
		kube,
		logger,
	}
}

func (s *Service) BuildImage(ctx context.Context, dto kubernetes.BuildImageDTO) *werror.WError {
	podClient := s.kube.NewPodClient(s.namespace)
	manifest := kubernetes.CreateBuilderPodManifest(dto)
	builderPod, err := podClient.Create(ctx, manifest)
	if err != nil {
		s.logger.Error("error occured while creating buidler pod", zap.Any("manifest", manifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	for {
		pod, err := podClient.Get(ctx, builderPod.Name)
		if errors.IsNotFound(err) {
			return werror.
				NewFromError(err).
				SetCode(http.StatusNotFound).
				SetMessage("pod not found")
		} else if err != nil {
			s.logger.Error("error occured while getting pod", zap.String("name", builderPod.Name), zap.Error(err))
			return werror.NewFromError(err)
		}

		if pod.Status.Phase == apicorev1.PodFailed ||
			pod.Status.Phase == apicorev1.PodUnknown {
			s.logger.Error("pod failed to start", zap.String("name", builderPod.Name))
			// podClient.GetLogString(ctx, builderPod.Name)
			break
		} else if pod.Status.Phase == apicorev1.PodSucceeded {
			break
		}

		time.Sleep(time.Second)
	}

	err = podClient.Delete(ctx, builderPod.Name)
	if err != nil {
		s.logger.Error("error occured while deleting pod", zap.String("name", builderPod.Name), zap.Error(err))
		return werror.NewFromError(err)
	}

	return nil
}
