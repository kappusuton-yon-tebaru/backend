package projectrepository

type UpdateRegistryProviderDTO struct {
	RegistryProviderId string `json:"registry_provider_id" validate:"required,min=1"`
}
