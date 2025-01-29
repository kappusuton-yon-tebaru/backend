package user

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDTO struct {
	Id       bson.ObjectID `bson:"_id"`
	Name     string        `bson:"name"`
	Password string        `bson:"password"`
}

type CreateUserDTO struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

func DTOToUser(user UserDTO) models.User {
	return models.User{
		Id:       user.Id.Hex(),
		Name:     user.Name,
		Password: user.Password,
	}
}
