package models

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type Permission struct {
	Id              string                 `json:"id"`
	Permission_name string                 `json:"permission_name"`
	Action          enum.PermissionActions `json:"action"`
	ResourceId      string                 `json:"resource_id"`
	ResourceType    enum.ResourceType      `json:"resource_type"`
}
