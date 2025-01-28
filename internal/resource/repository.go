package resource

import (
	"context"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	resource *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		resource: client.Database("Capstone").Collection("resources"),
	}
}

func (r *Repository) GetAllResources(ctx context.Context) ([]models.Resource, error) {
	cur, err := r.resource.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	resources := make([]models.Resource, 0)

	for cur.Next(ctx) {
		var resource ResourceDTO

		err = cur.Decode(&resource)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		resources = append(resources, DTOToResource(resource))
	}

	return resources, nil
}

func (r *Repository) DeleteResource(ctx context.Context, filter any) (int64, error) {
	result, err := r.resource.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
