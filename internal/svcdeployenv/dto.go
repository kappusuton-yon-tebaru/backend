package svcdeployenv

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type ServiceDeploymentEnvDTO struct {
	Id                  string `bson:"_id"`
	ServiceDeploymentId string `bson:"service_deployment_id"`
	EnvType             string `bson:"env_type"`
	Key                 string `bson:"key"`
	Value               string `bson:"value"`
}

type CreateServiceDeploymentEnvDTO struct {
	ServiceDeploymentId string `bson:"service_deployment_id"`
	EnvType             string `bson:"env_type"`
	Key                 string `bson:"key"`
	Value               string `bson:"value"`
}

func DTOToServiceDeploymentEnv(svcDeployEnv ServiceDeploymentEnvDTO) models.ServiceDeploymentEnv {
	return models.ServiceDeploymentEnv{
		Id:                  svcDeployEnv.Id,
		ServiceDeploymentId: svcDeployEnv.ServiceDeploymentId,
		EnvType:             svcDeployEnv.EnvType,
		Key:                 svcDeployEnv.Key,
		Value:               svcDeployEnv.Value,
	}
}
