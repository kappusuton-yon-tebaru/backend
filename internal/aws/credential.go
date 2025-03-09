package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type SessionBuilder struct {
	*session.Session
	Err error
}

func NewSession(cred models.ECRCredential, region string) SessionBuilder {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cred.AccessKey, cred.SecretAccessKey, ""),
	})

	return SessionBuilder{s, err}
}
