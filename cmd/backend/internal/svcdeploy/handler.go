package svcdeploy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/svcdeploy"
)

type Handler struct {
	service *svcdeploy.Service
}

func NewHandler(service *svcdeploy.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllServiceDeployments(ctx *gin.Context) {
	svcDeploys, err := h.service.GetAllServiceDeployments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, svcDeploys)
}

func (h *Handler) DeleteServiceDeployment(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteServiceDeployment(ctx, id)
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
