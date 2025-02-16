package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type AgentConfig struct {
	Port           int    `mapstructure:"AGENT_PORT"`
	WorkerImageUri string `mapstructure:"AGENT_WORKER_IMAGE_URI"`
}

type BackendConfig struct {
	Port          int    `mapstructure:"BACKEND_PORT"`
	AgentEndpoint string `mapstructure:"BACKEND_AGENT_ENDPOINT"`
}

type BuilderConfig struct {
	QueueUri  string `mapstructure:"BUILDER_QUEUE_URI"`
	QueueName string `mapstructure:"BUILDER_QUEUE_NAME"`
}

type Config struct {
	Agent             AgentConfig   `mapstructure:",squash"`
	Backend           BackendConfig `mapstructure:",squash"`
	BuilderConfig     BuilderConfig `mapstructure:",squash"`
	InCluster         bool          `mapstructure:"IN_CLUSTER"`
	Development       bool          `mapstructure:"DEVELOPMENT"`
	KubeNamespace     string        `mapstructure:"KUBE_NAMESPACE"`
	MongoUri          string        `mapstructure:"MONGO_URI"`
	MongoDatabaseName string        `mapstructure:"MONGO_DATABASE_NAME"`
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

func bindStruct(s interface{}) error {
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
