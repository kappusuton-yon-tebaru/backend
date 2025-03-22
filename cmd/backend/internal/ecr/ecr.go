package ecr

import (
	"context"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/internal/aws"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
)

type ECRRepository struct{}

func NewECRRepository() *ECRRepository {
	return &ECRRepository{}
}

func IsPublicRepo(repoURI string) bool {
	return strings.HasPrefix(repoURI, "public.ecr.aws")
}

func GetRepoName(repoURI string) string {
	return repoURI[strings.LastIndex(repoURI, "/")+1:]
}

func (r *ECRRepository) GetImages(ctx context.Context, registry models.RegistryProviders, serviceName string, queryParam query.QueryParam) (PaginatedECRImages, error) {
	repoURI := registry.Uri

	isPublic := IsPublicRepo(repoURI)
	repoName := GetRepoName(repoURI)

	if isPublic {
		client := aws.NewConfig(*registry.ECRCredential, registry.ECRCredential.Region).ECRPublic()
		tags, err := client.DescribeImageTags(ctx, repoName)
		if err != nil {
			return PaginatedECRImages{}, err
		}

		var images []ECRImageResponse
		for _, image := range tags {
			if strings.Contains(*image.ImageTag, serviceName) && strings.Contains(*image.ImageTag, queryParam.QueryFilter.Query) {
				images = append(images, ECRImageResponse{
					*image.ImageTag,
				})
			}
		}

		paginatedImages := utils.Paginate(images, queryParam.Pagination.Page, queryParam.Pagination.Limit)

		return PaginatedECRImages{
			Page:  queryParam.Pagination.Page,
			Limit: queryParam.Pagination.Limit,
			Total: len(paginatedImages),
			Data:  paginatedImages,
		}, nil
	} else {
		client := aws.NewConfig(*registry.ECRCredential, registry.ECRCredential.Region).ECRPrivate()
		tags, err := client.ListImages(ctx, repoName)
		if err != nil {
			return PaginatedECRImages{}, err
		}

		var images []ECRImageResponse
		for _, image := range tags {
			if strings.Contains(*image.ImageTag, serviceName) && strings.Contains(*image.ImageTag, queryParam.QueryFilter.Query) {
				images = append(images, ECRImageResponse{
					*image.ImageTag,
				})
			}
		}

		paginatedImages := utils.Paginate(images, queryParam.Pagination.Page, queryParam.Pagination.Limit)

		return PaginatedECRImages{
			Page:  queryParam.Pagination.Page,
			Limit: queryParam.Pagination.Limit,
			Total: len(paginatedImages),
			Data:  paginatedImages,
		}, nil
	}

}
