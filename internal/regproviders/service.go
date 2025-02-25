package regproviders

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
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

func ParseCredential(provider enum.RegistryProviderType, dto map[string]any) (interface{}, *werror.WError) {
	switch provider {
	case enum.ECR:
		var ecr models.ECRCredential
		if err := utils.MapToStruct(dto, &ecr); err != nil {
			return nil, werror.NewFromError(err).
				SetMessage("cannot parse ecr credential").
				SetCode(http.StatusBadRequest)
		}
		return ecr, nil

	default:
		return nil, werror.New().
			SetMessage("invalid provider type").
			SetCode(http.StatusBadRequest)
	}
}
