package models

import "time"

type ProviderType string

const (
	ECR       ProviderType = "ECR"
	DockerHub ProviderType = "DockerHub"
)

type RegistryProviders struct {
	Id             string       `json:"id"`
	Name           string       `json:"name"`
	Uri            string       `json:"uri"`
	ProviderType   ProviderType `json:"provider_type"`
	Credential     *Credential  `json:"credential"`
	OrganizationId string       `json:"organization_id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Credential struct {
	ECRCredential       *ECRCredential       `json:"ecr_credential,omitempty" bson:"ecr_credential,omitempty"`
	DockerhubCredential *DockerhubCredential `json:"dockerhub_credential,omitempty" bson:"dockerhub_credential,omitempty"`
}

type ECRCredential struct {
	AccessKey       string `json:"access_key"        bson:"access_key"        validate:"required"`
	SecretAccessKey string `json:"secret_access_key" bson:"secret_access_key" validate:"required"`
	Region          string `json:"aws_region"        bson:"aws_region"        validate:"required"`
}

type DockerhubCredential struct {
	Username            string `json:"username" bson:"username" validate:"required"`
	PersonalAccessToken string `json:"personal_access_token" bson:"personal_access_token" validate:"required"`
}
