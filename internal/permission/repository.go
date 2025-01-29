package permission

import (
	"context"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository struct {
	permission *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		permission: client.Database("Capstone").Collection("permissions"),
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

func (r *Repository) DeletePermission(ctx context.Context, filter any) (int64, error) {
	result, err := r.permission.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}