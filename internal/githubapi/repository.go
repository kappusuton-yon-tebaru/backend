package githubapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
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

// GetCommitMetadata fetches the commit metadata for a file in a specific branch
func (r *Repository) GetCommitMetadata(path string, branch string, fullname string, token string) (*models.CommitMetadata, error) {
    if token == "" {
        return nil, errors.New("No access token found")
    }

    url := fmt.Sprintf(
        "https://api.github.com/repos/%s/commits?path=%s&per_page=1&sha=%s&%d",
        fullname,
        path,
        branch,
        time.Now().Unix(),
    )

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
        return nil, fmt.Errorf("Failed to fetch commit metadata, status code: %d", resp.StatusCode)
    }

    var commitData []models.Commit
    err = json.NewDecoder(resp.Body).Decode(&commitData)
    if err != nil {
        return nil, err
    }

    // Extract metadata from the first commit (latest one)
    if len(commitData) > 0 {
        commit := commitData[0].Commit

        // Parse the date string
        parsedTime, err := time.Parse(time.RFC3339, commit.Author.Date)
        if err != nil {
            return nil, fmt.Errorf("Failed to parse commit date: %v", err)
        }

        return &models.CommitMetadata{
            LastEditTime: &parsedTime,
            CommitMessage: commit.Message,
        }, nil
    }

    return &models.CommitMetadata{
        LastEditTime: nil,
        CommitMessage: "No commits found",
    }, nil
}

func (r *Repository) FetchFileContent(ctx context.Context, fullname, filePath, branch, token string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", fullname, filePath, branch)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error fetching file content: %d %s", resp.StatusCode, resp.Status)
	}

	var fileData models.FileData
	err = json.NewDecoder(resp.Body).Decode(&fileData)
	if err != nil {
		return "", "", err
	}

	if fileData.Encoding == "base64" {
		decodedContent, err := base64.StdEncoding.DecodeString(fileData.Content)
		if err != nil {
			return "", "", errors.New("failed to decode base64 content")
		}
		return string(decodedContent), fileData.Sha, nil
	}

	return "", "", errors.New("unsupported file encoding")
}