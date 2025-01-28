package models

type Image struct {
	Id                 string `json:"id"`
	ImageName          string `json:"image_name"`
	ProjectId          string `json:"project_id"`
	JobId              string `json:"job_id"`
	RegistryProviderId string `json:"registry_provider_id"`
	Version            string `json:"version"`
	IsDeleted          bool   `json:"is_deleted"`
}
