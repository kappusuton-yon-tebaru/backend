package enum

type JobStatus string

const (
	JobStatusPending JobStatus = "pending"
	JobStatusRunning JobStatus = "running"
	JobStatusSuccess JobStatus = "success"
	JobStatusFailed  JobStatus = "failed"
)