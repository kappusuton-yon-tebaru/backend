package models

import "time"

type RegistryProviders struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	ProviderType   string    `json:"provider_type"` // enum
	Uri            string    `json:"uri"`
	JsonCredential string    `json:"json_credential"`
	OrganizationId string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
