package regproviders

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PaginatedRegistryProviders = models.Paginated[models.RegistryProviders]

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetAllRegistryProviders(ctx context.Context, queryParam query.QueryParam) (PaginatedRegistryProviders, error) {
	dtos, err := s.repo.GetAllRegistryProviders(ctx, queryParam)
	if err != nil {
		return PaginatedRegistryProviders{}, err
	}

	regProviders := make([]models.RegistryProviders, 0)
	for _, dto := range dtos.Data {
		regProviders = append(regProviders, DTOToRegistryProviders(dto))
	}

	return PaginatedRegistryProviders{
		Limit: dtos.Limit,
		Page:  dtos.Page,
		Total: dtos.Total,
		Data:  regProviders,
	}, nil
}

func (s *Service) GetAllRegistryProvidersWithoutProject(ctx context.Context) ([]models.RegistryProviders, error) {
	regProviders, err := s.repo.GetAllRegistryProvidersWithoutProject(ctx)
	if err != nil {
		return nil, err
	}

	return regProviders, nil
}

func (s *Service) GetRegistryProviderById(ctx context.Context, id string) (models.RegistryProviders, *werror.WError) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.RegistryProviders{}, werror.New().
			SetCode(http.StatusBadRequest).
			SetMessage("invalid registry provider id")
	}

	filter := map[string]any{
		"_id": objId,
	}

	regProvider, err := s.repo.GetRegistryProviderById(ctx, filter)

	if err != nil {
		return models.RegistryProviders{}, werror.NewFromError(err)
	}

	return regProvider, nil
}

func (s *Service) CreateRegistryProviders(ctx context.Context, dto CreateRegistryProvidersDTO) (string, *werror.WError) {
	id, err := s.repo.CreateRegistryProviders(ctx, dto)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return "", werror.NewFromError(err).SetMessage("registry name already exist in this organization").SetCode(http.StatusBadRequest)
	} else if err != nil {
		return "", werror.NewFromError(err)
	}

	return id, nil
}

func (s *Service) UpdateRegistryProviders(ctx context.Context, id string, dto UpdateRegistryProvidersDTO) *werror.WError {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusBadRequest).
			SetMessage("invalid registry provider id")
	}

	count, err := s.repo.UpdateRegistryProviders(ctx, objId, dto)
	if err != nil {
		return werror.NewFromError(err).
			SetCode(http.StatusInternalServerError).
			SetMessage("failed to update registry provider")
	}

	if count == 0 {
		return werror.New().
			SetCode(http.StatusNotFound).
			SetMessage("not found")
	}

	return nil
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
