package rolepermission

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RolePermissionDTO struct {
	Id           bson.ObjectID `bson:"_id"`
	RoleId       bson.ObjectID `bson:"role_id"`
	PermissionId bson.ObjectID `bson:"permission_id"`
}

type CreateRolePermissionDTO struct {
	Role_id       bson.ObjectID `bson:"role_id"`
	Permission_id bson.ObjectID `bson:"permission_id"`
}

func DTOToRolePermission(rolepermission RolePermissionDTO) models.RolePermission {
	return models.RolePermission{
		Id:           rolepermission.Id.Hex(),
		RoleId:       rolepermission.RoleId.Hex(),
		PermissionId: rolepermission.PermissionId.Hex(),
	}
}
