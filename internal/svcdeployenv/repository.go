package svcdeployenv

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	deploy *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		deploy: client.Database("Capstone").Collection("service_deployments"),
	}
}

func (r *Repository) GetAllServiceDeploymentEnvs(ctx context.Context) ([]models.ServiceDeploymentEnv, error) {
	cur, err := r.deploy.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	svcDeploys := make([]models.ServiceDeploymentEnv, 0)

	for cur.Next(ctx) {
		var dto ServiceDeploymentEnvDTO

		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		svcDeploys = append(svcDeploys, DTOToServiceDeploymentEnv(dto))
	}

	return svcDeploys, nil
}

func (r *Repository) CreateServiceDeploymentEnv(ctx context.Context, dto CreateServiceDeploymentEnvDTO) error {
	_, err := r.deploy.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteServiceDeploymentEnv(ctx context.Context, filter any) (int64, error) {
	result, err := r.deploy.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
