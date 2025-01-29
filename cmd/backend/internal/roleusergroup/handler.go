package roleusergroup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/roleusergroup"
)

type Handler struct {
	service *roleusergroup.Service
}

func NewHandler(service *roleusergroup.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllRoleUserGroups(ctx *gin.Context) {
	images, err := h.service.GetAllRoleUserGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (h *Handler) CreateRoleUserGroup(ctx *gin.Context) {
	var roleUserGroupDTO roleusergroup.CreateRoleUserGroupDTO

	if err := ctx.ShouldBindJSON(&roleUserGroupDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateRoleUserGroup(ctx,roleUserGroupDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create roleUserGroup",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "roleUserGroup created successfully",
		"Role_UserGroup_id": id,
	})
}

func (h *Handler) DeleteRoleUserGroupById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteRoleUserGroupById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}