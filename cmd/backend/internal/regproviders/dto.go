package regproviders

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateRegistryProvidersRequest struct {
	Name           string                    `json:"name"`
	ProviderType   enum.RegistryProviderType `json:"provider_type"`
	Uri            string                    `json:"uri"`
	Credential     map[string]any            `json:"credential"`
	OrganizationId bson.ObjectID             `json:"organization_id"`
}
