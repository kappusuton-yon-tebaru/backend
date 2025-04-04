package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

var InvalidSigningMethod = errors.New("invalid signing method")

func CreateJwtToken(sessionId string, expiresIn int, jwtSecret string) (string, error) {
	now := time.Now()
	issuedTime := now
	expiresTime := now.Add(time.Duration(expiresIn) * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.SessionClaim{
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedTime),
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string, jwtSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.SessionClaim{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, InvalidSigningMethod
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claim, ok := token.Claims.(*models.SessionClaim)
	if !ok {
		return "", errors.New("error parsing jwt")
	}

	return claim.SessionId, nil
}
