package resource

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResourceDTO struct {
	Id           bson.ObjectID     `bson:"_id"`
	ResourceName string            `bson:"resource_name"`
	ResourceType enum.ResourceType `bson:"resource_type"`
}

type CreateResourceDTO struct {
	ResourceName string            `bson:"resource_name"`
	ResourceType enum.ResourceType `bson:"resource_type"`
}

func DTOToResource(resource ResourceDTO) models.Resource {
	return models.Resource{
		Id:           resource.Id.Hex(),
		ResourceName: resource.ResourceName,
		ResourceType: resource.ResourceType,
	}
}
