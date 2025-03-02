package user

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDTO struct {
	Id       bson.ObjectID `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
}

type RegisterDTO struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func DTOToUser(user UserDTO) models.User {
	return models.User{
		Id:       user.Id.Hex(),
		Username: user.Username,
		Password: user.Password,
	}
}
