package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id             bson.ObjectID 	`bson:"_id"`
	Name           string 			`bson:"name"`
	ProviderType   string 			`bson:"provider_type"` // enum
	JsonCredential string 			`bson:"json_credential"`
	OrganizationId bson.ObjectID 	`bson:"organization_id"`
}

type CreateRegistryProvidersDTO struct {
	Name           string 			`bson:"name"`
	ProviderType   string 			`bson:"provider_type"` // enum
	JsonCredential string 			`bson:"json_credential"`
	OrganizationId bson.ObjectID 	`bson:"organization_id"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	return models.RegistryProviders{
		Id:             regProviders.Id.Hex(),
		Name:           regProviders.Name,
		ProviderType:   regProviders.ProviderType,
		JsonCredential: regProviders.JsonCredential,
		OrganizationId: regProviders.OrganizationId.Hex(),
	}
}