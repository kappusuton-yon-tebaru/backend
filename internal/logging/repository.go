package logging

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (r *Repository) GetLog(ctx context.Context, queryParam query.QueryParam) ([]models.Log, error) {
	filters := map[string]any{}
	sortDir := 1

	opts := options.Find()
	if queryParam.CursorPagination != nil {
		pagination := *queryParam.CursorPagination
		opts = opts.SetLimit(int64(pagination.Limit))

		switch pagination.Direction {
		case enum.Newer:
			sortDir = 1
		case enum.Older:
			sortDir = -1
		}

		if pagination.Cursor != nil {
			objId, _ := bson.ObjectIDFromHex(*pagination.Cursor)
			switch pagination.Direction {
			case enum.Newer:
				filters["_id"] = map[string]any{"$gt": objId}
			case enum.Older:
				filters["_id"] = map[string]any{"$lt": objId}
			}
		}
	}

	opts = opts.SetSort(map[string]any{"timestamp": sortDir})

	if queryParam.Filter != nil {
		filter := *queryParam.Filter
		filters["attribute"] = filter
	}

	cursor, err := r.log.Find(ctx, filters, opts)
	if err != nil {
		return nil, err
	}

	logs := []models.Log{}
	for cursor.Next(ctx) {
		var dto LogDTO
		err := cursor.Decode(&dto)
		if err != nil {
			return nil, err
		}

		if sortDir == 1 {
			logs = append(logs, DTOToLog(dto))
		} else if sortDir == -1 {
			logs = append([]models.Log{DTOToLog(dto)}, logs...)
		}
	}

	return logs, nil
}
