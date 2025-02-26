package resource

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResourceDTO struct {
	Id           bson.ObjectID     `bson:"_id"`
	ResourceName string            `bson:"resource_name"`
	ResourceType enum.ResourceType `bson:"resource_type"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type CreateResourceDTO struct {
	ResourceName string            `json:"resource_name" bson:"resource_name"`
	ResourceType enum.ResourceType `json:"resource_type" bson:"resource_type"`
}

type UpdateResourceDTO struct {
	ResourceName string            `json:"resource_name" bson:"resource_name"`
}

func DTOToResource(resource ResourceDTO) models.Resource {
	return models.Resource{
		Id:           resource.Id.Hex(),
		ResourceName: resource.ResourceName,
		ResourceType: resource.ResourceType,
		CreatedAt:    resource.CreatedAt,
		UpdatedAt:    resource.UpdatedAt,
	}
}
