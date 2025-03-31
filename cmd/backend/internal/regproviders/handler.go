package regproviders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	_ "github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *regproviders.Service
	validator *validator.Validator
}

func NewHandler(service *regproviders.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

// List all registry providers
//
//	@Router			/regproviders [get]
//	@Summary		List all registry providers
//	@Description	List all registry providers
//	@Tags			regproviders
//	@Produce		json
//	@Param			page		query		int		false	"Page"			Default(1)
//	@Param			limit		query		int		false	"Limit"			Default(10)
//	@Param			sort_by		query		string	false	"Sort by"		Enums(created_at, name)
//	@Param			sort_order	query		string	false	"Sort order"	Enums(asc, desc)
//	@Param			query		query		string	false	"Query on name"
//	@Success		200			{object}	models.Paginated[models.RegistryProviders]
//	@Failure		400			{object}	httputils.ErrResponse
//	@Failure		500			{object}	httputils.ErrResponse
func (h *Handler) GetAllRegProviders(ctx *gin.Context) {
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

	queryFilter := query.NewQueryFilter("name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid query filter",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	regProviders, err := h.service.GetAllRegistryProviders(ctx, queryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httputils.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, regProviders)
}

// List all registry providers without project binding to it
//
//	@Router			/regproviders/unbind [get]
//	@Summary		List all registry providers without project binding to it
//	@Description	List all registry providers without project binding to it
//	@Tags			regproviders
//	@Produce		json
//	@Success		200	{object}	models.RegistryProviders[]
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) GetAllRegProvidersWithoutProject(ctx *gin.Context) {
	regProviders, err := h.service.GetAllRegistryProvidersWithoutProject(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httputils.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"data": regProviders,
	})
}

// List registry provider by ID
//
//	@Router			/regproviders/{id} [get]
//	@Summary		List registry provider by ID
//	@Description	List registry provider by ID
//	@Tags			regproviders
//	@Produce		json
//	@Param			id	path		string	true	"Registry Provider ID"
//	@Success		200	{object}	models.RegistryProviders
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) GetRegProviderById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "empty id",
		})
		return
	}

	regProvider, err := h.service.GetRegistryProviderById(ctx, id)

	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), httputils.ErrResponse{
			Message: err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"data": regProvider,
	})
}

// Create registry provider
//
//	@Router			/regproviders [post]
//	@Summary		Create registry provider
//	@Description	Create registry provider
//	@Tags			regproviders
//	@Produce		json
//	@Param			body	body		CreateRegistryProvidersRequest	true	"Create Registry Provider Request"
//	@Success		201		{object}	CreateRegistryProvidersResponse
//	@Failure		400		{object}	httputils.ErrResponse
//	@Failure		500		{object}	httputils.ErrResponse
func (h *Handler) CreateRegProvider(ctx *gin.Context) {
	var req CreateRegistryProvidersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid body",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	dto := regproviders.CreateRegistryProvidersDTO{
		Name:           req.Name,
		Uri:            req.Uri,
		ECRCredential:  req.ECRCredential,
		OrganizationId: req.OrganizationId,
	}

	id, werr := h.service.CreateRegistryProviders(ctx, dto)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusCreated, CreateRegistryProvidersResponse{
		Message: "registry provider created successfully",
		Id:      id,
	})
}

// Delete registry provider
//
//	@Router			/regproviders/{id} [delete]
//	@Summary		Delete registry provider
//	@Description	Delete registry provider
//	@Tags			regproviders
//	@Produce		json
//	@Param			id	path		string	true	"Registry Provider ID"
//	@Success		200	{object}	map[string]any
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) DeleteRegProvider(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteRegistryProviders(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), httputils.ErrResponse{
			Message: err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
