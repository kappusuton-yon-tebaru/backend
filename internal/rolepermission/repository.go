package rolepermission

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
	rolepermission *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		rolepermission: db.Collection("role_permissions"),
	}
}

func (r *Repository) GetAllRolePermissions(ctx context.Context) ([]models.RolePermission, error) {
	cur, err := r.rolepermission.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	rolepermissions := make([]models.RolePermission, 0)

	for cur.Next(ctx) {
		var rolepermission RolePermissionDTO

		err = cur.Decode(&rolepermission)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		rolepermissions = append(rolepermissions, DTOToRolePermission(rolepermission))
	}

	return rolepermissions, nil
}

func (r *Repository) CreateRolePermission(ctx context.Context, dto CreateRolePermissionDTO) (string, error) {
	rolepermission := bson.M{
		"role_id":       dto.Role_id,
		"permission_id": dto.Permission_id,
	}

	result, err := r.rolepermission.InsertOne(ctx, rolepermission)
	if err != nil {
		log.Println("Error inserting rolepermission:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting rolepermission: %v", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *Repository) DeleteRolePermission(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.rolepermission.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
