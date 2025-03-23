package regproviders

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id                  bson.ObjectID               `bson:"_id"`
	Name                string                      `bson:"name"`
	ProviderType        string                      `bson:"provider_type"`
	ECRCredential       *models.ECRCredential       `bson:"ecr_credential"`
	DockerhubCredential *models.DockerhubCredential `bson:"dockerhub_credential"`
	Uri                 string                      `bson:"uri"`
	OrganizationId      bson.ObjectID               `bson:"organization_id"`
	CreatedAt           time.Time                   `bson:"created_at"`
	UpdatedAt           time.Time                   `bson:"updated_at"`
}

type CreateRegistryProvidersDTO struct {
	Name                string                      `bson:"name"                     json:"name"`
	Uri                 string                      `bson:"uri"                      json:"uri"`
	ProviderType        string                      `bson:"provider_type"            json:"provider_type"`
	ECRCredential       *models.ECRCredential       `bson:"ecr_credential,omitempty" json:"ecr_credential,omitempty"`
	DockerhubCredential *models.DockerhubCredential `bson:"dockerhub_credential,omitempty" json:"dockerhub_credential,omitempty"`
	OrganizationId      bson.ObjectID               `bson:"organization_id"          json:"organization_id"`
	CreatedAt           time.Time                   `bson:"created_at"`
	UpdatedAt           time.Time                   `bson:"updated_at"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	return models.RegistryProviders{
		Id:                  regProviders.Id.Hex(),
		Name:                regProviders.Name,
		ProviderType:        regProviders.ProviderType,
		ECRCredential:       regProviders.ECRCredential,
		DockerhubCredential: regProviders.DockerhubCredential,
		Uri:                 regProviders.Uri,
		OrganizationId:      regProviders.OrganizationId.Hex(),
		CreatedAt:           regProviders.CreatedAt,
		UpdatedAt:           regProviders.UpdatedAt,
	}
}
