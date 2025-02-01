package models

type RegistryProviders struct {
	Id		  		string 	`json:"id"`
	Name	  		string 	`json:"name"`
	ProviderType	string	`json:"provider_type"` // enum
	JsonCredential	string	`json:"json_credential"`
	OrganizationId	string	`json:"organization_id"`
}


