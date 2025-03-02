package job

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobDTO struct {
	Id          bson.ObjectID        `bson:"_id"`
	JobType     string               `bson:"job_type"`
	JobStatus   enum.JobStatus       `bson:"job_status"`
	ServiceName string               `bson:"service_name"`
	Project     resource.ResourceDTO `bson:"project"`
	CreatedAt   time.Time            `bson:"created_at"`
}

type CreateJobDTO struct {
	JobParentId bson.ObjectID `bson:"parent_job_id"`
	JobType     string        `bson:"job_type,omitempty"`
	JobStatus   string        `bson:"job_status,omitempty"`
	ProjectId   bson.ObjectID `bson:"project_id"`
	ServiceName string        `bson:"service_name,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"`
}

type CreateJobGroupDTO struct {
	ProjectId bson.ObjectID
	Jobs      []CreateJobDTO
}

func DTOToJob(job JobDTO) models.Job {
	return models.Job{
		Id:          job.Id.Hex(),
		JobType:     job.JobType,
		JobStatus:   job.JobStatus,
		CreatedAt:   job.CreatedAt,
		ServiceName: job.ServiceName,
		Project: models.JobProject{
			Id:   job.Project.Id.Hex(),
			Name: job.Project.ResourceName,
		},
	}
}

type JobParentDTO struct {
	Id        bson.ObjectID        `bson:"_id"`
	Jobs      []JobDTO             `bson:"jobs"`
	Project   resource.ResourceDTO `bson:"project"`
	CreatedAt time.Time            `bson:"created_at"`
}

type CreateGroupJobsResponse struct {
	ParentId string
	JobIds   []string
}
