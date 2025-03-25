package role

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	RoleName string        `bson:"role_name"`
	OrgId  bson.ObjectID 	`bson:"org_id"`
	Permissions []models.Permission `bson:"permissions"`
}

type CreateRoleDTO struct {
	RoleName string 		`json:"role_name" bson:"role_name"`
	OrgId  bson.ObjectID 	`json:"org_id" bson:"org_id"`
}

type UpdateRoleDTO struct {
	RoleName string 		`json:"role_name" bson:"role_name"`
}

func DTOToRole(role RoleDTO) models.Role {
	return models.Role{
		Id:        role.Id.Hex(),
		Role_name: role.RoleName,
		OrgId: role.OrgId.Hex(),
		Permissions: role.Permissions,
	}
}
