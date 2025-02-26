package models

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
)

type Resource struct {
	Id           string            `json:"id"`
	ResourceName string            `json:"resource_name"`
	ResourceType enum.ResourceType `json:"resource_type"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
