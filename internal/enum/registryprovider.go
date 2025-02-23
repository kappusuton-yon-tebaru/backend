package enum

type RegistryProviderType string

const (
	ECR       RegistryProviderType = "ECR"
	DOCKERHUB RegistryProviderType = "DOCKERHUB"
)
