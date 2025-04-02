package deploy

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/logging"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service        *Service
	loggingService *logging.Service
	validator      *validator.Validator
}

func NewHandler(service *Service, loggingService *logging.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		loggingService,
		validator,
	}
}

// Deploy services in project
//
//	@Router			/project/{projectId}/deploy [post]
//	@Summary		Deploy services in project
//	@Description	Deploy services in project
//	@Tags			Deployment
//	@Param			projectId	path	string			true	"Project Id"
//	@Param			request		body	DeployRequest	true	"Optional fields:\n - deployment_env (service will be deployed on __default__ if null)\n - services.\*.port\n - services.\*.secret_name\n - services.\*.health_check"
//	@Produce		json
//	@Success		200 {object} DeployResponse
func (h *Handler) Deploy(ctx *gin.Context) {
	req := DeployRequest{
		ProjectId:     ctx.Param("id"),
		DeploymentEnv: "default",
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

	parentId, werr := h.service.DeployService(ctx, req)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusCreated, DeployResponse{
		parentId,
	})
}

func (h *Handler) GetDeploymentLog(ctx *gin.Context) {
	getLogParam := GetLogQuerParam{DeployEnv: "default"}
	err := ctx.ShouldBindQuery(&getLogParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query params",
		})
		return
	}

	if err := h.validator.Struct(getLogParam); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	filter := query.Filter{
		"project_id":     ctx.Param("id"),
		"deployment_env": getLogParam.DeployEnv,
		"service_name":   utils.ToKebabCase(getLogParam.ServiceName),
	}

	pagination := query.NewCursorPaginationWithDefault(nil, 10, enum.Newer)
	err = ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid pagination",
		})
		return
	}

	if err := h.validator.Struct(pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithCursorPagination(pagination).
		WithFilter(filter)

	logs, werr := h.loggingService.GetLog(ctx, queryParam)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusOK, GetLogResponse{
		Data: logs,
	})
}
