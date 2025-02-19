package dockerhub

import "strings"

type Service struct {
	repo *DockerHubRepository
}

func NewService(repo *DockerHubRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetDockerHubImages(namespace string, repoName string, serviceName string) ([]DockerHubImageResponse, error) {
	images, err := s.repo.GetImages(namespace, repoName)
	if err != nil {
		return nil, err
	}

	var response []DockerHubImageResponse
	for _, img := range images {
		if strings.Contains(img.ImageTag, serviceName) {
			response = append(response, DockerHubImageResponse{
				ImageTag:    img.ImageTag,
				LastUpdated: img.LastUpdated,
			})
		}
	}

	return response, nil
}
