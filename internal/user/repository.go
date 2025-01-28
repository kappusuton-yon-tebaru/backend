package user

import (
	"context"
	"log"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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

func (r *Repository) DeleteUser(ctx context.Context, filter any) (int64, error) {
	result, err := r.user.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
