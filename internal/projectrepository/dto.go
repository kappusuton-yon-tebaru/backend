package projectrepository

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepositoryDTO struct {
	Id               bson.ObjectID                     `bson:"_id"`
	GitRepoUrl       string                            `bson:"git_repo_url"`
	ProjectId        bson.ObjectID                     `bson:"project_id"`
	RegistryProvider regproviders.RegistryProvidersDTO `bson:"registry_provider"`
}

type CreateProjectRepositoryDTO struct {
	GitRepoUrl string        `bson:"git_repo_url"`
	ProjectId  bson.ObjectID `bson:"project_id"`
}

func DTOToProjectRepository(projrepo ProjectRepositoryDTO) models.ProjectRepository {
	return models.ProjectRepository{
		Id:               projrepo.Id.Hex(),
		GitRepoUrl:       projrepo.GitRepoUrl,
		ProjectId:        projrepo.ProjectId.Hex(),
		RegistryProvider: regproviders.DTOToRegistryProviders(projrepo.RegistryProvider),
	}
}

type UpdateProjectRepositryDTO struct {
	RegistryProviderId bson.ObjectID `bson:"registry_provider_id"`
}
