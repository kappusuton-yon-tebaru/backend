package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/spf13/viper"
)

type AgentConfig struct {
	Port           int    `mapstructure:"AGENT_PORT"`
	WorkerImageUri string `mapstructure:"WORKER_IMAGE_URI"`
}

type BackendConfig struct {
	Port          int    `mapstructure:"BACKEND_PORT"`
	AgentEndpoint string `mapstructure:"BACKEND_AGENT_ENDPOINT"`
}

type ConsumerConfig struct {
	QueueUri         string `mapstructure:"CONSUMER_QUEUE_URI"`
	OrganizationName string `mapstructure:"CONSUMER_ORGANIZATION_NAME"`
}

type PodLogger struct {
	Mode               enum.PodLoggerMode `mapstructure:"POD_LOGGER_MODE"`
	LogExpiresInSecond int32              `mapstructure:"POD_LOGGER_LOG_EXPIRES_IN_SECOND"`
}

type DockerHubConfig struct {
	Token string `mapstructure:"DOCKERHUB_TOKEN"`
}

type Config struct {
	Agent                  AgentConfig     `mapstructure:",squash"`
	Backend                BackendConfig   `mapstructure:",squash"`
	DockerHub              DockerHubConfig `mapstructure:",squash"`
	MongoUri               string          `mapstructure:"MONGO_URI"`
	InCluster              bool            `mapstructure:"IN_CLUSTER"`
	Development            bool            `mapstructure:"DEVELOPMENT"`
	ConsumerConfig         ConsumerConfig  `mapstructure:",squash"`
	KubeNamespace          string          `mapstructure:"KUBE_NAMESPACE"`
	MongoDatabaseName      string          `mapstructure:"MONGO_DATABASE_NAME"`
	ClientID               string          `mapstructure:"GITHUB_CLIENT_ID"`
	ClientSecret           string          `mapstructure:"GITHUB_CLIENT_SECRET"`
	SessionExpiresInSecond int             `mapstructure:"SESSION_EXPIRES_IN_SECOND"`
	JwtSecret              string          `mapstructure:"JWT_SECRET"`
	PodLogger              PodLogger       `mapstructure:",squash"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	err := bindStruct(Config{})
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindStruct(s any) error {
	ct := reflect.TypeOf(s)

	if ct.Kind() != reflect.Struct {
		return fmt.Errorf("listStructKeys: %v is not a struct", ct)
	}

	for i := range ct.NumField() {
		field := ct.Field(i)
		tag := field.Tag.Get("mapstructure")

		if field.Type.Kind() == reflect.Struct {
			err := bindStruct(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return err
			}
		} else {
			if err := viper.BindEnv(strings.Split(tag, ",")[0]); err != nil {
				return err
			}
		}
	}

	return nil
}
