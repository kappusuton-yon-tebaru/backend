package logging

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/logging"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	service   *logging.Service
	validator *validator.Validator
}

func NewHandler(service *logging.Service, validator *validator.Validator) *Handler {
	return &Handler{
		service,
		validator,
	}
}

func (h *Handler) GetJobLog(ctx *gin.Context) {
	filter := query.Filter{"job_id": ctx.Param("id")}

	pagination := query.NewCursorPaginationWithDefault(nil, 10, enum.Newer)
	err := ctx.ShouldBindQuery(&pagination)
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

	logs, werr := h.service.GetLog(ctx, queryParam)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusOK, GetLogResponse{
		Data: logs,
	})
}
