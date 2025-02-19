package dockerhub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
)

type Handler struct {
	service            *Service
	projectRepoService *projectrepository.Service
}

func NewHandler(service *Service, projectRepoService *projectrepository.Service) *Handler {
	return &Handler{
		service,
		projectRepoService,
	}
}

func (h *Handler) GetDockerHubImages(ctx *gin.Context) {
	var req GetDockerHubImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	images, err := h.service.GetDockerHubImages(req.Namespace, req.RepositoryName, req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}
