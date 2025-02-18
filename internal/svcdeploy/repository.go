package svcdeploy

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	deploy *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		deploy: db.Collection("service_deployments"),
	}
}

func (r *Repository) GetAllServiceDeployments(ctx context.Context) ([]models.ServiceDeployment, error) {
	cur, err := r.deploy.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	svcDeploys := make([]models.ServiceDeployment, 0)

	for cur.Next(ctx) {
		var dto ServiceDeploymentDTO

		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		svcDeploys = append(svcDeploys, DTOToServiceDeployment(dto))
	}

	return svcDeploys, nil
}

func (r *Repository) CreateServiceDeployment(ctx context.Context, dto CreateServiceDeploymentDTO) error {
	_, err := r.deploy.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteServiceDeployment(ctx context.Context, filter map[string]any) (int64, error) {
	result, err := r.deploy.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
