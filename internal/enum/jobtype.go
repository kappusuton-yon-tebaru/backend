package enum

type JobType string

const (
	JobTypeBuild  JobType = "build"
	JobTypeDeploy JobType = "deploy"
)
