package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
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

func (h *Handler) GetAllJobParents(ctx *gin.Context) {
	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault(enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid sort query",
		})
		return
	}

	if err := h.validator.Var(sortFilter.SortBy, "oneof=created_at"); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "sort key can only be 'created_at'",
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter)

	jobs, err := h.service.GetAllJobParents(ctx, queryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

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

	sortFilter := query.NewSortQueryWithDefault(enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid sort query",
		})
		return
	}

	if err := h.validator.Var(sortFilter.SortBy, "oneof=created_at job_status"); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "sort key can only be 'created_at' or 'job_status'",
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10))

	jobs, werr := h.service.GetAllJobsByParentId(ctx, id, queryParam)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": werr.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) CreateJob(ctx *gin.Context) {
	var jobDTO job.CreateJobDTO

	if err := ctx.ShouldBindJSON(&jobDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateJob(ctx, jobDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create job",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "job created successfully",
		"id":      id,
	})
}

func (h *Handler) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteJob(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
