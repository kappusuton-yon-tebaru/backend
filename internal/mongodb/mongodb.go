package mongodb

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDB(cfg *config.Config) (*mongo.Database, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoUri))
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.MongoDatabaseName), nil
}
