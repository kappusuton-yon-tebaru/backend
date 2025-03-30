package logging

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	log *mongo.Collection
}

func NewRepository(config *config.Config, db *mongo.Database) (*Repository, error) {
	repo := &Repository{
		log: db.Collection("logs"),
	}

	if err := repo.ensureIndexes(context.Background(), config.PodLogger.LogExpiresInSecond); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) ensureIndexes(ctx context.Context, expiresIn int32) error {
	index := mongo.IndexModel{
		Keys: map[string]any{
			"timestamp": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(expiresIn),
	}

	_, err := r.log.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) BatchInsertLog(ctx context.Context, dtos []InsertLogDTO) error {
	_, err := r.log.InsertMany(ctx, dtos)
	if err != nil {
		return err
	}

	return nil
}
