package usergroup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/usergroup"
)

type Handler struct {
	service *usergroup.Service
}

func NewHandler(service *usergroup.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllUserGroups(ctx *gin.Context) {
	images, err := h.service.GetAllUserGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *Handler) CreateUserGroup(ctx *gin.Context) {
	var usergroupDTO usergroup.CreateUserGroupDTO

	if err := ctx.ShouldBindJSON(&usergroupDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateUserGroup(ctx, usergroupDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":       "user created successfully",
		"user_group_id": id,
	})
}

func (h *Handler) AddUserToUserGroup(ctx *gin.Context) {
	var addUserToUserGroupDTO usergroup.AddUserToUserGroupDTO

	if err := ctx.ShouldBindJSON(&addUserToUserGroupDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	_, err := h.service.AddUserToUserGroup(ctx, addUserToUserGroupDTO, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message":  "add user to group successfully",
		"user_id":  addUserToUserGroupDTO.UserId,
		"group_id": id,
	})
}

func (h *Handler) DeleteUserGroupById(ctx *gin.Context) {
	id := ctx.Param("group_id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteUserGroupById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}

func (h *Handler) DeleteUserFromUserGroupById(ctx *gin.Context) {
	gId := ctx.Param("group_id")
	if len(gId) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	uId := ctx.Param("user_id")
	if len(uId) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteUserFromUserGroupById(ctx, gId, uId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
