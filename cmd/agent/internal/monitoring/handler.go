package monitoring

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/monitoring"
	"go.uber.org/zap"
)

type Handler struct {
	service *monitoring.Service
	logger  *logger.Logger
}

func NewHandler(service *monitoring.Service, logger *logger.Logger) *Handler {
	return &Handler{
		service,
		logger,
	}
}

func (h *Handler) StreamJobLog(ctx *gin.Context) {
	name := fmt.Sprintf("worker-%s", ctx.Param("id"))

	s, werr := h.service.GetPodLogs(ctx, name)
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": werr.GetMessageOr("internal server error"),
		})
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "could not upgrade to websocket",
		})
	}

	defer func() {
		if err := conn.Close(); err != nil {
			h.logger.Error("error occured while closing websocket connection", zap.Error(err))
		}
	}()

	logCh, unsub := s.Listen()
	defer unsub()

	for {
		log, ok := <-logCh
		if !ok {
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(log))
		if err != nil {
			h.logger.Error("error occured while writing websocket message", zap.Error(err))
			break
		}
	}
}
