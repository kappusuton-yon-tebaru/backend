package projectenv

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectEnvDTO struct {
	Id        bson.ObjectID `bson:"id"`
	ProjectId bson.ObjectID `bson:"project_id"`
	EnvType   enum.EnvType  `bson:"env_type"`
	Key       string        `bson:"key"`
	Value     string        `bson:"value"`
}

type CreateProjectEnvDTO struct {
	ProjectId bson.ObjectID `bson:"project_id"`
	EnvType   enum.EnvType  `bson:"env_type"`
	Key       string        `bson:"key"`
	Value     string        `bson:"value"`
}

func DTOToProjectEnv(projectEnv ProjectEnvDTO) models.ProjectEnv {
	return models.ProjectEnv{
		Id:        projectEnv.Id.Hex(),
		ProjectId: projectEnv.ProjectId.Hex(),
		EnvType:   projectEnv.EnvType,
		Key:       projectEnv.Key,
		Value:     projectEnv.Value,
	}
}
