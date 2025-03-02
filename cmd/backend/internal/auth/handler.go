package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/auth"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *auth.Service
	validator *validator.Validator
}

func NewHandler(service *auth.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterReq
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

	credential := user.UserCredentialDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	werr := h.service.Register(ctx, credential)
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
