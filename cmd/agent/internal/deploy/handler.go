package deploy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sharedDeployEnv "github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
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

// List deployment in project
//
//	@Router			/project/{projectId}/deploy [get]
//	@Summary		List deployment in project
//	@Description	List deployment in project
//	@Tags			Deployment
//	@Param			projectId	path	string	true	"Project Id"
//	@Produce		json
//	@Param			page			query	int		false	"Page"															Default(1)
//	@Param			limit			query	int		false	"Limit"															Default(10)
//	@Param			sort_by			query	string	false	"Sort by"														Enums(age, service_name, status)
//	@Param			sort_order		query	string	false	"Sort order"													Enums(asc, desc)
//	@Param			deployment_env	query	string	false	"Deployment Environment defaults to 'default' if not specified"	Enums(asc, desc)
//	@Param			query			query	string	false	"Query on service_name"
//	@Success		200
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) ListDeployment(ctx *gin.Context) {
	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("service_name", enum.Asc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"age", "service_name", "status"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryFilter := query.NewQueryFilter("service_name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	deployFilter := ListDeploymentQuery{ProjectId: ctx.Param("id"), DeploymentEnv: "default"}
	err = ctx.ShouldBindQuery(&deployFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query",
		})
		return
	}

	deployments, werr := h.service.ListDeployment(ctx, queryParam, deployFilter)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
	}

	ctx.JSON(http.StatusOK, deployments)
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
