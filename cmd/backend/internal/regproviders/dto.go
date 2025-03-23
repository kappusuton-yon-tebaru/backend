package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateRegistryProvidersRequest struct {
	Name                string                      `json:"name"            validate:"required"`
	Uri                 string                      `json:"uri"             validate:"required"`
	ProviderType        string                      `json:"provider_type"   validate:"required"`
	ECRCredential       *models.ECRCredential       `json:"ecr_credential,omitempty"`
	DockerhubCredential *models.DockerhubCredential `json:"dockerhub_credential,omitempty"`
	OrganizationId      bson.ObjectID               `json:"organization_id" validate:"required"`
}

type CreateRegistryProvidersResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
