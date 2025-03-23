package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateRegistryProvidersRequest struct {
	Name           string                `json:"name"            validate:"required"`
	Uri            string                `json:"uri"             validate:"required"`
	ECRCredential  *models.ECRCredential `json:"ecr_credential"  validate:"required"`
	OrganizationId bson.ObjectID         `json:"organization_id" validate:"required"`
}

type CreateRegistryProvidersResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type UpdateRegistryProvidersRequest struct {
	Name           string                `json:"name"            validate:"omitempty"`
	Uri            string                `json:"uri"             validate:"omitempty"`
	ECRCredential  *models.ECRCredential `json:"ecr_credential"  validate:"omitempty"`
	OrganizationId bson.ObjectID         `json:"organization_id" validate:"omitempty"`
}

type UpdateRegistryProvidersResponse struct {
	Message string `json:"message"`
}
