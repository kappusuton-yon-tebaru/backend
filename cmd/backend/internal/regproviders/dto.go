package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateRegistryProvidersRequest struct {
	Name           string              `json:"name"            validate:"required"`
	Uri            string              `json:"uri"             validate:"required"`
	ProviderType   models.ProviderType `json:"provider_type"   validate:"required"`
	Credential     *models.Credential  `json:"credential"      validate:"required"`
	OrganizationId bson.ObjectID       `json:"organization_id" validate:"required"`
}

type UpdateRegistryProvidersRequest struct {
	Name           string              `json:"name"            validate:"omitempty"`
	Uri            string              `json:"uri"             validate:"omitempty"`
	ProviderType   models.ProviderType `json:"provider_type"   validate:"omitempty"`
	Credential     *models.Credential  `json:"credential"      validate:"omitempty"`
	OrganizationId bson.ObjectID       `json:"organization_id" validate:"omitempty"`
}

type UpdateRegistryProvidersResponse struct {
	Message string `json:"message"`
}

type CreateRegistryProvidersResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
