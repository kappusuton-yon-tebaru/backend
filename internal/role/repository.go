package role

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
	role *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		role: db.Collection("roles"),
	}
}

func (r *Repository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	cur, err := r.role.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	roles := make([]models.Role, 0)

	for cur.Next(ctx) {
		var role RoleDTO

		err = cur.Decode(&role)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		roles = append(roles, DTOToRole(role))
	}

	return roles, nil
}

func (r *Repository) CreateRole(ctx context.Context, dto CreateRoleDTO) (string, error) {
	role := bson.M{
		"role_name": dto.Role_name,
	}

	result, err := r.role.InsertOne(ctx, role)
	if err != nil {
		log.Println("Error inserting role:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting role: %v", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *Repository) DeleteRole(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.role.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
