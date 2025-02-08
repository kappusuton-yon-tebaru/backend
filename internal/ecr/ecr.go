package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ECRRepository struct {
	client *ecr.ECR
}

func NewECRRepository(cfg AWSConfig) *ECRRepository {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}))
	return &ECRRepository{
		client: ecr.New(session),
	}
}

func (r *ECRRepository) GetImages(repoName string) ([]string, error) {
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
		Filter: &ecr.ListImagesFilter{
			TagStatus: aws.String("TAGGED"),
		},
	}

	result, err := r.client.ListImages(input)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, image := range result.ImageIds {
		images = append(images, *image.ImageTag)
	}

	return images, nil
}

