package build

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type BuildContext struct {
	Id                 string                    `json:"id"`
	RepoUrl            string                    `json:"repo_url"`
	Dockerfile         string                    `json:"dockerfile"`
	RepoRoot           string                    `json:"repo_root"`
	Destination        string                    `json:"destination"`
	RegistryType       enum.RegistryProviderType `json:"registry_type"`
	RegistryCredential interface{}               `json:"registry_credential"`
}
