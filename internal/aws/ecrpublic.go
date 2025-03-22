package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
)

type ECRPublic struct {
	client *ecrpublic.Client
}

func (c ClientBuilder) ECRPublic() ECRPublic {
	return ECRPublic{ecrpublic.NewFromConfig(c.Config)}
}

func (ecr ECRPublic) DescribeImageTags(ctx context.Context, repoName string) ([]types.ImageTagDetail, error) {
	result, err := ecr.client.DescribeImageTags(ctx, &ecrpublic.DescribeImageTagsInput{
		RepositoryName: aws.String(repoName),
	})
	if err != nil {
		return nil, err
	}

	return result.ImageTagDetails, nil
}
