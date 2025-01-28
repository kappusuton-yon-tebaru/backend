package models

type ServiceDeployment struct {
	Id        string `json:"id"`
	JobId     string `json:"job_id"`
	ProjectId string `json:"project_id"`
	ImageId   string `json:"image_id"`
}
