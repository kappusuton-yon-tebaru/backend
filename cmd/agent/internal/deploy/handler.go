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

// List deployment environments in project
//
//	@Router			/project/{projectId}/deployenv [get]
//	@Summary		List deployment environments in project
//	@Description	List deployment environments in project
//	@Tags			Deployment Environment
//	@Param			projectId	path	string	true	"Project Id"
//	@Produce		json
//	@Success		200	{object}	ListDeploymentEnvResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) ListDeploymentEnv(ctx *gin.Context) {
	projectId := ctx.Param("id")

	namespaces, werr := h.service.ListDeploymentEnv(ctx, projectId)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(200, ListDeploymentEnvResponse{
		Data: namespaces,
	})
}

// Create deployment environment in project
//
//	@Router			/project/{projectId}/deployenv [post]
//	@Summary		Create deployment environment in project
//	@Description	Create deployment environment in project
//	@Tags			Deployment Environment
//	@Param			projectId	path	string						true	"Project Id"
//	@Param			request		body	ModifyDeploymentEnvRequest	true	"create deployment environment request"
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
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

// Delete deployment environment in project
//
//	@Router			/project/{projectId}/deployenv [delete]
//	@Summary		Delete deployment environment in project
//	@Description	Delete deployment environment in project
//	@Tags			Deployment Environment
//	@Param			projectId	path	string						true	"Project Id"
//	@Param			request		body	ModifyDeploymentEnvRequest	true	"delete deployment environment request"
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
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
