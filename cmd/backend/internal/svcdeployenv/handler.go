package svcdeployenv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/svcdeployenv"
)

type Handler struct {
	service *svcdeployenv.Service
}

func NewHandler(service *svcdeployenv.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllServiceDeploymentEnvs(ctx *gin.Context) {
	svcDeploys, err := h.service.GetAllServiceDeploymentEnvs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, svcDeploys)
}

func (h *Handler) DeleteServiceDeploymentEnv(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteServiceDeploymentEnv(ctx, id)
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
