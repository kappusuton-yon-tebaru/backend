package greeting

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Greeting(ctx *gin.Context) {
	ctx.String(200, "Hello from backend!")
}
