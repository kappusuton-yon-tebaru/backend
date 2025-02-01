package role

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	Role_name string        `bson:"role_name"`
}

type CreateRoleDTO struct {
	Role_name string `bson:"role_name"`
}

func DTOToRole(role RoleDTO) models.Role {
	return models.Role{
		Id:        role.Id.Hex(),
		Role_name: role.Role_name,
	}
}
