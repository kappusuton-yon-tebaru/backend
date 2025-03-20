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
	sharedService *sharedDeployEnv.Service
	service       *Service
	validator     *validator.Validator
}

func NewHandler(sharedService *sharedDeployEnv.Service, service *Service, validator *validator.Validator) *Handler {
	return &Handler{
		sharedService,
		service,
		validator,
	}
}

// Delete deployment in project
//
//	@Router			/project/{projectId}/deploy [delete]
//	@Summary		Delete deployment in project
//	@Description	Delete deployment in project
//	@Tags			Deployment
//	@Param			projectId	path	string					true	"Project Id"
//	@Param			request		body	DeleteDeploymentRequest	true	"Optional fields:\n - deployment_env"
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) DeleteDeployment(ctx *gin.Context) {
	req := DeleteDeploymentRequest{DeploymentEnv: "default", ProjectId: ctx.Param("id")}
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

	werr := h.service.DeleteDeployment(ctx, req)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.Status(http.StatusOK)
}
