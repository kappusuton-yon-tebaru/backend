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

func (s *Service) GetECRImages(repoURI string, serviceName string, search string, limit int, page int) ([]ECRImageResponse, error) {
	images, err := s.repo.GetImages(repoURI)
	if err != nil {
		return nil, err
	}

	var response []ECRImageResponse
	for _, img := range images {
		if strings.Contains(img, serviceName) && strings.Contains(img, search) {
			response = append(response, ECRImageResponse{
				ImageTag: img,
			})
		}
	}

	var paginationResponse []ECRImageResponse
	for i := (page - 1) * limit; i < page * limit && i < len(response); i++ {
		paginationResponse = append(paginationResponse, response[i])
	}

	return paginationResponse, nil
}
