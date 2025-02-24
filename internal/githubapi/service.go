package githubapi

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetUserRepos(ctx context.Context, token string) (map[string]interface{}, error) {
	repos, err := s.repo.GetUserRepos(ctx, token)
	if err != nil {
		return nil,err
	}
	response := map[string]interface{}{
		"data": repos,
	}

	return response,nil
}

func (s *Service) GetRepoContents(ctx context.Context, fullname string, path string, branch string, token string) (map[string]interface{}, error) {
	contents, err := s.repo.GetRepoContents(ctx, fullname, path, branch, token)
	if err != nil {
		return nil,err
	}
	response := map[string]interface{}{
		"data": contents,
	}

	return response,nil
}

// GetRepoBranches fetches the branches of a repository
func (s *Service) GetRepoBranches(ctx context.Context, fullname string, token string) (map[string]interface{}, error) {
	if fullname == "" {
		return nil, errors.New("Repository fullname is required")
	}

	branches, err := s.repo.GetRepoBranches(ctx, fullname, token)

	if err != nil {
		return nil,err
	}
	response := map[string]interface{}{
		"data": branches,
	}

	return response,nil
}

// GetCommitMetadata fetches the commit metadata for a file in a repository
func (s *Service) GetCommitMetadata(ctx context.Context, path string, branch string, fullname string, token string) (map[string]interface{}, error) {
	if fullname == "" || path == "" || branch == "" {
		return nil, errors.New("Repository fullname, path, and branch are required")
	}

	commitMetadata, err := s.repo.GetCommitMetadata(ctx, path, branch, fullname, token)

	if err != nil {
		return nil,err
	}
	response := map[string]interface{}{
		"data": commitMetadata,
	}

	return response,nil
}

func (s *Service) FetchFileContent(ctx context.Context, fullname, filePath, branch, token string) (map[string]interface{}, error) {
	if fullname == "" || filePath == "" || branch == "" || token == "" {
		return nil, errors.New("missing required parameters")
	}
	content, sha,err := s.repo.FetchFileContent(ctx, fullname, filePath, branch, token)
	
	if err != nil {
		return nil,err
	}

	response := map[string]interface{}{
		"data":map[string]interface{}{
			"sha":     sha,
			"content": content,
		},
	}

	return response,nil 
}

// GetBaseBranchSHA calls the repository to get the SHA of the base branch.
func (s *Service) GetBaseBranchSHA(ctx context.Context, fullname, branchName, token string) (string, error) {
	if fullname == "" || branchName == "" || token == "" {
		return "", errors.New("missing required parameters")
	}

	return s.repo.GetBaseBranchSHA(ctx, fullname, branchName, token)
}

// CreateBranch calls the repository to create a new branch.
func (s *Service) CreateBranch(ctx context.Context, fullname, branchName, baseBranchSHA, token string) (*models.Branch, error) {
	if fullname == "" || branchName == "" || baseBranchSHA == "" || token == "" {
		return nil, errors.New("missing required parameters")
	}

	return s.repo.CreateBranch(ctx, fullname, branchName, baseBranchSHA, token)
}

func (s *Service) UpdateFileContent(ctx context.Context, fullname, path, commitMsg, base64Content, sha, branch, token string) error {
	return s.repo.UpdateFileContent(ctx, fullname, path, commitMsg, base64Content, sha, branch, token)
}

func (s *Service) FindServices(ctx context.Context, fullname, token string) (map[string]interface{}, error) {
	files, err := s.repo.ListFiles(ctx, fullname, token)
	if err != nil {
		return nil, err
	}
	var services []models.Service
	seen := make(map[string]bool)

	for _, filePath := range files {

		if strings.HasPrefix(filePath, "apps/") && strings.HasSuffix(strings.ToLower(filePath), "/dockerfile") {
			parts := strings.Split(filePath, "/")
			if len(parts) >= 2 {
				serviceName := parts[1]

				// Ensure no duplicates
				if _, exists := seen[serviceName]; !exists {
					services = append(services, models.Service{
						Name:           serviceName,
						DockerfilePath: filePath,
					})
					seen[serviceName] = true
				}
			}
		}
	}

	response := map[string]interface{}{
		"data": services,
	}

	return response, nil
}

func (s *Service) CreateRepository(ctx context.Context, token string, repo models.CreateRepoRequest) (*models.CreateRepoResponse, error) {
	return s.repo.CreateRepository(ctx, token, repo)
}

// AuthenticateUser retrieves the GitHub access token and stores it
func (s *Service) AuthenticateUser(ctx context.Context, code string, c *gin.Context) error {
	token, err := s.repo.GetAccessToken(ctx, code)
	if err != nil {
		return err
	}

	// Store token securely in an HTTP-only cookie
	c.SetCookie("github_token", token, 3600, "/", "", false, true)

	return nil
}
