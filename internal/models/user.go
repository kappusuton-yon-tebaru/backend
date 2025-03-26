package models

// example model
type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	RoleIds 	 []string `json:"role_ids"`
}
