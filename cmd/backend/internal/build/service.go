package build

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sharedBuild "github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
)

type Service struct {
	rmq                *rmq.BuilderRmq
	jobService         *job.Service
	logger             *logger.Logger
	projectRepoService *projectrepository.Service
}

func NewService(rmq *rmq.BuilderRmq, jobService *job.Service, logger *logger.Logger, projectRepoService *projectrepository.Service) *Service {
	return &Service{
		rmq,
		jobService,
		logger,
		projectRepoService,
	}
}

func (s *Service) BuildImage(ctx context.Context, req BuildRequest) (string, *werror.WError) {
	projRepo, werr := s.projectRepoService.GetProjectRepositoryByProjectId(ctx, req.ProjectId)
	if werr != nil {
		return "", werr
	}

	if len(strings.TrimSpace(projRepo.GitRepoUrl)) == 0 {
		return "", werror.New().SetMessage("git repository url cannot be empty")
	}

	if len(strings.TrimSpace(projRepo.RegistryProvider.Uri)) == 0 {
		return "", werror.New().SetMessage("registry uri cannot be empty")
	}

	jobs := []job.CreateJobDTO{}
	for range len(req.Services) {
		jobs = append(jobs, job.CreateJobDTO{
			JobType:   string(enum.JobTypeBuild),
			JobStatus: string(enum.JobStatusPending),
		})
	}

	resp, werr := s.jobService.CreateGroupJobs(ctx, jobs)
	if werr != nil {
		s.logger.Error("error occured while creating jobs", zap.Error(werr.Err))
		return "", werr
	}

	for i, service := range req.Services {
		jobId := resp.JobIds[i]

		buildCtx := sharedBuild.BuildContext{
			Id:                 jobId,
			RepoUrl:            projRepo.GitRepoUrl,
			RepoRoot:           fmt.Sprintf("apps/%s", service.ServiceName),
			Destination:        fmt.Sprintf("%s:%s", projRepo.RegistryProvider.Uri, service.Tag),
			Dockerfile:         "Dockerfile",
			RegistryType:       projRepo.RegistryProvider.ProviderType,
			RegistryCredential: projRepo.RegistryProvider.Credential,
		}

		bs, err := json.Marshal(buildCtx)
		if err != nil {
			return "", nil
		}

		if err := s.rmq.Publish(ctx, bs); err != nil {
			s.logger.Error("error occured while publishing build context", zap.Error(err))
			return "", werror.NewFromError(err).SetMessage("error occured while publishing build context")
		}
	}

	return resp.ParentId, nil
}
