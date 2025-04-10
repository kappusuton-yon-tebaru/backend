package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"errors"
)

func GetUserID(ctx *gin.Context) (string, error) {
	sessionValue, exists := ctx.Get("session")
	if !exists {
		return "", errors.New("no session found")
	}

	session, ok := sessionValue.(models.Session)
	if !ok {
		return "", errors.New("invalid session type")
	}

	return session.UserId, nil
}