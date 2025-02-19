package ecr

import (
	"context"
	"strings"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
)

type ECRRepository struct {
	client *ecr.ECR
	publicClient *ecrpublic.Client
}

func NewECRRepository(cfg *config.Config) *ECRRepository {
	privateSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.ECR.Region),
		Credentials: credentials.NewStaticCredentials(cfg.ECR.AccessKey, cfg.ECR.SecretKey, ""),
	}))
	publicConfig, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion(cfg.ECR.Region))
	if err != nil {
		panic(err)
	}

	return &ECRRepository{
		client: ecr.New(privateSession),
		publicClient: ecrpublic.NewFromConfig(publicConfig),
	}
}

func IsPublicRepo(repoURI string) bool {
	return strings.HasPrefix(repoURI, "public.ecr.aws")
}

func GetRepoName(repoURI string) string {
	return repoURI[strings.LastIndex(repoURI, "/")+1:]
}

func (r *ECRRepository) GetImages(repoURI string) ([]string, error) {
	isPublic := IsPublicRepo(repoURI)
	repoName := GetRepoName(repoURI)

	if isPublic {
		input := &ecrpublic.DescribeImageTagsInput{
			RepositoryName: aws.String(repoName),
		}

		result, err := r.publicClient.DescribeImageTags(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		var images []string
		for _, image := range result.ImageTagDetails {
			images = append(images, *image.ImageTag)
		}

		return images, nil
	} else {
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

}

