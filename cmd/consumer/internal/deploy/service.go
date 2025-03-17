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
	deployManifest := kubernetes.ApplyDeploymentManifest(dto)

	_, err := deployClient.Apply(ctx, deployManifest)
	if err != nil {
		s.logger.Error("error occured while deploying service", zap.Any("manifest", deployManifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	svcManifest := kubernetes.ApplyLoadBalancerService(dto)
	svcClient := s.kube.NewServiceClient(dto.Namespace)

	_, err = svcClient.Apply(ctx, svcManifest)
	if err != nil {
		s.logger.Error("error occured while deploying service", zap.Any("manifest", deployManifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	return nil
}
