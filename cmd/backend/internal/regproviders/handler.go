package regproviders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
)

type Handler struct {
	service *regproviders.Service
}

func NewHandler(service *regproviders.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllRegProviders(ctx *gin.Context) {
	regProviders, err := h.service.GetAllRegistryProviders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, regProviders)
}

func (h *Handler) CreateRegProvider(ctx *gin.Context) {
	var regprovidersDTO regproviders.CreateRegistryProvidersDTO

	if err := ctx.ShouldBindJSON(&regprovidersDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
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
