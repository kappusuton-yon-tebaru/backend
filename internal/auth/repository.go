package auth

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	db      *mongo.Database
	session *mongo.Collection
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	repo := &Repository{
		db:      db,
		session: db.Collection("sessions"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) ensureIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys: map[string]any{
			"expires_at": -1,
		},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := r.session.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateSession(ctx context.Context, dto CreateSessionDTO) (string, error) {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		filter := map[string]any{
			"user_id": dto.UserId,
		}

		_, err := r.session.DeleteOne(ctx, filter)
		if err != nil {
			return nil, err
		}

		result, err := r.session.InsertOne(ctx, dto)
		if err != nil {
			return "", err
		}

		sessionId := result.InsertedID.(bson.ObjectID)
		return sessionId.Hex(), nil
	})

	sessionId := result.(string)
	return sessionId, nil
}

func (r *Repository) GetSessionById(ctx context.Context, sessionId bson.ObjectID) (models.Session, error) {
	filter := map[string]any{
		"_id": sessionId,
	}

	result := r.session.FindOne(ctx, filter)
	if result.Err() != nil {
		return models.Session{}, result.Err()
	}

	var dto SessionDTO
	if err := result.Decode(&dto); err != nil {
		return models.Session{}, err
	}

	return DTOToSession(dto), nil
}
