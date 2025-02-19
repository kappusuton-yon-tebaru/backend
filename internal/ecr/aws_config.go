package ecr

import (
	"os"
	
	// "github.com/kappusuton-yon-tebaru/backend/internal/config"
)

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

func LoadAWSConfig() AWSConfig {
	return AWSConfig{
		Region:          os.Getenv("AWS_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}