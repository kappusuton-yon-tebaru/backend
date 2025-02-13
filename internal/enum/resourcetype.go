package enum

type ResourceType string

const (
	ResourceTypeOrganization ResourceType = "organization"
	ResourceTypeProjectSpace ResourceType = "project_space"
	ResourceTypeProject      ResourceType = "project"
)

func IsValidResourceType(t ResourceType) bool {
	switch t {
	case ResourceTypeOrganization, ResourceTypeProjectSpace, ResourceTypeProject:
		return true
	default:
		return false
	}
}
