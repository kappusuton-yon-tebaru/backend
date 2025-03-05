package build

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
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

// Build services in project
//
//	@Router			/build [post]
//	@Summary		Build services in project
//	@Description	Build services in project
//	@Tags			Build
//	@Produce		json
//	@Success		201	{object}	BuildResponse
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) Build(ctx *gin.Context) {
	var req BuildRequest
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

	parentId, werr := h.service.BuildImage(ctx, req)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusCreated, BuildResponse{
		parentId,
	})
}
