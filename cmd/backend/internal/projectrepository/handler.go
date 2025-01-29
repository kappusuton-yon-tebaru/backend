package projectrepository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
)

type Handler struct {
	service *projectrepository.Service
}

func NewHandler(service *projectrepository.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllProjectRepositories(ctx *gin.Context) {
	projRepos, err := h.service.GetAllProjectRepositories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, projRepos)
}

func (h *Handler) CreateProjectRepository(ctx *gin.Context) {
	var projRepoDTO projectrepository.CreateProjectRepositoryDTO

	if err := ctx.ShouldBindJSON(&projRepoDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateProjectRepository(ctx, projRepoDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create project repository",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":               "project repository created successfully",
		"project_repository_id": id,
	})
}

func (h *Handler) DeleteProjectRepository(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteProjectRepository(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
