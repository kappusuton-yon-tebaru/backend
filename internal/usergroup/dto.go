package usergroup

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserGroupDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	GroupName string        `bson:"group_name"`
}

type CreateUserGroupDTO struct {
	GroupName string `json:"group_name" bson:"group_name"`
}

func DTOToUserGroup(usergroup UserGroupDTO) models.UserGroup {
	return models.UserGroup{
		Id:        usergroup.Id.Hex(),
		GroupName: usergroup.GroupName,
	}
}
