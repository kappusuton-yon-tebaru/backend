package regproviders

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegistryProvidersDTO struct {
	Id             bson.ObjectID         `bson:"_id"`
	Name           string                `bson:"name"`
	ECRCredential  *models.ECRCredential `bson:"ecr_credential"`
	Uri            string                `bson:"uri"`
	OrganizationId bson.ObjectID         `bson:"organization_id"`
	CreatedAt      time.Time             `bson:"created_at"`
	UpdatedAt      time.Time             `bson:"updated_at"`
}

type CreateRegistryProvidersDTO struct {
	Name           string                `bson:"name"                     json:"name"`
	Uri            string                `bson:"uri"                      json:"uri"`
	ECRCredential  *models.ECRCredential `bson:"ecr_credential,omitempty" json:"ecr_credential,omitempty"`
	OrganizationId bson.ObjectID         `bson:"organization_id"          json:"organization_id"`
	Id             bson.ObjectID         `bson:"_id"`
	CreatedAt      time.Time             `bson:"created_at"`
	UpdatedAt      time.Time             `bson:"updated_at"`
}

type UpdateRegistryProvidersDTO struct {
	Name           string                `bson:"name,omitempty"                     json:"name,omitempty"`
	Uri            string                `bson:"uri,omitempty"                      json:"uri,omitempty"`
	ECRCredential  *models.ECRCredential `bson:"ecr_credential,omitempty"           json:"ecr_credential,omitempty"`
	OrganizationId bson.ObjectID         `bson:"organization_id,omitempty"          json:"organization_id,omitempty"`
	UpdatedAt      time.Time             `bson:"updated_at"`
}

func DTOToRegistryProviders(regProviders RegistryProvidersDTO) models.RegistryProviders {
	return models.RegistryProviders{
		Id:             regProviders.Id.Hex(),
		Name:           regProviders.Name,
		ECRCredential:  regProviders.ECRCredential,
		Uri:            regProviders.Uri,
		OrganizationId: regProviders.OrganizationId.Hex(),
		CreatedAt:      regProviders.CreatedAt,
		UpdatedAt:      regProviders.UpdatedAt,
	}
}
