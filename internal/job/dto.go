package job

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobDTO struct {
	Id          bson.ObjectID	`bson:"_id"`
	JobType     string    		`bson:"job_type"`
	JobStatus   enum.JobStatus	`bson:"job_status"`
	JobDuration int       		`bson:"job_duration"`
	JsonLogs    string    		`bson:"json_logs"`
}

type CreateJobDTO struct {
	JobType     string    		`bson:"job_type"`
	JobStatus   string	 		`bson:"job_status"`
	JobDuration int       		`bson:"job_duration"`
	JsonLogs    string    		`bson:"json_logs"`
}

func DTOToJob(job JobDTO) models.Job {
	return models.Job{
		Id:          job.Id.Hex(),
		JobType:     job.JobType,
		JobStatus:   job.JobStatus,
		JobDuration: job.JobDuration,
		JsonLogs:    job.JsonLogs,
	}
}