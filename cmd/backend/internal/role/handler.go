package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/role"
)

type Handler struct {
	service *role.Service
}

func NewHandler(service *role.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllRoles(ctx *gin.Context) {
	roles, err := h.service.GetAllRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (h *Handler) CreateRole(ctx *gin.Context) {
	var roleDTO role.CreateRoleDTO

	if err := ctx.ShouldBindJSON(&roleDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateRole(ctx, roleDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create role",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "role created successfully",
		"Role_id": id,
	})
}

func (h *Handler) AddPermissionToRole(ctx *gin.Context) {
	id := ctx.Param("role_id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	var permissionDTO role.CreatePermissionDTO

	if err := ctx.ShouldBindJSON(&permissionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	roleId, err := h.service.AddPermissionToRole(ctx, permissionDTO, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to add permission to role",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message":    "permission added successfully",
		"roleId": roleId,
	})
}

func (h *Handler) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	var roleDTO role.UpdateRoleDTO

	if err := ctx.ShouldBindJSON(&roleDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	roleId, err := h.service.UpdateRole(ctx, roleDTO, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to update role",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message":    "role updated successfully",
		"roleId": roleId,
	})
}

func (h *Handler) DeleteRoleById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteRoleById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
