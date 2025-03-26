package resource

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *resource.Service
	validator *validator.Validator
}

func NewHandler(service *resource.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) GetAllResources(ctx *gin.Context) {
	resources, err := h.service.GetAllResources(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resources)
}

func (h *Handler) GetResourceByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	resource, err := h.service.GetResourceByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resource)
}

func (h *Handler) GetChildrenResourcesByParentID(ctx *gin.Context) {
	parentId := ctx.Param("parent_id")
	if parentId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parent_id is required"})
	}

	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "pagination should be integer",
		})
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

	availableSortKey := []string{"created_at", "resource_name", "updated_at"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	queryFilter := query.NewQueryFilter("resource_name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query",
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
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	resources, err := h.service.GetChildrenResourcesByParentID(ctx, queryParam, parentId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resources)
}

func (h *Handler) CreateResource(ctx *gin.Context) {
	id := ctx.DefaultQuery("parent_id", "")

	var resourceDTO resource.CreateResourceDTO

	if err := ctx.ShouldBindJSON(&resourceDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	resourceId, err := h.service.CreateResource(ctx, resourceDTO, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create resource",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":    "resource created successfully",
		"resourceId": resourceId,
	})
}

func (h *Handler) UpdateResource(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	var resourceDTO resource.UpdateResourceDTO

	if err := ctx.ShouldBindJSON(&resourceDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	resourceId, err := h.service.UpdateResource(ctx, resourceDTO, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to update resource",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message":    "resource updated successfully",
		"resourceId": resourceId,
	})
}

func (h *Handler) DeleteResource(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteResource(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}

func (h *Handler) CascadeDeleteResource(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.CascadeDeleteResource(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
