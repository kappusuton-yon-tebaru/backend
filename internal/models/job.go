package models

import (
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
)

type Job struct {
	Id          string         `json:"id"`
	JobParentId string         `json:"job_parent_id,omitempty"`
	JobType     string         `json:"job_type,omitempty"`
	JobStatus   enum.JobStatus `json:"job_status"`
	ServiceName string         `json:"service_name,omitempty"`
	Project     JobProject     `json:"project"`
	CreatedAt   time.Time      `json:"created_at"`
}

type JobProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
