package regproviders

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id             bson.ObjectID             `bson:"_id"`
	Name           string                    `bson:"name"`
	ProviderType   enum.RegistryProviderType `bson:"provider_type"` // enum
	Credential     interface{}               `bson:"credential"`
	Uri            string                    `bson:"uri"`
	OrganizationId bson.ObjectID             `bson:"organization_id"`
	CreatedAt	  	time.Time      	`bson:"created_at"`
	UpdatedAt	  	time.Time      	`bson:"updated_at"`

}

type CreateRegistryProvidersDTO struct {
	Name           string                    `bson:"name"`
	ProviderType   enum.RegistryProviderType `bson:"provider_type"`
	Uri            string                    `bson:"uri"`
	Credential     interface{}               `bson:"credential"`
	OrganizationId bson.ObjectID             `bson:"organization_id"`
	Id             	bson.ObjectID 	`bson:"_id"`
	CreatedAt	  	time.Time      	`bson:"created_at"`
	UpdatedAt	  	time.Time      	`bson:"updated_at"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	return models.RegistryProviders{
		Id:             regProviders.Id.Hex(),
		Name:           regProviders.Name,
		ProviderType:   regProviders.ProviderType,
		Credential:     regProviders.Credential,
		Uri:            regProviders.Uri,
		OrganizationId: regProviders.OrganizationId.Hex(),
		CreatedAt:      regProviders.CreatedAt,
		UpdatedAt:      regProviders.UpdatedAt,
	}
}
