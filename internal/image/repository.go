package image

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	image *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		image: client.Database("Capstone").Collection("images"),
	}
}

func (r *Repository) GetAllImages(ctx context.Context) ([]models.Image, error) {
	cur, err := r.image.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	imgs := make([]models.Image, 0)

	for cur.Next(ctx) {
		var img ImageDTO

		err = cur.Decode(&img)
		if err != nil {
			return nil, err
		}

		imgs = append(imgs, DTOToImage(img))
	}

	return imgs, nil
}

func (r *Repository) CreateImage(ctx context.Context, img CreateImageDTO) error {
	_, err := r.image.InsertOne(ctx, img)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteImage(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.image.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
