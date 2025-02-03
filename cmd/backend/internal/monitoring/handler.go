package monitoring

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"go.uber.org/zap"
)

type Handler struct {
	logger *logger.Logger
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		logger,
	}
}

func (h *Handler) StreamJobLog(ctx *gin.Context) {
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

	upConn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws/job/e7a26802-eb3a-4001-ab26-5d4adb05c99c/log", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "could not connect to upstream websocket",
		})
	}

	defer func() {
		if err := upConn.Close(); err != nil {
			h.logger.Error("error occured while closing websocket connection", zap.Error(err))
		}
	}()

	for {
		t, bs, err := upConn.ReadMessage()
		if err != nil {
			h.logger.Error("error occured while reading websocket message from upstream", zap.Error(err))
			break
		}

		err = conn.WriteMessage(t, bs)
		if err != nil {
			h.logger.Error("error occured while writing websocket message", zap.Error(err))
			break
		}
	}
}

