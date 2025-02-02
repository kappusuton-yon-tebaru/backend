package ecr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/ecr"
)

type ImageHandler struct {
	service *ecr.Service
}

func NewImageHandler(service *ecr.Service) *ImageHandler {
	return &ImageHandler{
		service,
	}
}

func (h *ImageHandler) GetECRImages(ctx *gin.Context) {
	var req ecr.GetECRImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error": err.Error(),
		})
		return
	}

	images, err := h.service.GetECRImages(req.RepositoryName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *ImageHandler) PushImageToECR(ctx *gin.Context) {
	var req ecr.PushECRImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error": err.Error(),
		})
		return
	}

	id, err := h.service.PushECRImage(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "push successfully",
		"id": id,
	})
}