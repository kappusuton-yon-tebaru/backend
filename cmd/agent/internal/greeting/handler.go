package greeting

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/greeting"
)

type Handler struct {
	svc greeting.Service
}

func New(svc greeting.Service) *Handler {
	return &Handler{
		svc,
	}
}

func (h *Handler) Greeting(ctx *gin.Context) {
	ctx.String(200, h.svc.GetGreeting())
}
