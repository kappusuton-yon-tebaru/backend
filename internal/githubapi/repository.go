package githubapi

import (
	"context"
    "encoding/json"
    "fmt"
    "net/http"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
    "errors"
)

type Repository struct{}

func NewRepository() *Repository {
    return &Repository{}
}

func (r *Repository) GetUserRepos(ctx context.Context, token string) ([]models.Repository, error) {
    if token == "" {
        return nil, errors.New("No access token found")
    }

    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/repos", nil)
    if err != nil {
        return nil, err
    }

    // Set Authorization header with the Bearer token
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Failed to fetch repositories, status code: %d", resp.StatusCode)
    }

    var repos []models.Repository
    err = json.NewDecoder(resp.Body).Decode(&repos)
    if err != nil {
        return nil, err
    }

    return repos, nil
}

func (r *Repository) GetRepoContents(fullname string, path string, token string) ([]models.File, error) {
    if token == "" {
        return nil, errors.New("No access token found")
    }

    url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", fullname, path)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Set Authorization header with the Bearer token
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Failed to fetch repository contents, status code: %d", resp.StatusCode)
    }

    var contents []models.File
    err = json.NewDecoder(resp.Body).Decode(&contents)
    if err != nil {
        return nil, err
    }

    return contents, nil
}

// GetRepoBranches fetches the branches of a GitHub repository
func (r *Repository) GetRepoBranches(fullname string, token string) ([]models.Branch, error) {
    if token == "" {
        return nil, errors.New("No access token found")
    }

    url := fmt.Sprintf("https://api.github.com/repos/%s/branches", fullname)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Set Authorization header with the Bearer token
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Failed to fetch branches, status code: %d", resp.StatusCode)
    }

    var branches []models.Branch
    err = json.NewDecoder(resp.Body).Decode(&branches)
    if err != nil {
        return nil, err
    }

    return branches, nil
}