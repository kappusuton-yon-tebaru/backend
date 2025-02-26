package utils

import "github.com/kappusuton-yon-tebaru/backend/internal/models"

func NewPaginationAggregationPipeline(pagination models.Pagination, pipelines []map[string]any) []map[string]any {
	return append(pipelines,
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
					"path": "$metadata",
				},
			}, {
				"$project": map[string]any{
					"limit": "$metadata.limit",
					"page":  "$metadata.page",
					"total": "$metadata.total",
					"data":  1,
				},
			},
		}...,
	)
}
