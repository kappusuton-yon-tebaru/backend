package ecr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service            *Service
	projectRepoService *projectrepository.Service
	validator          *validator.Validator
}

func NewHandler(service *Service, projectRepoService *projectrepository.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		projectRepoService,
		validator,
	}
}

func (h *Handler) GetECRImages(ctx *gin.Context) {
	projectId := ctx.Query("project_id")
	serviceName := ctx.Query("service_name")

	if projectId == "" || serviceName == "" {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "missing required query parameters",
		})
		return
	}

	projectRepo, projectRepoErr := h.projectRepoService.GetProjectRepositoryByProjectId(ctx, projectId)
	if projectRepoErr != nil {
		ctx.JSON(http.StatusNotFound, map[string]any{
			"message": "project repository not found",
			"error":   projectRepoErr.Error(),
		})
		return
	}

	if projectRepo.RegistryProvider.ECRCredential == nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid ecr credential",
		})
	}

	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("created_at", enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"created_at", "name"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryFilter := query.NewQueryFilter("image_tag")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid query",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	images, err := h.service.GetECRImages(ctx, *projectRepo.RegistryProvider, serviceName, queryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}
