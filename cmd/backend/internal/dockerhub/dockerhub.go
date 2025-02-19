package dockerhub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
)

type DockerHubRepository struct {
	apiURL string
	token  string
}

func NewDockerHubRepository(cfg *config.Config) *DockerHubRepository {
	return &DockerHubRepository{
		apiURL: "https://hub.docker.com/v2",
		token:  cfg.DockerHub.Token,
	}
}

func (r *DockerHubRepository) GetImages(namespace string, repoName string) ([]DockerHubImageResponse, error) {
	url := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags", r.apiURL, namespace, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var result struct {
		Results []struct {
			Count int    `json:"count"`
			Next string `json:"next"`
			Previous string `json:"previous"`
			Results struct {
				Name string `json:"name"`
				LastUpdated string `json:"last_updated"`
				Status string `json:"status"`
			} `json:"results"`
		} `json:"results"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var images []DockerHubImageResponse
	for _, resp := range result.Results {
		if resp.Results.Status != "active" {
			continue
		}

		images = append(images, DockerHubImageResponse{
			ImageTag: resp.Results.Name,
			LastUpdated: resp.Results.LastUpdated,
		})
	}
	return images, nil
	
}