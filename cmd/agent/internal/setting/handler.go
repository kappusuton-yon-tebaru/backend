package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	validator *validator.Validator
	service   *Service
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{
		validator,
		service,
	}
}

func (h *Handler) SetMaxWorker(ctx *gin.Context) {
	var req SetMaxWorkerDTO
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

	werr := h.service.SetMaxWorker(ctx, req.MaxWorker)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": werr.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(200, map[string]any{
		"max_worker": req.MaxWorker,
	})
}
