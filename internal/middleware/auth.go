package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/auth"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/role"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
)

type Middleware struct {
	config      *config.Config
	authService *auth.Service
	logger      *logger.Logger
	roleService *role.Service
}

func NewMiddleware(config *config.Config, authService *auth.Service, logger *logger.Logger, roleService *role.Service) *Middleware {
	return &Middleware{
		config,
		authService,
		logger,
		roleService,
	}
}

func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("session_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionId, err := utils.ParseToken(token, m.config.JwtSecret)
		if errors.Is(err, utils.InvalidSigningMethod) {
			ctx.Status(http.StatusUnauthorized)
			return
		} else if err != nil {
			m.logger.Error("error occured while parsing token", zap.Error(err), zap.String("token", token))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
			return
		}

		session, werr := m.authService.GetSession(ctx, sessionId)
		if werr != nil {
			m.logger.Error("error occured while authenticating request", zap.Error(werr.Err), zap.String("session_id", sessionId))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("session", session)
	}
}

func (m *Middleware) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
