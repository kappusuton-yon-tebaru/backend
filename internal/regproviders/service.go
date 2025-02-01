package regproviders

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

func (s *Service) GetAllRegistryProviders(ctx context.Context) ([]models.RegistryProviders, error) {
	regProviders, err := s.repo.GetAllRegistryProviders(ctx)
	if err != nil {
		return nil, err
	}

	return regProviders, nil
}

func (s *Service) CreateRegistryProviders(ctx context.Context, dto CreateRegistryProvidersDTO) (string, error) {
	id, err := s.repo.CreateRegistryProviders(ctx, dto)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeleteRegistryProviders(ctx context.Context, id string) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid registry provider id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	count, err := s.repo.DeleteRegistryProviders(ctx, filter)
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