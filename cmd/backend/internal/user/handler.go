package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *user.Service
	validator *validator.Validator
}

func NewHandler(service *user.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) GetAllUsers(ctx *gin.Context) {
	images, err := h.service.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *Handler) GetUsersByRoleID(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	if len(roleID) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	users, err := h.service.GetUsersByRoleID(ctx, roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}

func (h *Handler) AddRole(ctx *gin.Context) {
	role_id := ctx.Param("role_id")
	if len(role_id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}
	user_id := ctx.Param("user_id")
	if len(user_id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	roleId, err := h.service.AddRole(ctx, user_id, role_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to add role to user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "role added successfully",
		"roleId":  roleId,
	})
}

func (h *Handler) RemoveRole(ctx *gin.Context) {
	role_id := ctx.Param("role_id")
	if len(role_id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}
	user_id := ctx.Param("user_id")
	if len(user_id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.RemoveRole(ctx, user_id, role_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}

func (h *Handler) UpdateUserRoles(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	if len(roleID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "empty role_id",
		})
		return
	}

	var req struct {
		UserIDs []string `json:"user_ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	if len(req.UserIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "user_ids is required",
		})
		return
	}

	if err := h.service.UpdateUserRoles(ctx, roleID, req.UserIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "users updated successfully",
	})
}
