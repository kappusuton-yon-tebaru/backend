package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManager struct {
	client *secretsmanager.Client
}

func (c ClientBuilder) SecretsManager() SecretsManager {
	return SecretsManager{secretsmanager.NewFromConfig(c.Config)}
}

func (sm SecretsManager) GetSecretValue(ctx context.Context, secretName string, v any) error {
	output, err := sm.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(*output.SecretString), v)
	if err != nil {
		return err
	}

	return nil
}
