package resourcerelationship

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resourcerelationship"
)

type Handler struct {
	service *resourcerelationship.Service
}

func NewHandler(service *resourcerelationship.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllResourceRelationships(ctx *gin.Context) {
	resourceRelas, err := h.service.GetAllResourceRelationships(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resourceRelas)
}

func (h *Handler) GetChildrenResourceRelationshipByParentID(ctx *gin.Context) {
	parentID := ctx.Param("parent_id")

	if parentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parent_id is required"})
		return
	}

	pagination := query.Pagination{
		Limit: 10,
		Page:  1,
	}

	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pagination should be integer"})
		return
	}

	pagination.Page = max(pagination.Page, 1)
	pagination.Limit = max(pagination.Limit, 10)

	projRepo, total, werr := h.service.GetChildrenResourceRelationshipByParentID(ctx, parentID, pagination.Page, pagination.Limit)
	if werr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": werr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": projRepo,
		"pagination": gin.H{
			"page":  pagination.Page,
			"limit": pagination.Limit,
			"total": total,
		},
	})
}

func (h *Handler) CreateResourceRelationship(ctx *gin.Context) {
	var resourceRelaDTO resourcerelationship.CreateResourceRelationshipDTO

	if err := ctx.ShouldBindJSON(&resourceRelaDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateResourceRelationship(ctx, resourceRelaDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create resource relationship",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":                  "resource relationship created successfully",
		"resource_relationship_id": id,
	})
}

func (h *Handler) DeleteResourceRelationship(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteResourceRelationship(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
