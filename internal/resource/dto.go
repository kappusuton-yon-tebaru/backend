package resource

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResourceDTO struct {
	Id           bson.ObjectID     `json:"_id"           bson:"_id"`
	ResourceName string            `json:"resource_name" bson:"resource_name "`
	ResourceType enum.ResourceType `json:"resource_type" bson:"resource_type "`
}

type CreateResourceDTO struct {
	ResourceName string            `json:"resource_name" bson:"resource_name"`
	ResourceType enum.ResourceType `json:"resource_type" bson:"resource_type"`
}

func DTOToResource(resource ResourceDTO) models.Resource {
	return models.Resource{
		Id:           resource.Id.Hex(),
		ResourceName: resource.ResourceName,
		ResourceType: resource.ResourceType,
	}
}
