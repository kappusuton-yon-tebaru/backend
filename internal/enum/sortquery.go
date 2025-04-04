package enum

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type CursorPaginationDirection string

const (
	Newer CursorPaginationDirection = "newer"
	Older CursorPaginationDirection = "older"
)
