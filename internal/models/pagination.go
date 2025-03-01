package models

type Paginated[T any] struct {
	Page  int `json:"page"  bson:"page" `
	Limit int `json:"limit" bson:"limit"`
	Total int `json:"total" bson:"total"`
	Data  []T `json:"data"  bson:"data" `
}
