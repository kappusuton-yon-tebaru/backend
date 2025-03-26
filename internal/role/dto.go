package role

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PermDTO struct {
	Id            bson.ObjectID `bson:"_id"`
	PermissionName string       `bson:"permission_name"`
	Action        enum.PermissionActions        `bson:"action"`
	ResourceId    bson.ObjectID `bson:"resource_id"`
}
type RoleDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	RoleName string        `bson:"role_name"`
	OrgId  bson.ObjectID 	`bson:"org_id"`
	Permissions []PermDTO `bson:"permissions"`
}

type CreateRoleDTO struct {
	RoleName string 		`json:"role_name" bson:"role_name"`
	OrgId  bson.ObjectID 	`json:"org_id" bson:"org_id"`
}

type UpdateRoleDTO struct {
	RoleName string 		`json:"role_name" bson:"role_name"`
}

type CreatePermissionDTO struct {
	PermissionName string                 `json:"permission_name" bson:"permission_name"`
	Action          enum.PermissionActions `json:"action" bson:"action"`
	ResourceId     bson.ObjectID          `json:"resource_id" bson:"resource_id"`
}

func DTOToRole(role RoleDTO) models.Role {
	return models.Role{
		Id:        role.Id.Hex(),
		RoleName: role.RoleName,
		OrgId: role.OrgId.Hex(),
		Permissions: mapPermissions(role.Permissions),
	}
}

func mapPermissions(perms []PermDTO) []models.Perm {
	mapped := make([]models.Perm, len(perms))
	for i, perm := range perms {
		mapped[i] = models.Perm{
			Id:            perm.Id.Hex(), // Convert ObjectID to string
			PermissionName: perm.PermissionName,
			Action:        perm.Action,
			ResourceId:    perm.ResourceId.Hex(), // Convert ObjectID to string
		}
	}
	return mapped
}
