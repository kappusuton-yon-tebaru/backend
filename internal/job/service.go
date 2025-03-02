package job

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaginatedJobs = models.Paginated[models.Job]

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetAllJobParents(ctx context.Context, queryParam query.QueryParam) (PaginatedJobs, error) {
	dtos, err := s.repo.GetAllJobParents(ctx, queryParam)
	if err != nil {
		return PaginatedJobs{}, err
	}

	jobs := make([]models.Job, 0)
	for _, dto := range dtos.Data {
		numJob := len(dto.Jobs)
		statusMap := map[enum.JobStatus]int{}

		for _, job := range dto.Jobs {
			statusMap[job.JobStatus] += 1
		}

		var finalStatus enum.JobStatus
		for status, count := range statusMap {
			if status == enum.JobStatusPending && count == numJob {
				finalStatus = enum.JobStatusPending
			}

			if status == enum.JobStatusSuccess && count == numJob {
				finalStatus = enum.JobStatusSuccess
			}

			if status == enum.JobStatusRunning && count > 0 {
				finalStatus = enum.JobStatusRunning
			}

			if status == enum.JobStatusFailed && count > 0 {
				finalStatus = enum.JobStatusFailed
			}
		}

		job := models.Job{
			Id:        dto.Id.Hex(),
			CreatedAt: dto.CreatedAt,
			JobStatus: finalStatus,
			Project: models.JobProject{
				Id:   dto.Project.Id.Hex(),
				Name: dto.Project.ResourceName,
			},
		}

		jobs = append(jobs, job)
	}

	return PaginatedJobs{
		Page:  dtos.Page,
		Limit: dtos.Limit,
		Total: dtos.Total,
		Data:  jobs,
	}, nil
}

func (s *Service) GetAllJobsByParentId(ctx context.Context, id string, queryParam query.QueryParam) (models.Paginated[models.Job], *werror.WError) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Paginated[models.Job]{}, werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid parent job id")
	}

	dtos, err := s.repo.GetAllJobsByParentId(ctx, objId, queryParam)
	if err != nil {
		return models.Paginated[models.Job]{}, werror.NewFromError(err)
	}

	jobs := make([]models.Job, 0)
	for _, dto := range dtos.Data {
		jobs = append(jobs, DTOToJob(dto))
	}

	return PaginatedJobs{
		Page:  dtos.Page,
		Limit: dtos.Limit,
		Total: dtos.Total,
		Data:  jobs,
	}, nil
}

func (s *Service) CreateJob(ctx context.Context, dto CreateJobDTO) (string, error) {
	id, err := s.repo.CreateJob(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) CreateGroupJobs(ctx context.Context, dto CreateJobGroupDTO) (CreateGroupJobsResponse, *werror.WError) {
	resp, err := s.repo.CreateGroupJobs(ctx, dto)
	if err != nil {
		return CreateGroupJobsResponse{}, werror.NewFromError(err).SetMessage("error occured while creating jobs")
	}

	return resp, nil
}

func (s *Service) UpdateJobStatus(ctx context.Context, jobId string, jobStatus enum.JobStatus) *werror.WError {
	err := s.repo.UpdateJobStatus(ctx, jobId, string(jobStatus))
	if err != nil {
		return werror.NewFromError(err).SetMessage("error occured while updating job status")
	}

	return nil
}

func (s *Service) DeleteJob(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid job id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteJob(ctx, filter)
	if err != nil {
		return werror.NewFromError(err)
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("not found")
	}

	return nil
}
