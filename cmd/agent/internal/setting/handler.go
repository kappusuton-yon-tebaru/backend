package setting

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	validator *validator.Validator
	service   *Service
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{
		validator,
		service,
	}
}

// Set worker pool setting
//
//	@Router			/setting/workerpool [post]
//	@Summary		Set worker pool setting
//	@Description	Set worker pool setting
//	@Tags			Setting
//	@Param			request	body	SetWorkerPoolRequest	true	"worker pool setting request"
//	@Produce		json
//	@Success		200	{object}	WorkerPoolResponse
//	@Failure		400	{object}	httputils.ErrResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) SetWorkerPoolSetting(ctx *gin.Context) {
	var req SetWorkerPoolRequest
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

	werr := h.service.SetMaxWorker(ctx, req.PoolSize)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(200, WorkerPoolResponse(req))
}

// Get worker pool setting
//
//	@Router			/setting/workerpool [get]
//	@Summary		Get worker pool setting
//	@Description	Get worker pool setting
//	@Tags			Setting
//	@Produce		json
//	@Success		200	{object}	WorkerPoolResponse
//	@Failure		500	{object}	httputils.ErrResponse
func (h *Handler) GetWorkerPoolSetting(ctx *gin.Context) {
	maxWorker, werr := h.service.GetMaxWorker(ctx)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusOK, WorkerPoolResponse{
		PoolSize: maxWorker,
	})
}
