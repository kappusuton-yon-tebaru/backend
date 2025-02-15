package ecr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetECRImages(ctx *gin.Context) {
	var req GetECRImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error": err.Error(),
		})
		return
	}

	images, err := h.service.GetECRImages(req.RepositoryURI, req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}