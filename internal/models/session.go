package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	Id        string
	UserId    string
	ExpiresAt time.Time
}

type SessionClaim struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}
