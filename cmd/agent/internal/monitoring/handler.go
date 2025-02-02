package monitoring

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/monitoring"
)

type Handler struct {
	service *monitoring.Service
}

func NewHandler(service *monitoring.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) StreamJobLog(ctx *gin.Context) {
	name := fmt.Sprintf("worker-%s", ctx.Param("id"))

	s, err := h.service.GetPodLogs(ctx, name)
	if err != nil {
		// ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]any{
		// 	"message": err.GetMessageOr("internal server error"),
		// })
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
		return
	}

	ctx.Stream(func(w io.Writer) bool {
		logCh, unsub := s.Listen()
		defer unsub()

		for {
			select {
			case <-ctx.Writer.CloseNotify():
				return false
			case log, ok := <-logCh:
				if !ok {
					ctx.SSEvent("close", "")
					return false
				}

				ctx.SSEvent("message", log)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	})
}
