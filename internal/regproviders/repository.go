package regproviders

import (
	"context"
	"errors"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	regProviders *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		regProviders: db.Collection("registry_providers"),
	}
}

func (r *Repository) GetAllRegistryProviders(ctx context.Context, queryParam query.QueryParam) (models.Paginated[RegistryProvidersDTO], error) {
	pipeline := utils.NewFilterAggregationPipeline(queryParam, []map[string]any{})

	cur, err := r.regProviders.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Paginated[RegistryProvidersDTO]{}, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return models.Paginated[RegistryProvidersDTO]{}, errors.New("not found")
	}

	var dto models.Paginated[RegistryProvidersDTO]
	err = cur.Decode(&dto)
	if err != nil {
		return models.Paginated[RegistryProvidersDTO]{}, err
	}

	return dto, nil
}

func (r *Repository) GetAllRegistryProvidersWithoutProject(ctx context.Context) ([]models.RegistryProviders, error) {
	pipeline := []map[string]any{
		{
			"$lookup": map[string]any{
				"from":         "projects_repositories",
				"localField":   "_id",
				"foreignField": "registry_provider_id",
				"as":           "projects",
			},
		},
		{
			"$match": map[string]any{
				"projects": map[string]any{
					"$size": 0,
				},
			},
		},
	}

	cur, err := r.regProviders.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	registryProviders := make([]models.RegistryProviders, 0)

	for cur.Next(ctx) {
		var dto RegistryProvidersDTO

		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		registryProviders = append(registryProviders, DTOToRegistryProviders(dto))
	}

	return registryProviders, nil
}

func (r *Repository) GetRegistryProviderById(ctx context.Context, filter map[string]any) (models.RegistryProviders, error) {
	var dto RegistryProvidersDTO

	err := r.regProviders.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		return models.RegistryProviders{}, err
	}

	return DTOToRegistryProviders(dto), nil
}

func (r *Repository) CreateRegistryProviders(ctx context.Context, dto CreateRegistryProvidersDTO) (string, error) {
	request := CreateRegistryProvidersDTO{
		Name:           dto.Name,
		ProviderType:   dto.ProviderType,
		Uri:            dto.Uri,
		Credential:     dto.Credential,
		OrganizationId: dto.OrganizationId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := r.regProviders.InsertOne(ctx, request)
	if err != nil {
		return primitive.NilObjectID.Hex(), err
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) DeleteRegistryProviders(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.regProviders.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
