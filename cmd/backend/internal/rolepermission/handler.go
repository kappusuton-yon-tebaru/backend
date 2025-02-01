package rolepermission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/rolepermission"
)

type Handler struct {
	service *rolepermission.Service
}

func NewHandler(service *rolepermission.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllRolePermissions(ctx *gin.Context) {
	images, err := h.service.GetAllRolePermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *Handler) CreateRolePermission(ctx *gin.Context) {
	var rolepermissionDTO rolepermission.CreateRolePermissionDTO

	if err := ctx.ShouldBindJSON(&rolepermissionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateRolePermission(ctx, rolepermissionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create rolepermission",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":            "rolepermission created successfully",
		"Role_permission_id": id,
	})
}

func (h *Handler) DeleteRolePermissionById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteRolePermissionById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
