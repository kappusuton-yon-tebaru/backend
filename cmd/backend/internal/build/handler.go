package build

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func NewHandler(validator *validator.Validator, service *Service) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) Build(ctx *gin.Context) {
	var req BuildRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "cannot parse json",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	err := h.service.BuildImage(ctx, req)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "image building job submitted",
	})
}
