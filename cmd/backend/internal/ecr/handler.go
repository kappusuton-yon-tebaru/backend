package ecr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
)

type Handler struct {
	service *Service
	projectRepoService *projectrepository.Service
}

func NewHandler(service *Service, projectRepoService *projectrepository.Service) *Handler {
	return &Handler{
		service,
		projectRepoService,
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

	projectRepo, projectRepoErr := h.projectRepoService.GetProjectRepositoryByProjectId(ctx, req.ProjectId);
	if projectRepoErr != nil {
		ctx.JSON(http.StatusNotFound, map[string]any{
			"message": "project repository not found",
			"error": projectRepoErr.Error(),
		})
		return
	}

	images, err := h.service.GetECRImages(projectRepo.RegistryProvider.Uri, req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}