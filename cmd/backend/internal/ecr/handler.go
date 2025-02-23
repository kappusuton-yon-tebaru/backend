package ecr

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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
	search := ctx.Query("search")
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	if projectId == "" || serviceName == "" {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "missing required query parameters",
		})
		return
	}

	pageNumber := 1;

	if _, err := strconv.Atoi(page); err == nil {
		pageNumber, _ = strconv.Atoi(page)
	}

	limitNumber := 10;

	if _, err := strconv.Atoi(limit); err == nil {
		limitNumber, _ = strconv.Atoi(limit)
	}

	projectRepo, projectRepoErr := h.projectRepoService.GetProjectRepositoryByProjectId(ctx, projectId)
	if projectRepoErr != nil {
		ctx.JSON(http.StatusNotFound, map[string]any{
			"message": "project repository not found",
			"error":   projectRepoErr.Error(),
		})
		return
	}

	pagination := models.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "pagination should be integer",
		})
		return
	}

	images, err := h.service.GetECRImages(projectRepo.RegistryProvider.Uri, serviceName, pagination.WithMinimum(1, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}
