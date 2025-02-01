package roleusergroup

import (
	"context"
	"fmt"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository struct {
	roleUsergroup *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		roleUsergroup: client.Database("Capstone").Collection("role_user_groups"),
	}
}

func (r *Repository) GetAllRoleUserGroups(ctx context.Context) ([]models.RoleUserGroup, error) {
	cur, err := r.roleUsergroup.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	roleUserGroups := make([]models.RoleUserGroup, 0)

	for cur.Next(ctx) {
		var roleUserGroup RoleUserGroupDTO

		err = cur.Decode(&roleUserGroup)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		roleUserGroups = append(roleUserGroups, DTOToRoleUserGroup(roleUserGroup))
	}

	return roleUserGroups, nil
}

func (r *Repository) CreateRoleUserGroup(ctx context.Context, dto CreateRoleUserGroupDTO) (any, error) {
	roleuserGroup := bson.M{
		"role_id": dto.Role_id,
		"user_group_id":dto.UserGroup_id,
	}

	result, err := r.roleUsergroup.InsertOne(ctx, roleuserGroup)
	if err != nil {
		log.Println("Error inserting roleUserGroup:", err)
		return primitive.NilObjectID, fmt.Errorf("error inserting roleUserGroup: %v", err)
	}

	return result.InsertedID, nil
}

func (r *Repository) DeleteRoleUserGroup(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.roleUsergroup.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}