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

func NewSortQueryWithDefault(defaultKey string, defaultOrder enum.SortOrder) SortFilter {
	return SortFilter{
		SortBy:    defaultKey,
		SortOrder: defaultOrder,
	}
}

type CursorPagination struct {
	Cursor    *string                        `form:"cursor"`
	Direction enum.CursorPaginationDirection `form:"direction" validate:"oneof=newer older"`
	Limit     int                            `form:"limit"`
}

func NewCursorPaginationWithDefault(defaultCursor *string, defaultLimit int, defaultDirection enum.CursorPaginationDirection) CursorPagination {
	return CursorPagination{
		Cursor:    defaultCursor,
		Limit:     defaultLimit,
		Direction: defaultDirection,
	}
}

type Filter map[string]string
