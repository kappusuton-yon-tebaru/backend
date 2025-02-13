package projectrepository

import (
	"context"
	"fmt"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	projRepo *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		projRepo: client.Database("Capstone").Collection("projects_repositories"),
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

func (r *Repository) GetProjectRepositoryByFilter(ctx context.Context, filter map[string]any) (models.ProjectRepository, error) {
	var projRepoDTO ProjectRepositoryDTO

	err := r.projRepo.FindOne(ctx, filter).Decode(&projRepoDTO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document found, return an empty Resource
			return models.ProjectRepository{}, nil
		}
		log.Println("Error finding project repository:", err)
		return models.ProjectRepository{}, err
	}

	return DTOToProjectRepository(projRepoDTO), nil
}

func (r *Repository) CreateProjectRepository(ctx context.Context, dto CreateProjectRepositoryDTO) (string, error) {
	// Check if project_id already exists
	filter := bson.M{"project_id": dto.Project_Id}
	count, err := r.projRepo.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error checking existing project repository:", err)
		return "", fmt.Errorf("error checking existing project repository: %v", err)
	}

	if count > 0 {
		return "", fmt.Errorf("project repository already exists")
	}
	
	projRepo := bson.M{
		"git_repo_url": dto.Git_Repo_Url,
		"project_id":   dto.Project_Id,
	}

	result, err := r.projRepo.InsertOne(ctx, projRepo)
	if err != nil {
		log.Println("Error inserting project repository:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting project repository: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) DeleteProjectRepository(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.projRepo.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
