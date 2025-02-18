package regproviders

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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

func (r *Repository) GetAllRegistryProviders(ctx context.Context) ([]models.RegistryProviders, error) {
	cur, err := r.regProviders.Find(ctx, bson.D{})
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

func (r *Repository) CreateRegistryProviders(ctx context.Context, dto CreateRegistryProvidersDTO) (string, error) {
	result, err := r.regProviders.InsertOne(ctx, dto)
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
