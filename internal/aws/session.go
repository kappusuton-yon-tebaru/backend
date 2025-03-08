package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Session struct {
	*session.Session
}

func NewSession(cred SessionCredentials) (Session, error) {
	s, err := session.NewSession(&aws.Config{
		Region:      &cred.Region,
		Credentials: credentials.NewStaticCredentials(cred.Key, cred.Secret, ""),
	})

	return Session{s}, err
}
