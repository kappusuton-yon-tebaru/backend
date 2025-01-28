package projectrepository

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepositoryDTO struct {
	Id         bson.ObjectID `bson:"_id"`
	GitRepoUrl string        `bson:"git_repo_url"`
	ProjectId  bson.ObjectID `bson:"project_id"`
}

func DTOToProjectRepository(projrepo ProjectRepositoryDTO) models.ProjectRepository {
	return models.ProjectRepository{
		Id:         projrepo.Id.Hex(),
		GitRepoUrl: projrepo.GitRepoUrl,
		ProjectId:  projrepo.ProjectId.Hex(),
	}
}
