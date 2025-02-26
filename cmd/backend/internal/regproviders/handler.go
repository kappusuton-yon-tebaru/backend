package regproviders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
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
	pagination := models.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "pagination should be integer",
		})
		return
	}

	regProviders, err := h.service.GetAllRegistryProviders(ctx, pagination.WithMinimum(1, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": regProviders,
	})
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
