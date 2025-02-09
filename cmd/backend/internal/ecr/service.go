package ecr

import "strings"

type Service struct {
	repo *ECRRepository
}

func NewService(repo *ECRRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetECRImages(repoName string, serviceName string) ([]ECRImageResponse, error) {
	images, err := s.repo.GetImages(repoName)
	if err != nil {
		return nil, err
	}

	var response []ECRImageResponse
	for _, img := range images {
		if strings.Contains(img, serviceName) {
			response = append(response, ECRImageResponse{
				ImageTag: img,
			})
		}
	}

	return response, nil
}
