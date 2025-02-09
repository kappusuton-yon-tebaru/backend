package job

import (
	"context"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	client *mongo.Client
	job    *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
		job:    client.Database("Capstone").Collection("jobs"),
	}
}

func (r *Repository) GetAllJobParents(ctx context.Context) ([]JobParentDTO, error) {
	pipeline := []map[string]any{
		{
			"$lookup": map[string]any{
				"from":         "jobs",
				"localField":   "_id",
				"foreignField": "parent_job_id",
				"as":           "jobs",
			},
		},
		{
			"$match": map[string]any{
				"parent_job_id": bson.NilObjectID,
			},
		},
		{
			"$project": map[string]any{
				"$id":        true,
				"created_at": true,
				"jobs":       true,
			},
		},
	}

	cur, err := r.job.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	jobParents := make([]JobParentDTO, 0)

	for cur.Next(ctx) {
		var dto JobParentDTO

		err := cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		jobParents = append(jobParents, dto)
	}

	return jobParents, nil
}

func (r *Repository) GetAllJobsByParentId(ctx context.Context, id bson.ObjectID) ([]models.Job, error) {
	filter := map[string]any{
		"parent_job_id": id,
	}

	cur, err := r.job.Find(ctx, filter)
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

func (r *Repository) CreateGroupJobs(ctx context.Context, dtos []CreateJobDTO) (CreateGroupJobsResponse, error) {
	session, err := r.client.StartSession()
	if err != nil {
		return CreateGroupJobsResponse{}, err
	}

	var parentId bson.ObjectID

	defer session.EndSession(ctx)

	results, err := session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		now := time.Now()

		parentJob := CreateJobDTO{
			JobParentId: bson.NilObjectID,
			CreatedAt:   now,
		}

		result, err := r.job.InsertOne(ctx, parentJob)
		if err != nil {
			return nil, err
		}

		parentId = result.InsertedID.(bson.ObjectID)

		for i := range len(dtos) {
			dtos[i].JobParentId = parentId
			dtos[i].CreatedAt = now
		}

		results, err := r.job.InsertMany(ctx, dtos)
		if err != nil {
			return nil, err
		}

		return results.InsertedIDs, err
	})

	if err != nil {
		return CreateGroupJobsResponse{}, err
	}

	ids := []string{}
	for _, result := range results.([]interface{}) {
		id := result.(bson.ObjectID)
		ids = append(ids, id.Hex())
	}

	return CreateGroupJobsResponse{
		ParentId: parentId.Hex(),
		JobIds:   ids,
	}, nil
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
