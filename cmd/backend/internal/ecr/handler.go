package ecr

import (
	"net/http"
	"strings"

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

func (h *Handler) GetECRImages(ctx *gin.Context) {
	projectId := ctx.Query("project_id")
	serviceName := ctx.Query("service_name")

	if projectId == "" || serviceName == "" {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "missing required query parameters",
		})
		return
	}

	projectRepo, projectRepoErr := h.projectRepoService.GetProjectRepositoryByProjectId(ctx, projectId)
	if projectRepoErr != nil {
		ctx.JSON(http.StatusNotFound, map[string]any{
			"message": "project repository not found",
			"error":   projectRepoErr.Error(),
		})
		return
	}

	if projectRepo.RegistryProvider == nil || len(strings.TrimSpace(projectRepo.RegistryProvider.Uri)) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "registry provider uri must not be empty",
		})
		return
	}

	images, err := h.service.GetECRImages(projectRepo.RegistryProvider.Uri, serviceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}
