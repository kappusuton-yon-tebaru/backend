package build

import (
	"context"
	"net/http"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type Service struct {
	namespace  string
	kube       *kubernetes.Kubernetes
	logger     *logger.Logger
	jobService *job.Service
}

func NewService(cfg *config.Config, kube *kubernetes.Kubernetes, logger *logger.Logger, jobService *job.Service) *Service {
	return &Service{
		cfg.KubeNamespace,
		kube,
		logger,
		jobService,
	}
}

func (s *Service) BuildImage(ctx context.Context, dto kubernetes.BuildImageDTO) *werror.WError {
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

	podClient := s.kube.NewPodClient(s.namespace)
	manifest, err := kubernetes.CreateBuilderPodManifest(dto)
	if err != nil {
		s.logger.Error("error occured while creating buidler pod", zap.Error(err))
		return werror.NewFromError(err)
	}

	builderPod, err := podClient.Create(ctx, manifest)
	if err != nil {
		s.logger.Error("error occured while creating buidler pod", zap.Any("manifest", manifest), zap.Error(err))
		return werror.NewFromError(err)
	}

	werr := s.jobService.UpdateJobStatus(ctx, dto.Id, enum.JobStatusRunning)
	if werr != nil {
		s.logger.Error("error occured while updating job status", zap.Error(werr.Err), zap.String("job_id", dto.Id))
		return werr
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
			log, err := podClient.GetLog(builderPod.Name, kubernetes.WithContainer("kaniko")).String(ctx)
			if err != nil {
				s.logger.Error("error occured while log", zap.String("name", builderPod.Name), zap.Error(err))
				return werror.NewFromError(err)
			}

			s.logger.Error(log)

			break
		} else if pod.Status.Phase == apicorev1.PodSucceeded {
			success = true
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
