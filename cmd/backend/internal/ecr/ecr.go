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
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
)

type ECRRepository struct {
	client       *ecr.ECR
	publicClient *ecrpublic.Client
}

func NewECRRepository(cfg *config.Config) *ECRRepository {
	privateSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.ECR.Region),
		Credentials: credentials.NewStaticCredentials(cfg.ECR.AccessKey, cfg.ECR.SecretKey, ""),
	}))
	publicConfig, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion(cfg.ECR.Region))
	if err != nil {
		panic(err)
	}

	return &ECRRepository{
		client:       ecr.New(privateSession),
		publicClient: ecrpublic.NewFromConfig(publicConfig),
	}
}

func IsPublicRepo(repoURI string) bool {
	return strings.HasPrefix(repoURI, "public.ecr.aws")
}

func GetRepoName(repoURI string) string {
	return repoURI[strings.LastIndex(repoURI, "/")+1:]
}

func (r *ECRRepository) GetImages(repoURI string, serviceName string, queryParam query.QueryParam) (PaginatedECRImages, error) {
	isPublic := IsPublicRepo(repoURI)
	repoName := GetRepoName(repoURI)

	offset := (queryParam.Pagination.Page - 1) * queryParam.Pagination.Limit;
	end := offset + queryParam.Pagination.Limit;

	if isPublic {
		input := &ecrpublic.DescribeImageTagsInput{
			RepositoryName: aws.String(repoName),
		}

		result, err := r.publicClient.DescribeImageTags(context.TODO(), input)
		if err != nil {
			return PaginatedECRImages{}, err
		}

		var images []ECRImageResponse
		for _, image := range result.ImageTagDetails {
			if strings.Contains(*image.ImageTag, serviceName) && strings.Contains(*image.ImageTag, queryParam.QueryFilter.Query) {
				images = append(images, ECRImageResponse{
					*image.ImageTag,
				})
			}
		}

		return PaginatedECRImages{
			Page:  queryParam.Pagination.Page,
			Limit: queryParam.Pagination.Limit,
			Total: len(images),
			Data:  images[offset : end],
		}, nil
	} else {
		input := &ecr.ListImagesInput{
			RepositoryName: aws.String(repoName),
			Filter: &ecr.ListImagesFilter{
				TagStatus: aws.String("TAGGED"),
			},
		}
		result, err := r.client.ListImages(input)
		if err != nil {
			return PaginatedECRImages{}, err
		}

		var images []ECRImageResponse
		for _, image := range result.ImageIds {
			if strings.Contains(*image.ImageTag, serviceName) && strings.Contains(*image.ImageTag, queryParam.QueryFilter.Query) {
				images = append(images, ECRImageResponse{
					*image.ImageTag,
				})
			}
		}

		return PaginatedECRImages{
			Page:  queryParam.Pagination.Page,
			Limit: queryParam.Pagination.Limit,
			Total: len(images),
			Data:  images[offset : end],
		}, nil
	}

}
