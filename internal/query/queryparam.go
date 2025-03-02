package query

type QueryParam struct {
	Pagination  *Pagination
	SortFilter  *SortFilter
	QueryFilter *QueryFilter
}

func NewQueryParam() QueryParam {
	return QueryParam{
		Pagination:  nil,
		SortFilter:  nil,
		QueryFilter: nil,
	}
}

func (q QueryParam) WithPagination(p Pagination) QueryParam {
	q.Pagination = &p
	return q
}

func (q QueryParam) WithSortQuery(s SortFilter) QueryParam {
	q.SortFilter = &s
	return q
}

func (q QueryParam) WithQueryFilter(s QueryFilter) QueryParam {
	q.QueryFilter = &s
	return q
}
