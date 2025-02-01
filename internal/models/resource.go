package models

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type Resource struct {
	Id           string            `json:"id"`
	ResourceName string            `json:"resource_name"`
	ResourceType enum.ResourceType `json:"resource_type"`
}
