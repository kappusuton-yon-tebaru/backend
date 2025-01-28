package image

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/image"
)

type Handler struct {
	service *image.Service
}

func NewHandler(service *image.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllImages(ctx *gin.Context) {
	images, err := h.service.GetAllImages(ctx)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, images)
}

func (h *Handler) DeleteImage(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(400, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteImage(ctx, id)
	if err != nil {
		ctx.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(299, map[string]any{
		"message": "deleted",
	})
}
