package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id             bson.ObjectID `bson:"_id"`
	Name           string        `bson:"name"`
	ProviderType   string        `bson:"provider_type"` // enum
	JsonCredential string        `bson:"json_credential"`
	Uri            string        `bson:"uri"`
	OrganizationId bson.ObjectID `bson:"organization_id"`
}

type CreateRegistryProvidersDTO struct {
	Name           string        `bson:"name" json:"name"`
	ProviderType   string        `bson:"provider_type" json:"provider_type"` // enum
	Uri            string        `bson:"uri" json:"uri"`
	JsonCredential string        `bson:"json_credential" json:"json_credential"`
	OrganizationId bson.ObjectID `bson:"organization_id" json:"organization_id"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	return models.RegistryProviders{
		Id:             regProviders.Id.Hex(),
		Name:           regProviders.Name,
		ProviderType:   regProviders.ProviderType,
		JsonCredential: regProviders.JsonCredential,
		Uri:            regProviders.Uri,
		OrganizationId: regProviders.OrganizationId.Hex(),
	}
}
