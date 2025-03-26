package resource

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	resource     *mongo.Collection
	resourceRela *mongo.Collection
	projectRepo  *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		resource:     db.Collection("resources"),
		resourceRela: db.Collection("resource_relationships"),
		projectRepo:  db.Collection("projects_repositories"),
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

func (r *Repository) GetResourceByFilter(ctx context.Context, filter map[string]any) (models.Resource, error) {
	var resource ResourceDTO

	err := r.resource.FindOne(ctx, filter).Decode(&resource)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document found, return an empty Resource
			return models.Resource{}, nil
		}
		log.Println("Error in FindOne:", err)
		return models.Resource{}, err
	}

	// Convert DTO to the actual model
	return DTOToResource(resource), nil
}

func (r *Repository) GetResourcesByFilter(ctx context.Context, queryParam query.QueryParam, parentId string) (models.Paginated[ResourceDTO], error) {
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return models.Paginated[ResourceDTO]{}, err
	}
	pipeline := utils.NewFilterAggregationPipeline(queryParam,
		[]map[string]any{
			{
				"$lookup": map[string]any{
					"from":         "resource_relationships",
					"localField":   "_id",
					"foreignField": "child_resource_id",
					"as":           "relationships",
				},
			},
			{
				"$unwind": map[string]any{
					"path": "$relationships",
				},
			},
			{
				"$match": map[string]any{
					"relationships.parent_resource_id": objId,
				},
			},
		},
	)

	cur, err := r.resource.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Paginated[ResourceDTO]{}, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return models.Paginated[ResourceDTO]{}, errors.New("not found")
	}

	var dto models.Paginated[ResourceDTO]
	err = cur.Decode(&dto)
	if err != nil {
		return models.Paginated[ResourceDTO]{}, err
	}

	return dto, nil
}

func (r *Repository) CreateResource(ctx context.Context, dto CreateResourceDTO) (string, error) {
	if !enum.IsValidResourceType(dto.ResourceType) {
		return "", fmt.Errorf("invalid resource type: %v", dto.ResourceType)
	}

	resource := bson.M{
		"resource_name": dto.ResourceName,
		"resource_type": dto.ResourceType,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	result, err := r.resource.InsertOne(ctx, resource)
	if err != nil {
		log.Println("Error inserting resource:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting resource: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) UpdateResource(ctx context.Context, dto UpdateResourceDTO, resourceID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(resourceID)
	if err != nil {
		log.Println("ObjectIDFromHex err")
		return "", err
	}
	update := map[string]any{
		"$set": map[string]any{
			"resource_name": dto.ResourceName,
			"updated_at":    time.Now(), // Ensure `updated_at` is refreshed
		},
	}
	// Update the resource in MongoDB
	result, err := r.resource.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error updating resource:", err)
		return "", fmt.Errorf("error updating resource: %v", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		return "", fmt.Errorf("resource not found")
	}

	return objID.Hex(), nil
}

func (r *Repository) DeleteResource(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.resource.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// Cascade delete function for organization → project_space → project
func (r *Repository) CascadeDeleteResource(ctx context.Context, resourceID string, resourceType enum.ResourceType) error {

	objId, err := bson.ObjectIDFromHex(resourceID)
	if err != nil {
		log.Println("ObjectIDFromHex err")

		return err
	}

	findChildFilter := map[string]any{
		"parent_resource_id": objId,
	}

	// Find all direct children of this resource
	cur, err := r.resourceRela.Find(ctx, findChildFilter)
	if err != nil {
		log.Println("Error in Find:", err)
		return err
	}

	defer cur.Close(ctx)

	childRelationships := make([]models.ResourceRelationship, 0)

	for cur.Next(ctx) {
		var resourceRela resourcerelationship.ResourceRelationshipDTO

		err = cur.Decode(&resourceRela)
		if err != nil {
			log.Println("Error in Find2:", err)
			return err
		}

		childRelationships = append(childRelationships, resourcerelationship.DTOToResourceRelationship(resourceRela))
	}

	// Recursively delete all child resources
	for _, rel := range childRelationships {
		// Get child resource type before deleting (needed for repo deletion)
		filter := map[string]any{
			"_id": rel.ChildResourceId,
		}
		childResource, err := r.GetResourceByFilter(ctx, filter)
		if err != nil {
			log.Println("resource.FindOne err")

			return err // Return if the child resource is not found
		}

		if err := r.CascadeDeleteResource(ctx, rel.ChildResourceId, childResource.ResourceType); err != nil {
			log.Println("r.CascadeDeleteResource err")

			return err
		}
	}

	// If the resource is a PROJECT, delete related repositories
	if resourceType == enum.ResourceTypeProject {
		_, err := r.projectRepo.DeleteOne(ctx, map[string]any{"project_id": objId})
		if err != nil {
			log.Println("projectRepo.DeleteOne err")

			return err
		}
		log.Println("Deleted repositories for project:", objId)
	}

	// Delete all relationships where the resource is parent or child
	_, err = r.resourceRela.DeleteMany(ctx, map[string]any{
		"$or": []map[string]any{
			{"parent_resource_id": objId},
			{"child_resource_id": objId},
		},
	})
	if err != nil {
		log.Println("r.resourceRela.DeleteMany err")

		return err
	}

	// Delete the resource itself
	_, err = r.resource.DeleteOne(ctx, map[string]any{"_id": objId})
	if err != nil {
		log.Println("r.resource.DeleteOne err")

		return err
	}

	log.Println("Deleted resource:", objId)
	return nil
}
