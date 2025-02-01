package job

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository struct {
	job *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		job: client.Database("Capstone").Collection("jobs"),
	}
}

func (r* Repository) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	cur, err := r.job.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	jobs := make([]models.Job, 0)

	for cur.Next(ctx) {
		var dto JobDTO

		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, DTOToJob(dto))
	}

	return jobs, nil
}

func (r *Repository) CreateJob(ctx context.Context, dto CreateJobDTO) (any, error) {
	result, err := r.job.InsertOne(ctx, dto)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID, nil
}

func (r * Repository) DeleteJob(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.job.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}