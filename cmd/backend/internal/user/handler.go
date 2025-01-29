package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{
		service,
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

func (h *Handler) CreateUser(ctx *gin.Context) {
	var userDTO user.CreateUserDTO

	if err := ctx.ShouldBindJSON(&userDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateUser(ctx, userDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "user created successfully",
		"user_id": id,
	})
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
