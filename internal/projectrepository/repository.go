package projectrepository

import (
	"context"
	"errors"
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

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		projRepo: db.Collection("projects_repositories"),
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

func (r *Repository) GetProjectRepositoryByFilter(ctx context.Context, filter map[string]any) (ProjectRepositoryDTO, error) {
	pipeline := []map[string]any{
		{
			"$match": filter,
		},
		{
			"$lookup": map[string]any{
				"from":         "registry_providers",
				"localField":   "registry_provider_id",
				"foreignField": "_id",
				"as":           "registry_provider",
			},
		},
		{
			"$unwind": map[string]any{
				"path": "$registry_provider",
			},
		},
		{
			"$project": map[string]any{
				"registry_provider_id": false,
			},
		},
		{
			"$limit": 1,
		},
	}

	cur, err := r.projRepo.Aggregate(ctx, pipeline)
	if err != nil {
		return ProjectRepositoryDTO{}, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return ProjectRepositoryDTO{}, errors.New("not found")
	}

	var projectRepo ProjectRepositoryDTO
	err = cur.Decode(&projectRepo)
	if err != nil {
		return ProjectRepositoryDTO{}, err
	}

	return projectRepo, nil
}

func (r *Repository) CreateProjectRepository(ctx context.Context, dto CreateProjectRepositoryDTO) (string, error) {
	// Check if project_id already exists
	filter := bson.M{"project_id": dto.ProjectId}
	count, err := r.projRepo.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error checking existing project repository:", err)
		return "", fmt.Errorf("error checking existing project repository: %v", err)
	}

	if count > 0 {
		return "", fmt.Errorf("project repository already exists")
	}

	projRepo := bson.M{
		"git_repo_url": dto.GitRepoUrl,
		"project_id":   dto.ProjectId,
	}

	result, err := r.projRepo.InsertOne(ctx, projRepo)
	if err != nil {
		log.Println("Error inserting project repository:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting project repository: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) UpdateProjectRepository(ctx context.Context, projectId bson.ObjectID, registryProviderId bson.ObjectID) (int64, error) {
	filter := map[string]any{
		"project_id": projectId,
	}

	update := map[string]any{
		"$set": map[string]any{
			"registry_provider_id": registryProviderId,
		},
	}

	result, err := r.projRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *Repository) DeleteProjectRepository(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.projRepo.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
