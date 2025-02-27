package projectrepository

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepositoryDTO struct {
	Id               bson.ObjectID                      `bson:"_id"`
	GitRepoUrl       string                             `bson:"git_repo_url"`
	ProjectId        bson.ObjectID                      `bson:"project_id"`
	RegistryProvider *regproviders.RegistryProvidersDTO `bson:"registry_provider"`
}

type CreateProjectRepositoryDTO struct {
	GitRepoUrl string `json:"git_repo_url" bson:"git_repo_url"`
	// ProjectId  bson.ObjectID `bson:"project_id"`
}

func DTOToProjectRepository(projrepo ProjectRepositoryDTO) models.ProjectRepository {
	regProvider := (*models.RegistryProviders)(nil)
	if projrepo.RegistryProvider != nil {
		regProvider = utils.Pointer(regproviders.DTOToRegistryProviders(*projrepo.RegistryProvider))
	}

	return models.ProjectRepository{
		Id:               projrepo.Id.Hex(),
		GitRepoUrl:       projrepo.GitRepoUrl,
		ProjectId:        projrepo.ProjectId.Hex(),
		RegistryProvider: regProvider,
	}
}

type UpdateProjectRepositryDTO struct {
	RegistryProviderId bson.ObjectID `bson:"registry_provider_id"`
}
