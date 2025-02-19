package enum

type ResourceType string

const (
	ResourceTypeOrganization ResourceType = "ORGANIZATION"
	ResourceTypeProjectSpace ResourceType = "PROJECT_SPACE"
	ResourceTypeProject      ResourceType = "PROJECT"
)

func IsValidResourceType(t ResourceType) bool {
	switch t {
	case ResourceTypeOrganization, ResourceTypeProjectSpace, ResourceTypeProject:
		return true
	default:
		return false
	}
}
