package permission

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
)

type PermissionDTO struct {
	Id        			bson.ObjectID 			`bson:"_id"`
	Permission_name 	string        			`bson:"permission_name"`
	Action 				enum.PermissionActions 	`bson:"action"`
	ResourceId			bson.ObjectID			`bson:"resource_id"`
	ResourceType		enum.ResourceType 		`bson:"resource_type"`
}

func DTOToPermission(permission PermissionDTO) models.Permission {
	return models.Permission{
		Id:       			permission.Id.Hex(),
		Permission_name:  	permission.Permission_name,
		Action:  			permission.Action,
		ResourceId:  		permission.ResourceId.Hex(),
		ResourceType:  		permission.ResourceType,
	}
}