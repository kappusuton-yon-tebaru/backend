package setting

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
)

type Service struct {
	logger *logger.Logger
	config *config.Config
	kube   *kubernetes.Kubernetes
}

func NewService(config *config.Config, kube *kubernetes.Kubernetes, logger *logger.Logger) *Service {
	return &Service{
		logger,
		config,
		kube,
	}
}

func (s *Service) SetMaxWorker(ctx context.Context, numberWorker int32) *werror.WError {
	dto := kubernetes.ConfigureMaxWorkerDTO{
		WorkerNumber:   numberWorker,
		WorkerImageUri: s.config.Agent.WorkerImageUri,
	}

	deployClient := s.kube.NewDeploymentClient(s.config.KubeNamespace)
	manifest := kubernetes.ApplyBuilderConsumerDeploymentManifest(dto, s.config)
	_, err := deployClient.Apply(ctx, manifest)
	if err != nil {
		s.logger.Error("error occured while configurating max worker", zap.Any("manifest", manifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	return nil
}
