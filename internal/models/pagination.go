package models

type Pagination struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func NewPaginationWithDefault(defaultPage, defaultLimit int) Pagination {
	return Pagination{
		Page:  defaultPage,
		Limit: defaultLimit,
	}
}

func (p Pagination) WithMinimum(minimumPage, minimumLimit int) Pagination {
	return Pagination{
		Page:  max(p.Page, minimumPage),
		Limit: max(p.Limit, minimumLimit),
	}
}

type Paginated[T any] struct {
	Page  int `json:"page"  bson:"page" `
	Limit int `json:"limit" bson:"limit"`
	Total int `json:"total" bson:"total"`
	Data  []T `json:"data"  bson:"data" `
}
