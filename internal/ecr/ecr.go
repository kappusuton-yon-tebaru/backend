package ecr

import (
	"fmt"
	"os"

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
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}))
	return &ECRRepository{
		client: ecr.New(session),
	}
}

func (r *ECRRepository) GetImages(repoName string) ([]string, error) {
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
	}

	result, err := r.client.ListImages(input)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, image := range result.ImageIds {
		images = append(images, fmt.Sprintf("%s:%s", *image.ImageDigest, *image.ImageTag))
	}

	return images, nil
}

func (r *ECRRepository) PushImage(repoName, imageName, tag string) (string, error) {
	input := &ecr.PutImageInput{
		RepositoryName: aws.String(repoName),
		ImageManifest:  aws.String(fmt.Sprintf(`{"image_name": "%s", "tag": "%s"}`, imageName, tag)),
		ImageTag:       aws.String(tag),
	}

	result, err := r.client.PutImage(input)
	if err != nil {
		return "", err
	}

	return *result.Image.RegistryId, nil
}

