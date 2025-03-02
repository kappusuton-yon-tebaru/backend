package utils

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
)

func NewFilterAggregationPipeline(queryParam query.QueryParam, pipelines []map[string]any) []map[string]any {
	filters := []map[string]any{}

	if queryParam.QueryFilter != nil {
		query := queryParam.QueryFilter

		filters = append(filters,
			map[string]any{
				"$match": map[string]any{
					query.Key: map[string]any{
						"$regex":   query.Query,
						"$options": "i",
					},
				},
			},
		)
	}

	if queryParam.SortFilter != nil {
		sort := queryParam.SortFilter

		direction := map[enum.SortOrder]int{
			enum.Asc:  1,
			enum.Desc: -1,
		}

		filters = append(filters,
			map[string]any{
				"$sort": map[string]any{
					sort.SortBy: direction[sort.SortOrder],
				},
			},
		)
	}

	if queryParam.Pagination != nil {
		pagination := queryParam.Pagination
		filters = append(filters,
			[]map[string]any{
				{
					"$facet": map[string]any{
						"metadata": []map[string]any{
							{
								"$count": "total",
							},
							{
								"$addFields": map[string]any{
									"page":  pagination.Page,
									"limit": pagination.Limit,
								},
							},
						},
						"data": []map[string]any{
							{"$skip": (pagination.Page - 1) * pagination.Limit},
							{"$limit": pagination.Limit},
						},
					},
				},
				{
					"$unwind": map[string]any{
						"path":                       "$metadata",
						"preserveNullAndEmptyArrays": true,
					},
				}, {
					"$project": map[string]any{
						"limit": map[string]any{
							"$ifNull": []any{
								"$metadata.limit",
								pagination.Limit,
							},
						},
						"page": map[string]any{
							"$ifNull": []any{
								"$metadata.page",
								pagination.Page,
							},
						},
						"total": map[string]any{
							"$ifNull": []any{
								"$metadata.total",
								0,
							},
						},
						"data": true,
					},
				},
			}...,
		)
	}

	return append(pipelines, filters...)
}
