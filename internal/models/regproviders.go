package models

import "time"

type RegistryProviders struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	Uri            string         `json:"uri"`
	ECRCredential  *ECRCredential `json:"ecr_credential"`
	OrganizationId string         `json:"organization_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type ECRCredential struct {
	AccessKey       string `json:"access_key"        bson:"access_key"        validate:"required"`
	SecretAccessKey string `json:"secret_access_key" bson:"secret_access_key" validate:"required"`
	Region          string `json:"aws_region"        bson:"aws_region"        validate:"required"`
}
