package logging

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InsertLogDTO struct {
	Timestamp time.Time         `bson:"timestamp"`
	Log       string            `bson:"log"`
	Attribute map[string]string `bson:"attribute"`
}

type LogDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	Timestamp time.Time     `bson:"timestamp"`
	Log       string        `bson:"log"`
}

func DTOToLog(log LogDTO) models.Log {
	return models.Log{
		Id:        log.Id.Hex(),
		Log:       log.Log,
		Timestamp: log.Timestamp,
	}
}
