package deploy

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
)

type Service struct {
	kube       *kubernetes.Kubernetes
	logger     *logger.Logger
	jobService *job.Service
}

func NewService(kube *kubernetes.Kubernetes, logger *logger.Logger, jobService *job.Service) *Service {
	return &Service{
		kube,
		logger,
		jobService,
	}
}

func (s *Service) Deploy(ctx context.Context, dto kubernetes.DeployDTO) *werror.WError {
	deployClient := s.kube.NewDeploymentClient(dto.Namespace)
	manifest := kubernetes.ApplyDeploymentManifest(dto)

	_, err := deployClient.Apply(ctx, manifest)
	if err != nil {
		s.logger.Error("error occured while deploying service", zap.Any("manifest", manifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	return nil
}
