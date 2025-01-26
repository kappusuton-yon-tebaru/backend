package image

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ImageDTO struct {
	Id                 bson.ObjectID `bson:"_id"`
	ImageName          string        `bson:"image_name"`
	ProjectId          bson.ObjectID `bson:"project_id"`
	JobId              bson.ObjectID `bson:"job_id"`
	RegistryProviderId bson.ObjectID `bson:"registry_provider_id"`
	Version            string        `bson:"version"`
	IsDeleted          bool          `bson:"is_deleted"`
}

func (img ImageDTO) ToImage() models.Image {
	return models.Image{
		Id:                 img.Id.Hex(),
		ImageName:          img.ImageName,
		ProjectId:          img.ProjectId.Hex(),
		JobId:              img.JobId.Hex(),
		RegistryProviderId: img.RegistryProviderId.Hex(),
		Version:            img.Version,
		IsDeleted:          img.IsDeleted,
	}
}
