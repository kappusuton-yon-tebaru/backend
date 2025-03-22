package build

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sharedBuild "github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	projectId, err := bson.ObjectIDFromHex(req.ProjectId)
	if err != nil {
		return "", werror.NewFromError(err).SetMessage("invalid project id").SetCode(400)
	}

	projRepo, werr := s.projectRepoService.GetProjectRepositoryByProjectId(ctx, req.ProjectId)
	if werr != nil {
		return "", werr
	}

	if len(strings.TrimSpace(projRepo.GitRepoUrl)) == 0 {
		return "", werror.New().
			SetMessage("git repository url cannot be empty").
			SetCode(http.StatusBadRequest)
	}

	if projRepo.RegistryProvider == nil || len(strings.TrimSpace(projRepo.RegistryProvider.Uri)) == 0 {
		return "", werror.New().
			SetMessage("registry uri cannot be empty").
			SetCode(http.StatusBadRequest)
	}

	jobs := []job.CreateJobDTO{}
	for _, service := range req.Services {
		jobs = append(jobs, job.CreateJobDTO{
			JobType:     string(enum.JobTypeBuild),
			JobStatus:   string(enum.JobStatusPending),
			ProjectId:   projectId,
			ServiceName: service.ServiceName,
		})
	}

	dto := job.CreateJobGroupDTO{
		ProjectId: projectId,
		Jobs:      jobs,
	}

	resp, werr := s.jobService.CreateGroupJobs(ctx, dto)
	if werr != nil {
		s.logger.Error("error occured while creating jobs", zap.Error(werr.Err))
		return "", werr
	}

	repoUrl, err := projRepo.GetGitRepoUrl()
	if err != nil {
		return "", werror.NewFromError(err).SetMessage("invalid git repo url").SetCode(http.StatusBadRequest)
	}

	for i, service := range req.Services {
		jobId := resp.JobIds[i]

		buildCtx := sharedBuild.BuildContext{
			Id:            jobId,
			RepoUrl:       repoUrl,
			RepoRoot:      fmt.Sprintf("apps/%s", service.ServiceName),
			Destination:   fmt.Sprintf("%s:%s", projRepo.RegistryProvider.Uri, utils.ToKebabCase(service.Tag)),
			Dockerfile:    "Dockerfile",
			ECRCredential: projRepo.RegistryProvider.ECRCredential,
		}

		bs, err := json.Marshal(buildCtx)
		if err != nil {
			return "", nil
		}

		if err := s.rmq.Publish(ctx, enum.BuildContextRoutingKey, bs); err != nil {
			s.logger.Error("error occured while publishing build context", zap.Error(err))
			return "", werror.NewFromError(err).SetMessage("error occured while publishing build context")
		}
	}

	return resp.ParentId, nil
}
