package resourcerelationship

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

func (r *Repository) GetChildrenResourceRelationshipByParentID(ctx context.Context, filter map[string]any) ([]models.ResourceRelationship, error) {	
	cur, err := r.resourceRela.Find(ctx, filter)
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}
	defer cur.Close(ctx)

	childrenResources := make([]models.ResourceRelationship, 0)

	for cur.Next(ctx) {
		var childrenRe ResourceRelationshipDTO

		err = cur.Decode(&childrenRe)
		if err != nil {
			log.Println("Error in Decode:", err)
			return nil, err
		}

		childrenResources = append(childrenResources, DTOToResourceRelationship(childrenRe))
	}

	if err := cur.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return childrenResources, nil
}

func (r *Repository) CreateResourceRelationship(ctx context.Context, dto CreateResourceRelationshipDTO) (string, error) {
	resourceRela := bson.M{
		"parent_resource_id": dto.Parent_Resource_Id,
		"child_resource_id":  dto.Child_Resource_Id,
	}

	result, err := r.resourceRela.InsertOne(ctx, resourceRela)
	if err != nil {
		log.Println("Error inserting resource relationship:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting resource relationship: %v", err)
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) DeleteResourceRelationship(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.resourceRela.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
