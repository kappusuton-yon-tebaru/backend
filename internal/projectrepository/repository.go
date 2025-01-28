package projectrepository

import (
	"context"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	projRepo *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		projRepo: client.Database("Capstone").Collection("project_repositories"),
	}
}

func (r *Repository) GetAllProjectRepositories(ctx context.Context) ([]models.ProjectRepository, error) {
	cur, err := r.projRepo.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	projRepos := make([]models.ProjectRepository, 0)

	for cur.Next(ctx) {
		var projrepo ProjectRepositoryDTO

		err = cur.Decode(&projrepo)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		projRepos = append(projRepos, DTOToProjectRepository(projrepo))
	}

	return projRepos, nil
}

func (r *Repository) DeleteProjectRepository(ctx context.Context, filter any) (int64, error) {
	result, err := r.projRepo.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
