package job

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	job *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		job: client.Database("Capstone").Collection("jobs"),
	}
}

func (r *Repository) GetAllJobs(ctx context.Context) ([]models.Job, error) {
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

func (r *Repository) CreateJob(ctx context.Context, dto CreateJobDTO) (string, error) {
	result, err := r.job.InsertOne(ctx, dto)
	if err != nil {
		return primitive.NilObjectID.Hex(), err
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) CreateGroupJobs(ctx context.Context, dtos []CreateJobDTO) ([]string, error) {
	results, err := r.job.InsertMany(ctx, dtos)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	for _, result := range results.InsertedIDs {
		id := result.(bson.ObjectID)
		ids = append(ids, id.Hex())
	}

	return ids, nil
}

func (r *Repository) UpdateJobStatus(ctx context.Context, jobId string, jobStatus string) error {
	id, err := bson.ObjectIDFromHex(jobId)
	if err != nil {
		return err
	}

	update := map[string]any{
		"$set": map[string]any{
			"job_status": jobStatus,
		},
	}

	_, err = r.job.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteJob(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.job.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
