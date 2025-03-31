package ecr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	_ "github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service            *Service
	projectRepoService *projectrepository.Service
	validator          *validator.Validator
}

func NewHandler(service *Service, projectRepoService *projectrepository.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		projectRepoService,
		validator,
	}
}

// List all ecr images by project id and service name
//
//	@Router			/ecr/images [get]
//	@Summary		List all ecr images by project id and service name
//	@Description	List all ecr images by project id and service name
//	@Tags			ECR
//	@Produce		json
//	@Param			project_id		query		string	true	"Project ID"
//	@Param			service_name	query		string	true	"Service Name"
//	@Param			page			query		int		false	"Page"			Default(1)
//	@Param			limit			query		int		false	"Limit"			Default(10)
//	@Param			sort_by			query		string	false	"Sort by"		Enums(created_at, name)
//	@Param			sort_order		query		string	false	"Sort order"	Enums(asc, desc)
//	@Param			query			query		string	false	"Query on image_tag"
//	@Success		200				{object}	models.Paginated[ECRImageResponse]
//	@Failure		400				{object}	httputils.ErrResponse
//	@Failure		404				{object}	httputils.ErrResponse
//	@Failure		500				{object}	httputils.ErrResponse
func (h *Handler) GetECRImages(ctx *gin.Context) {
	projectId := ctx.Query("project_id")
	serviceName := ctx.Query("service_name")

	if projectId == "" || serviceName == "" {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "missing required query parameters",
		})
		return
	}

	projectRepo, projectRepoErr := h.projectRepoService.GetProjectRepositoryByProjectId(ctx, projectId)
	if projectRepoErr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(projectRepoErr))
		return
	}

	if projectRepo.RegistryProvider.ECRCredential == nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid ecr credential",
		})
	}

	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("created_at", enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"created_at", "name"}
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

	queryFilter := query.NewQueryFilter("image_tag")
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

	images, err := h.service.GetECRImages(ctx, *projectRepo.RegistryProvider, serviceName, queryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httputils.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, images)
}
