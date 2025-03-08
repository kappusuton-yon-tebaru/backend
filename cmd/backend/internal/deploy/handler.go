package deploy

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	validator *validator.Validator
}

func NewHandler(validator *validator.Validator) *Handler {
	return &Handler{
		validator,
	}
}

// Deploy services in project
//
//	@Router			/deploy [post]
//	@Summary		Deploy services in project
//	@Description	Deploy services in project
//	@Tags			Build
//	@Param			request	body	DeployRequest	true "deploy request"
//	@Produce		json
//	@Success		200
func (h *Handler) Build(ctx *gin.Context) {
	var req DeployRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "cannot parse json",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	ctx.Status(http.StatusOK)
}
