package roleusergroup

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleUserGroupDTO struct {
	Id        	 	bson.ObjectID 	`bson:"_id"`
	RoleId 	 		bson.ObjectID  	`bson:"role_id"`
	UserGroupId 	bson.ObjectID 	`bson:"user_group_id"`
}

type CreateRoleUserGroupDTO struct {
	Role_id 	 	bson.ObjectID  	`bson:"role_id"`
	UserGroup_id 	bson.ObjectID 	`bson:"user_group_id"`
}

func DTOToRoleUserGroup(RoleUserGroup RoleUserGroupDTO) models.RoleUserGroup {
	return models.RoleUserGroup{
		Id:       		RoleUserGroup.Id.Hex(),
		RoleId:       	RoleUserGroup.RoleId.Hex(),
		UserGroupId:    RoleUserGroup.UserGroupId.Hex(),
	}
}