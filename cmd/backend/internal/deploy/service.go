package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/kappusuton-yon-tebaru/backend/internal/aws"
	sharedDeploy "github.com/kappusuton-yon-tebaru/backend/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
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
	resourceService    *resource.Service
}

func NewService(rmq *rmq.BuilderRmq, jobService *job.Service, logger *logger.Logger, projectRepoService *projectrepository.Service, resourceService *resource.Service) *Service {
	return &Service{
		rmq,
		jobService,
		logger,
		projectRepoService,
		resourceService,
	}
}

func (s *Service) DeployService(ctx context.Context, req DeployRequest) (string, *werror.WError) {
	projectId, err := bson.ObjectIDFromHex(req.ProjectId)
	if err != nil {
		return "", werror.NewFromError(err).SetMessage("invalid project id").SetCode(http.StatusBadRequest)
	}

	project, werr := s.resourceService.GetResourceByID(ctx, req.ProjectId)
	if werr != nil {
		return "", werr
	}

	projRepo, werr := s.projectRepoService.GetProjectRepositoryByProjectId(ctx, req.ProjectId)
	if werr != nil {
		return "", werr
	}

	if projRepo.RegistryProvider == nil || len(strings.TrimSpace(projRepo.RegistryProvider.Uri)) == 0 {
		return "", werror.New().
			SetMessage("registry uri cannot be empty").
			SetCode(http.StatusBadRequest)
	}

	secretManager := aws.NewConfig(*projRepo.RegistryProvider.ECRCredential, "ap-southeast-1").SecretsManager()

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
		return "", werr
	}

	for i, service := range req.Services {
		jobId := resp.JobIds[i]

		envs := make(map[string]string)
		if service.SecretName != nil {
			err := secretManager.GetSecretValue(ctx, *service.SecretName, &envs)

			var notFoundErr *types.ResourceNotFoundException
			if errors.As(err, &notFoundErr) {
				return "", werror.NewFromError(err).SetMessage("secret not found").SetCode(http.StatusNotFound)
			} else if err != nil {
				s.logger.Error("error occured while creating aws secret manager session", zap.Error(err))
				return "", werror.NewFromError(err)
			}
		}

		deployCtx := sharedDeploy.DeployContext{
			Id:            jobId,
			ProjectId:     req.ProjectId,
			ServiceName:   service.ServiceName,
			ImageUri:      fmt.Sprintf("%s:%s", projRepo.RegistryProvider.Uri, service.Tag),
			Port:          service.Port,
			Namespace:     deployenv.GetNamespaceName(project.ResourceName, req.DeploymentEnv),
			Environments:  envs,
			DeploymentEnv: req.DeploymentEnv,
			HealthCheck:   (*sharedDeploy.DeployHealthCheckContext)(service.Healthcheck),
		}

		bs, err := json.Marshal(deployCtx)
		if err != nil {
			return "", nil
		}

		fmt.Println(string(bs))

		if err := s.rmq.Publish(ctx, enum.DeployContextRoutingKey, bs); err != nil {
			s.logger.Error("error occured while publishing deploy context", zap.Error(err))
			return "", werror.NewFromError(err).SetMessage("error occured while publishing deploy context")
		}
	}

	return resp.ParentId, nil
}
