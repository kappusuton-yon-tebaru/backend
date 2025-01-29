package job

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	jobs, err := s.repo.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *Service) CreateJob(ctx context.Context, dto CreateJobDTO) (any, error) {
	id, err := s.repo.CreateJob(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
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