package user

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDTO struct {
	Id       bson.ObjectID `bson:"_id"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
}

type UserCredentialDTO struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func DTOToUser(user UserDTO) models.User {
	return models.User{
		Id:       user.Id.Hex(),
		Email:    user.Email,
		Password: user.Password,
	}
}
