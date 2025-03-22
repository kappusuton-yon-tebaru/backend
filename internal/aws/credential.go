package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type ClientBuilder struct {
	aws.Config
}

func NewConfig(cred models.ECRCredential, region string) ClientBuilder {
	return ClientBuilder{aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(cred.AccessKey, cred.SecretAccessKey, ""),
	}}
}
