package deploy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/internal/aws"
	sharedDeploy "github.com/kappusuton-yon-tebaru/backend/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
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

func (s *Service) DeployService(ctx context.Context, req DeployRequest) *werror.WError {
	projectId, err := bson.ObjectIDFromHex(req.ProjectId)
	if err != nil {
		return werror.NewFromError(err).SetMessage("invalid project id").SetCode(http.StatusBadRequest)
	}

	projRepo, werr := s.projectRepoService.GetProjectRepositoryByProjectId(ctx, req.ProjectId)
	if werr != nil {
		return werr
	}

	if projRepo.RegistryProvider == nil || len(strings.TrimSpace(projRepo.RegistryProvider.Uri)) == 0 {
		return werror.New().
			SetMessage("registry uri cannot be empty").
			SetCode(http.StatusBadRequest)
	}

	secretManager, err := aws.NewSession(*projRepo.RegistryProvider.ECRCredential, "ap-southeast-1").SecretsManager()
	if err != nil {
		s.logger.Error("error occured while creating aws session", zap.Error(err))
		return werror.NewFromError(err)
	}

	envs := make(map[string]string)
	err = secretManager.GetSecretValue(ctx, "covid-summary", &envs)
	if err != nil {
		s.logger.Error("error occured while creating aws secret manager session", zap.Error(err))
		return werror.NewFromError(err)
	}

	jobs := []job.CreateJobDTO{}
	for _, service := range req.Services {
		jobs = append(jobs, job.CreateJobDTO{
			JobType:     string(enum.JobTypeDeploy),
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
		return werr
	}

	for i, service := range req.Services {
		jobId := resp.JobIds[i]

		deployCtx := sharedDeploy.DeployContext{
			Id:           jobId,
			ServiceName:  service.ServiceName,
			ImageUri:     fmt.Sprintf("%s:%s", projRepo.RegistryProvider.Uri, service.Tag),
			Environments: envs,
		}

		fmt.Println(deployCtx)
	}

	return nil
}
