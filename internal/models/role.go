package models

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type Permission struct {
	Id              string                 `json:"id"`
	PermissionName string                 `json:"permission_name"`
	Action          enum.PermissionActions `json:"action"`
	ResourceId      string                 `json:"resource_id"`
}

type Role struct {
	Id          string       `json:"id"`
	RoleName   string       `json:"role_name"`
	OrgId       string       `json:"orgId"`
	Permissions []Permission `json:"permissions"`
}
