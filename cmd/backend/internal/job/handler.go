package job

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *job.Service
	validator *validator.Validator
}

func NewHandler(service *job.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

// List all job parents
//
//	@Router			/jobs [get]
//	@Summary		List all job parents
//	@Description	List all job parents
//	@Tags			Jobs
//	@Produce		json
//	@Param			page		query		int		false	"Page"			Default(1)
//	@Param			limit		query		int		false	"Limit"			Default(10)
//	@Param			sort_by		query		string	false	"Sort by"		Enums(created_at, project.name)
//	@Param			sort_order	query		string	false	"Sort order"	Enums(asc, desc)
//	@Param			query		query		string	false	"Query on project.resource_name"
//	@Success		200			{object}	job.PaginatedJobs
//	@Failure		400			{object}	httputils.ErrResponse
//	@Failure		500			{object}	httputils.ErrResponse
func (h *Handler) GetAllJobParents(ctx *gin.Context) {
	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("created_at", enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"created_at", "project.name"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryFilter := query.NewQueryFilter("project.resource_name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	jobs, werr := h.service.GetAllJobParents(ctx, queryParam)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

// List jobs by job parent id
//
//	@Router				/jobs/{jobParentId} [get]
//	@Summary			List jobs by job parent id
//	@DescriptionList	jobs by job parent id
//	@Tags				Jobs
//	@Produce			json
//	@Param				jobParentId	path		string	true	"Job Parent Id"
//	@Param				page		query		int		false	"Page"			Default(1)
//	@Param				limit		query		int		false	"Limit"			Default(10)
//	@Param				sort_by		query		string	false	"Sort by"		Enums(created_at, job_status, service_name, project.name)
//	@Param				sort_order	query		string	false	"Sort order"	Enums(asc, desc)
//	@Param				query		query		string	false	"Query on service_name"
//	@Success			200			{object}	job.PaginatedJobs
//	@Failure			400			{object}	httputils.ErrResponse
//	@Failure			500			{object}	httputils.ErrResponse
func (h *Handler) GetAllJobsByParentId(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pagination should be integer"})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("created_at", enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"created_at", "job_status", "service_name", "project.name"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryFilter := query.NewQueryFilter("service_name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid query",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	jobs, werr := h.service.GetAllJobsByParentId(ctx, id, queryParam)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": werr.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}
