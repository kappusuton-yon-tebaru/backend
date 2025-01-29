package resourcerelationship

import (
	"net/http"

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
