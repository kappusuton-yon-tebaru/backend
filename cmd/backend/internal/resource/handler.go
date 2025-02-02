package resource

import (
	"net/http"

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

func (h *Handler) CreateResource(ctx *gin.Context) {
	var resourceDTO resource.CreateResourceDTO

	if err := ctx.ShouldBindJSON(&resourceDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateResource(ctx, resourceDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create resource",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":     "resource created successfully",
		"resource_id": id,
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
