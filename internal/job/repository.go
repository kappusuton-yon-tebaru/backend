package job

import (
	"context"
	"errors"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db  *mongo.Database
	job *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db:  db,
		job: db.Collection("jobs"),
	}
}

func (r *Repository) GetAllJobParents(ctx context.Context, queryParam query.QueryParam) (models.Paginated[JobParentDTO], error) {
	pipeline := utils.NewFilterAggregationPipeline(queryParam,
		[]map[string]any{
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
				"$lookup": map[string]any{
					"from":         "resources",
					"localField":   "project_id",
					"foreignField": "_id",
					"as":           "project",
				},
			},
			{
				"$unwind": map[string]any{
					"path":                       "$project",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$project": map[string]any{
					"$id":        true,
					"created_at": true,
					"jobs":       true,
					"project":    true,
				},
			},
		},
	)

	cur, err := r.job.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Paginated[JobParentDTO]{}, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return models.Paginated[JobParentDTO]{}, errors.New("not found")
	}

	var dto models.Paginated[JobParentDTO]
	err = cur.Decode(&dto)
	if err != nil {
		return models.Paginated[JobParentDTO]{}, err
	}

	return dto, nil
}

func (r *Repository) GetAllJobsByParentId(ctx context.Context, id bson.ObjectID, queryParam query.QueryParam) (models.Paginated[JobDTO], error) {
	pipeline := utils.NewFilterAggregationPipeline(queryParam, []map[string]any{
		{
			"$match": map[string]any{
				"parent_job_id": id,
			},
		},
		{
			"$lookup": map[string]any{
				"from":         "resources",
				"localField":   "project_id",
				"foreignField": "_id",
				"as":           "project",
			},
		},
		{
			"$unwind": map[string]any{
				"path":                       "$project",
				"preserveNullAndEmptyArrays": true,
			},
		},
	})

	cur, err := r.job.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Paginated[JobDTO]{}, err
	}

	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return models.Paginated[JobDTO]{}, errors.New("not found")
	}

	var dto models.Paginated[JobDTO]
	err = cur.Decode(&dto)
	if err != nil {
		return models.Paginated[JobDTO]{}, err
	}

	return dto, nil
}

func (r *Repository) CreateJob(ctx context.Context, dto CreateJobDTO) (string, error) {
	result, err := r.job.InsertOne(ctx, dto)
	if err != nil {
		return primitive.NilObjectID.Hex(), err
	}

	id := result.InsertedID.(bson.ObjectID)

	return id.Hex(), nil
}

func (r *Repository) CreateGroupJobs(ctx context.Context, dto CreateJobGroupDTO) (CreateGroupJobsResponse, error) {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return CreateGroupJobsResponse{}, err
	}

	var parentId bson.ObjectID

	defer session.EndSession(ctx)

	results, err := session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		now := time.Now()

		parentJob := CreateJobDTO{
			JobParentId: bson.NilObjectID,
			ProjectId:   dto.ProjectId,
			CreatedAt:   now,
		}

		result, err := r.job.InsertOne(ctx, parentJob)
		if err != nil {
			return nil, err
		}

		parentId = result.InsertedID.(bson.ObjectID)

		for i := range len(dto.Jobs) {
			dto.Jobs[i].JobParentId = parentId
			dto.Jobs[i].CreatedAt = now
		}

		results, err := r.job.InsertMany(ctx, dto.Jobs)
		if err != nil {
			return nil, err
		}

		return results.InsertedIDs, err
	})

	if err != nil {
		return CreateGroupJobsResponse{}, err
	}

	ids := []string{}
	for _, result := range results.([]any) {
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
