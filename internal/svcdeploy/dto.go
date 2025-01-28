package svcdeploy

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ServiceDeploymentDTO struct {
	Id        bson.ObjectID `bson:"_id"`
	JobId     bson.ObjectID `bson:"job_id"`
	ProjectId bson.ObjectID `bson:"project_id"`
	ImageId   bson.ObjectID `bson:"image_id"`
}

type CreateServiceDeploymentDTO struct {
	JobId     bson.ObjectID `bson:"job_id"`
	ProjectId bson.ObjectID `bson:"project_id"`
	ImageId   bson.ObjectID `bson:"image_id"`
}

func DTOToServiceDeployment(svcDeploy ServiceDeploymentDTO) models.ServiceDeployment {
	return models.ServiceDeployment{
		Id:        svcDeploy.Id.Hex(),
		JobId:     svcDeploy.JobId.Hex(),
		ProjectId: svcDeploy.ProjectId.Hex(),
		ImageId:   svcDeploy.ImageId.Hex(),
	}
}
