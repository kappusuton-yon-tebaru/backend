package projectenv

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	repo *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		repo: client.Database("Capstone").Collection("projectenv"),
	}
}

func (r *Repository) GetAllProjectEnvs(ctx context.Context) ([]models.ProjectEnv, error) {
	cur, err := r.repo.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	projectenvs := make([]models.ProjectEnv, 0)

	for cur.Next(ctx) {
		var dto ProjectEnvDTO

		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		projectenvs = append(projectenvs, DTOToProjectEnv(dto))
	}

	return projectenvs, nil
}

func (r *Repository) CreateProjectEnv(ctx context.Context, dto CreateProjectEnvDTO) (string, error) {
	result, err := r.repo.InsertOne(ctx, dto)
	if err != nil {
		return primitive.NilObjectID.Hex(), err
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) DeleteProjectEnv(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.repo.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
