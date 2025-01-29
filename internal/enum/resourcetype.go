package enum

type ResourceType string

const (
	ResourceTypeOrganization ResourceType = "organization"
	ResourceTypeProjectSpace ResourceType = "project_space"
	ResourceTypeProject ResourceType = "project"

)