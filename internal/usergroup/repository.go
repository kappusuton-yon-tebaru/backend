package usergroup

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
	usergroup           *mongo.Collection
	usergroupMembership *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		usergroup:           client.Database("Capstone").Collection("user_groups"),
		usergroupMembership: client.Database("Capstone").Collection("user_group_memberships"),
	}
}

func (r *Repository) GetAllUserGroups(ctx context.Context) ([]models.UserGroup, error) {
	cur, err := r.usergroup.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	usergroups := make([]models.UserGroup, 0)

	for cur.Next(ctx) {
		var usergroup UserGroupDTO

		err = cur.Decode(&usergroup)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		usergroups = append(usergroups, DTOToUserGroup(usergroup))
	}

	return usergroups, nil
}

func (r *Repository) CreateUserGroup(ctx context.Context, dto CreateUserGroupDTO) (any, error) {
	usergroup := bson.M{
		"group_name": dto.GroupName,
	}

	result, err := r.usergroup.InsertOne(ctx, usergroup)
	if err != nil {
		log.Println("Error inserting user group:", err)
		return primitive.NilObjectID, fmt.Errorf("error inserting user group: %v", err)
	}

	insertedID := result.InsertedID

	return insertedID, nil
}

func (r *Repository) AddUserToUserGroup(ctx context.Context, dto AddUserToUserGroupDTO, id bson.ObjectID) (any, error) {
	addUser := bson.M{
		"group_id": id,
		"user_id":  dto.UserId,
	}

	result, err := r.usergroupMembership.InsertOne(ctx, addUser)
	if err != nil {
		log.Println("Error inserting user to user group:", err)
		return primitive.NilObjectID, fmt.Errorf("error inserting user to user group: %v", err)
	}

	insertedID := result.InsertedID

	return insertedID, nil
}

func (r *Repository) DeleteUserGroup(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.usergroup.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *Repository) DeleteUserFromUserGroup(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.usergroupMembership.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
