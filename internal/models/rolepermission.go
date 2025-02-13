package models

type RolePermission struct {
	Id           string `json:"id"`
	RoleId       string `json:"role_id"`
	PermissionId string `json:"permission_id"`
}
