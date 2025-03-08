package auth

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateSessionDTO struct {
	UserId    bson.ObjectID `bson:"user_id"`
	ExpiresAt time.Time     `bson:"expires_at"`
}

type SessionDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	UserId    bson.ObjectID `bson:"user_id"`
	ExpiresAt time.Time     `bson:"expires_at"`
}

func DTOToSession(dto SessionDTO) models.Session {
	return models.Session{
		Id:        dto.Id.Hex(),
		UserId:    dto.UserId.Hex(),
		ExpiresAt: dto.ExpiresAt,
	}
}
