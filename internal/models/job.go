package models

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type Job struct {
	Id          string         `json:"id"`
	JobType     string         `json:"job_type"`
	JobStatus   enum.JobStatus `json:"job_status"`
	JobDuration int            `json:"job_duration"`
	JsonLogs    string         `json:"json_logs"`
}
