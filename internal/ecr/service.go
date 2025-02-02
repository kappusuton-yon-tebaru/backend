package ecr

type Service struct {
	repo *ECRRepository
}

func NewService(repo *ECRRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetECRImages(repoName string) ([]ECRImageResponse, error) {
	images, err := s.repo.GetImages(repoName)
	if err != nil {
		return nil, err
	}

	var response []ECRImageResponse
	for _, img := range images {
		response = append(response, ECRImageResponse{
			RepositoryName: repoName,
			ImageDigest:    img,
			ImageTag:       "latest",
		})
	}

	return response, nil
}

func (s *Service) PushECRImage(req PushECRImageRequest) (string, error) {
	id, err := s.repo.PushImage(req.RepositoryName, req.ImageName, req.Tag)
	if err != nil {
		return "", err
	}

	return id, nil
}
