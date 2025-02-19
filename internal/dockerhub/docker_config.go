package dockerhub

import "os"

type DockerConfig struct {
	DockerHubToken string
}

func LoadDockerConfig() DockerConfig {
	return DockerConfig{
		DockerHubToken: os.Getenv("DOCKERHUB_TOKEN"),
	}
}
