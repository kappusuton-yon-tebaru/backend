package user

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDTO struct {
	Id       bson.ObjectID `bson:"_id"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
	RoleIds []bson.ObjectID `bson:"role_ids"`
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
		RoleIds:  mapRoles(user.RoleIds),
	}
}

func mapRoles(roles []bson.ObjectID) []string {
	mapped := make([]string, len(roles))
	for i, role := range roles {
		mapped[i] = role.Hex()
	}
	return mapped
}
