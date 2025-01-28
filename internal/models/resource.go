package models

type ResourceType string

const (
	Organization ResourceType = "organization"
	ProjectSpace ResourceType = "project_space"
	Project      ResourceType = "project"
)

type Resource struct {
	Id           string `json:"id"`
	ResourceName string `json:"resource_name"`
	ResourceType ResourceType
}
