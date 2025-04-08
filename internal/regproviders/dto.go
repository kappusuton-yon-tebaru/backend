package regproviders

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id             bson.ObjectID       `bson:"_id"`
	Name           string              `bson:"name"`
	ProviderType   models.ProviderType `bson:"provider_type"`
	Credential     *models.Credential  `bson:"credential"`
	Uri            string              `bson:"uri"`
	OrganizationId bson.ObjectID       `bson:"organization_id"`
	CreatedAt      time.Time           `bson:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at"`
}

type CreateRegistryProvidersDTO struct {
	Name           string              `bson:"name"                     json:"name"`
	Uri            string              `bson:"uri"                      json:"uri"`
	ProviderType   models.ProviderType `bson:"provider_type"            json:"provider_type"`
	Credential     *models.Credential  `bson:"credential"              json:"credential"`
	OrganizationId bson.ObjectID       `bson:"organization_id"          json:"organization_id"`
	CreatedAt      time.Time           `bson:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at"`
}

type UpdateRegistryProvidersDTO struct {
	Name           string              `bson:"name,omitempty"                     json:"name,omitempty"`
	Uri            string              `bson:"uri,omitempty"                      json:"uri,omitempty"`
	ProviderType   models.ProviderType `bson:"provider_type,omitempty"            json:"provider_type,omitempty"`
	Credential     *models.Credential  `bson:"credential,omitempty"              json:"credential,omitempty"`
	OrganizationId bson.ObjectID       `bson:"organization_id,omitempty"          json:"organization_id,omitempty"`
	UpdatedAt      time.Time           `bson:"updated_at"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	var credential models.Credential
	bsonBytes, _ := bson.Marshal(regProviders.Credential)
	bson.Unmarshal(bsonBytes, &credential)

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
