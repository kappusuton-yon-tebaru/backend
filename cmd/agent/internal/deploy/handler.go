package deploy

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sharedDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *sharedDeployEnv.Service
	validator *validator.Validator
}

func NewHandler(service *sharedDeployEnv.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) ListDeploymentEnv(ctx *gin.Context) {
	projectId := ctx.Param("id")

	namespaces, werr := h.service.ListDeploymentEnv(ctx, projectId)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(200, map[string]any{
		"data": namespaces,
	})
}

func (h *Handler) CreateDeploymentEnv(ctx *gin.Context) {
	req := ModifyDeploymentEnvRequest{
		ProjectId: ctx.Param("id"),
	}

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

	werr := h.service.CreateDeploymentEnv(ctx, sharedDeployEnv.ModifyDeploymentEnvDTO(req))
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.Status(201)
}

func (h *Handler) DeleteDeploymentEnv(ctx *gin.Context) {
	req := ModifyDeploymentEnvRequest{
		ProjectId: ctx.Param("id"),
	}

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

	werr := h.service.DeleteDeploymentEnv(ctx, sharedDeployEnv.ModifyDeploymentEnvDTO(req))
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteDeployment(ctx *gin.Context) {
	var req DeleteDeploymentRequest
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

	werr := h.service.DeleteDeployment(ctx, sharedDeployEnv.DeleteDeploymentDTO(req))
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}
}
