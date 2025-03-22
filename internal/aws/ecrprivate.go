package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type ECRPrivate struct {
	client *ecr.Client
}

func (c ClientBuilder) ECRPrivate() ECRPrivate {
	return ECRPrivate{ecr.NewFromConfig(c.Config)}
}

func (e ECRPrivate) ListImages(ctx context.Context, repoName string) ([]types.ImageIdentifier, error) {
	result, err := e.client.ListImages(ctx, &ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
		Filter: &types.ListImagesFilter{
			TagStatus: types.TagStatusTagged,
		},
	})
	if err != nil {
		return nil, err
	}

	return result.ImageIds, nil
}
