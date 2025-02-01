package models

type RoleUserGroup struct {
	Id       	string `json:"id"`
	RoleId 		string `json:"role_id"`
	UserGroupId 		string `json:"user_group_id"`
}