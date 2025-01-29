package permission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/permission"
)

type Handler struct {
	service *permission.Service
}

func NewHandler(service *permission.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllPermissions(ctx *gin.Context) {
	images, err := h.service.GetAllPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *Handler) DeletePermissionById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeletePermissionById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}