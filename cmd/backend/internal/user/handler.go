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

func (h *Handler) Register(ctx *gin.Context) {
	var req user.RegisterDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid body",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"errors": h.validator.Translate(err),
		})
		return
	}

	werr := h.service.Register(ctx, req)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"error": werr.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "user created",
	})
}

func (h *Handler) GetAllUsers(ctx *gin.Context) {
	images, err := h.service.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, images)
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
