package monitoring

import (
	"fmt"
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

	done := make(chan struct{})

	// TODO: properly handle websocket connection

	go func() {
		for {
			select {
			case _, ok := <-done:
				if !ok {
					return
				}
			default:
				t, bs, err := upConn.ReadMessage()
				if err != nil {
					h.logger.Error("error occured while reading websocket message from upstream", zap.Error(err))
					close(done)
					return
				}

				err = conn.WriteMessage(t, bs)
				if err != nil {
					h.logger.Error("error occured while writing websocket message to downstream", zap.Error(err))
					close(done)
					return
				}
			}
		}
	}()

	func() {
		for {
			select {
			case _, ok := <-done:
				if !ok {
					return
				}
			default:
				t, bs, err := conn.ReadMessage()
				if err != nil {
					h.logger.Error("error occured while reading websocket message from downstream", zap.Error(err))
					close(done)
					return
				}

				err = upConn.WriteMessage(t, bs)
				if err != nil {
					h.logger.Error("error occured while writing websocket message to upstream", zap.Error(err))
					close(done)
					return
				}
			}
		}
	}()

	fmt.Println("ws closed")
}
