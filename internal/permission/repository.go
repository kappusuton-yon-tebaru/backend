package permission

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
	permission *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		permission: db.Collection("permissions"),
	}
}

func (r *Repository) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	cur, err := r.permission.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	permissions := make([]models.Permission, 0)

	for cur.Next(ctx) {
		var permission PermissionDTO

		err = cur.Decode(&permission)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		permissions = append(permissions, DTOToPermission(permission))
	}

	return permissions, nil
}

func (r *Repository) CreatePermission(ctx context.Context, dto CreatePermissionDTO) (string, error) {
	permission := bson.M{
		"permission_name": dto.Permission_name,
		"action":          dto.Action,
		"resource_id":     dto.Resource_id,
		"resource_type":   dto.Resource_type,
	}

	result, err := r.permission.InsertOne(ctx, permission)

	if err != nil {
		log.Println("Error inserting permission:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting permission: %v", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *Repository) DeletePermission(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.permission.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
