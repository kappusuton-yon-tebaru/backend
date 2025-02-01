package mongodb

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoUri))
	defer func() {
		err = client.Disconnect(context.Background())
	}()

	return client, err
}
