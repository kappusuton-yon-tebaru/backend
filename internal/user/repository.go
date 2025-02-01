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

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		user: client.Database("Capstone").Collection("users"),
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

func (r *Repository) DeleteUser(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.user.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
