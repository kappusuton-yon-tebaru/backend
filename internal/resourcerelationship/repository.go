package resourcerelationship

import (
	"context"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	resourceRela *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		resourceRela: client.Database("Capstone").Collection("resource_relationships"),
	}
}

func (r *Repository) GetAllResourceRelationships(ctx context.Context) ([]models.ResourceRelationship, error) {
	cur, err := r.resourceRela.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	resourceRelas := make([]models.ResourceRelationship, 0)

	for cur.Next(ctx) {
		var resourceRela ResourceRelationshipDTO

		err = cur.Decode(&resourceRela)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		resourceRelas = append(resourceRelas, DTOToResourceRelationship(resourceRela))
	}

	return resourceRelas, nil
}

func (r *Repository) DeleteResourceRelationship(ctx context.Context, filter any) (int64, error) {
	result, err := r.resourceRela.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
