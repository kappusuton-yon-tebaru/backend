package resourcerelationship

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResourceRelationshipDTO struct {
	Id               bson.ObjectID `bson:"_id"`
	ParentResourceId bson.ObjectID `bson:"parent_resource_id"`
	ChildResourceId  bson.ObjectID `bson:"child_resource_id"`
}

func DTOToResourceRelationship(resourceRela ResourceRelationshipDTO) models.ResourceRelationship {
	return models.ResourceRelationship{
		Id:               resourceRela.Id.Hex(),
		ParentResourceId: resourceRela.ParentResourceId.Hex(),
		ChildResourceId:  resourceRela.ChildResourceId.Hex(),
	}
}
