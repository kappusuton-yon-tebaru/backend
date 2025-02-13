package projectenv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectenv"
)

type Handler struct {
	service *projectenv.Service
}

func NewHandler(service *projectenv.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllProjectEnvs(ctx *gin.Context) {
	projectEnvs, err := h.service.GetAllProjectEnvs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, projectEnvs)
}

func (h *Handler) CreateProjectEnv(ctx *gin.Context) {
	var projectEnvDTO projectenv.CreateProjectEnvDTO

	if err := ctx.ShouldBindJSON(&projectEnvDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateProjectEnv(ctx, projectEnvDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to create project env",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "project env created successfully",
		"id":      id,
	})
}

func (h *Handler) DeleteProjectEnv(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteProjectEnv(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]interface{}{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "deleted",
	})
}
