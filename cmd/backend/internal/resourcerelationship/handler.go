package resourcerelationship

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    projRepo, total, err := h.service.GetChildrenResourceRelationshipByParentID(ctx, parentID, page, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": projRepo,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
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
