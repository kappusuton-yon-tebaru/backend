package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	apicorev1 "k8s.io/api/core/v1"
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
	success := false

	defer func() {
		var status enum.JobStatus
		if success {
			status = enum.JobStatusSuccess
		} else {
			status = enum.JobStatusFailed
		}

		ctx := context.Background()
		werr := s.jobService.UpdateJobStatus(ctx, dto.Id, status)
		if werr != nil {
			s.logger.Error("error occured while updating job status", zap.Error(werr.Err), zap.String("job_id", dto.Id))
		}
	}()

	deployClient := s.kube.NewDeploymentClient(dto.Namespace)
	deployManifest := kubernetes.ApplyDeploymentManifest(dto)

	deployedService, err := deployClient.Apply(ctx, deployManifest)
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

	for {
		deployedService, err := deployClient.Get(ctx, deployedService.Name)
		if err != nil {
			s.logger.Error("error occured while getting deployment status", zap.Any("manifest", deployManifest), zap.Error(err))
			return werror.NewFromError(err)
		}

		condition := deployClient.GetCondition(deployedService)

		replicaDeplyed := condition.Progressing != nil &&
			condition.Progressing.Reason == enum.ProgressingReasonNewReplicaSetAvailable

		replicaMatched := deployedService.Spec.Replicas == nil ||
			(deployedService.Spec.Replicas != nil && deployedService.Status.ReadyReplicas == *deployedService.Spec.Replicas)

		available := condition.Available != nil &&
			condition.Available.Status == apicorev1.ConditionTrue

		if replicaDeplyed && available && replicaMatched {
			success = true
			break
		}

		time.Sleep(time.Second)
	}

	fmt.Println("DONE")

	return nil
}
