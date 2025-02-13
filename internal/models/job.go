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
	JobDuration int            `json:"job_duration,omitempty"`
	JsonLogs    string         `json:"json_logs,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}
