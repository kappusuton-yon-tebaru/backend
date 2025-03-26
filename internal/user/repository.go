package user

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
	user *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		user: db.Collection("users"),
	}
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	cur, err := r.user.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in Find:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	users := make([]models.User, 0)

	for cur.Next(ctx) {
		var user UserDTO

		err = cur.Decode(&user)
		if err != nil {
			log.Println("Error in Find:", err)
			return nil, err
		}

		users = append(users, DTOToUser(user))
	}

	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, dto CreateUserDTO) (string, error) {
	user := bson.M{
		"name":     dto.Name,
		"password": dto.Password,
	}

	result, err := r.user.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error inserting user:", err)
		return primitive.NilObjectID.Hex(), fmt.Errorf("error inserting user: %v", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID)

	return insertedID.Hex(), nil
}

func (r *Repository) DeleteUser(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.user.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *Repository) AddRole(ctx context.Context, userID string, roleID string) (string, error) {
	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("userID FromHex err")
		return "", err
	}

	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("roleID FromHex err")
		return "", err
	}

	update := bson.M{
		"$addToSet": bson.M{
			"role_ids": roleObjID, 
		},
	}

	result, err := r.user.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
	if err != nil {
		log.Println("Error adding role to user:", err)
		return "", fmt.Errorf("error adding user to role: %v", err)
	}

	if result.MatchedCount == 0 {
		return "", fmt.Errorf("user not found")
	}

	return userObjID.Hex(), nil
}
func (r *Repository) RemoveRole(ctx context.Context, userID string, roleID string) (int64, error) {
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	update := bson.M{
		"$pull": bson.M{
			"role_ids": roleObjID,
		},
	}

	result, err := r.user.UpdateOne(ctx, bson.M{"_id": userObjID}, update)
	if err != nil {
		log.Println("Error removing role from user:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *Repository) RemoveRoleInAllUser(ctx context.Context, roleID string) (int64, error) {
	roleObjID, err := bson.ObjectIDFromHex(roleID)
	if err != nil {
		log.Println("ObjectID FromHex err")
		return 0, err
	}

	update := bson.M{
		"$pull": bson.M{
			"role_ids": roleObjID,
		},
	}

	result, err := r.user.UpdateMany(ctx, bson.M{"role_ids": roleObjID}, update)
	if err != nil {
		log.Println("Error removing role from users:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}