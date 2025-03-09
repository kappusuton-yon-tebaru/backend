package regproviders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
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

func (h *Handler) GetAllRegProviders(ctx *gin.Context) {
	pagination := query.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "pagination should be integer",
		})
		return
	}

	sortFilter := query.NewSortQueryWithDefault("created_at", enum.Desc)
	err = ctx.ShouldBindQuery(&sortFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid sort query",
		})
		return
	}

	availableSortKey := []string{"created_at", "name"}
	if err := h.validator.Var(sortFilter.SortBy, fmt.Sprintf("omitempty,oneof=%s", strings.Join(availableSortKey, " "))); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": fmt.Sprintf("sort key can only be one of the field: %s", utils.ArrayWithComma(availableSortKey, "or")),
		})
		return
	}

	if err := h.validator.Struct(sortFilter); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "sort order can only be 'asc' or 'desc'",
		})
		return
	}

	queryFilter := query.NewQueryFilter("name")
	err = ctx.ShouldBindQuery(&queryFilter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid query filter",
		})
		return
	}

	queryParam := query.NewQueryParam().
		WithPagination(pagination.WithMinimum(1, 10)).
		WithSortQuery(sortFilter).
		WithQueryFilter(queryFilter)

	regProviders, err := h.service.GetAllRegistryProviders(ctx, queryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, regProviders)
}

func (h *Handler) GetAllRegProvidersWithoutProject(ctx *gin.Context) {
	regProviders, err := h.service.GetAllRegistryProvidersWithoutProject(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": regProviders,
	})
}

func (h *Handler) GetRegProviderById(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty id",
		})
		return
	}

	regProvider, err := h.service.GetRegistryProviderById(ctx, id)

	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]interface{}{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": regProvider,
	})
}

func (h *Handler) CreateRegProvider(ctx *gin.Context) {
	var req CreateRegistryProvidersRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	credential, werr := regproviders.ParseCredential(req.ProviderType, req.Credential)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusBadRequest), map[string]interface{}{
			"message": werr.GetMessageOr("bad request"),
		})
	}

	if err := h.validator.Struct(credential); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"messages": h.validator.Translate(err),
		})
		return
	}

	regprovidersDTO := regproviders.CreateRegistryProvidersDTO{
		Name:           req.Name,
		ProviderType:   req.ProviderType,
		Uri:            req.Uri,
		Credential:     credential,
		OrganizationId: req.OrganizationId,
	}

	id, err := h.service.CreateRegistryProviders(ctx, regprovidersDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to create registry provider",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "registry provider created successfully",
		"id":      id,
	})
}

func (h *Handler) DeleteRegProvider(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteRegistryProviders(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]interface{}{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "deleted",
	})
}
