package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManager struct {
	session *secretsmanager.SecretsManager
}

func (s Session) SecretManager() SecretsManager {
	return SecretsManager{secretsmanager.New(s)}
}

func (sm SecretsManager) GetSecretValue(ctx context.Context, secretName string, v any) error {
	output, err := sm.session.GetSecretValueWithContext(ctx, &secretsmanager.GetSecretValueInput{
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
