package query

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

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

type QueryFilter struct {
	Query string `form:"query"`
	Key   string
}

func NewQueryFilter(key string) QueryFilter {
	return QueryFilter{
		Key: key,
	}
}

type SortFilter struct {
	SortBy    string         `form:"sort_by"`
	SortOrder enum.SortOrder `form:"sort_order" validate:"oneof=asc desc"`
}

func NewSortQueryWithDefault(defaultOrder enum.SortOrder) SortFilter {
	return SortFilter{
		SortBy:    "",
		SortOrder: defaultOrder,
	}
}
