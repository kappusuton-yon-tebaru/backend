package mongodb

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(cfg *config.Config) (*mongo.Client, error) {
	return mongo.Connect(options.Client().ApplyURI(cfg.MongoUri))
}
