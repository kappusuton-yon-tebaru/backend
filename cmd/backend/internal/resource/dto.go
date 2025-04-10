package resource

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type CreateResourceRequest struct {
	ParentId     string            `validate:"required"`
	ResourceName string            `validate:"required,kebabnum"`
	ResourceType enum.ResourceType `validate:"required"`
}

type UpdateResourceRequest struct {
	ResourceName string            `validate:"required,kebabnum"`
}
