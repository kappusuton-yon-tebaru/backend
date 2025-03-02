package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtToken(sessionId string, expiresIn int, jwtSecret string) (string, error) {
	now := time.Now()
	issuedTime := now
	expiresTime := now.Add(time.Duration(expiresIn) * time.Second)

	fmt.Println(issuedTime, expiresTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        sessionId,
		IssuedAt:  jwt.NewNumericDate(issuedTime),
		ExpiresAt: jwt.NewNumericDate(expiresTime),
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string, jwtSecret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claim, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("error parsing jwt")
	}

	return claim.ID, nil
}
