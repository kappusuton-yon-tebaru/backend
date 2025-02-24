package resource

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
)

type Handler struct {
	service *resource.Service
}

func NewHandler(service *resource.Service) *Handler {
	return &Handler{
		service,
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
	id := ctx.Param("parent_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parent_id is required"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	resources, total, err := h.service.GetChildrenResourcesByParentID(ctx, id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  resources,
		"page":  page,
		"limit": limit,
		"total": total,
	})
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
