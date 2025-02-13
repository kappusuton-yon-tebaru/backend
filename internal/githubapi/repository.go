package githubapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type Repository struct{
	cfg *config.Config
}

func NewRepository(cfg *config.Config) *Repository {
    return &Repository{
		cfg ,
	}
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

func (r *Repository) GetRepoContents(ctx context.Context, fullname string, path string, branch string, token string) ([]models.File, error) {
    if token == "" {
        return nil, errors.New("No access token found")
    }

    url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", fullname, path, branch)
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
func (r *Repository) GetRepoBranches(ctx context.Context, fullname string, token string) ([]models.Branch, error) {
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
func (r *Repository) GetCommitMetadata(ctx context.Context,path string, branch string, fullname string, token string) (*models.CommitMetadata, error) {
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

// GetBaseBranchSHA fetches the SHA of the base branch for the given branch name.
func (r *Repository) GetBaseBranchSHA(ctx context.Context, fullname, branchName, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/refs/heads/%s", fullname, branchName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch base branch SHA: %d %s", resp.StatusCode, resp.Status)
	}

	var resData models.BaseBranchResponse
	err = json.NewDecoder(resp.Body).Decode(&resData)
	if err != nil {
		return "", err
	}

	return resData.Object.Sha, nil
}

// CreateBranch creates a new branch from the selected branch.
func (r *Repository) CreateBranch(ctx context.Context, fullname, branchName, baseBranchSHA, token string) (*models.Branch, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/refs", fullname)

	branchRequest := models.CreateBranchRequest{
		Ref: fmt.Sprintf("refs/heads/%s", branchName),
		Sha: baseBranchSHA,
	}

	body, err := json.Marshal(branchRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create branch: %d %s", resp.StatusCode, resp.Status)
	}

	var branchData models.Branch
	err = json.NewDecoder(resp.Body).Decode(&branchData)
	if err != nil {
		return nil, err
	}

	return &branchData, nil
}

// UpdateFileContent updates the file content on GitHub
func (r *Repository) UpdateFileContent(ctx context.Context, fullname, path, commitMsg, base64Content, sha, branch, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", fullname, path)

	body := map[string]interface{}{
		"message": commitMsg,
		"content": base64Content,
		"sha":     sha,
		"branch":  branch,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update file: status code %d", resp.StatusCode)
	}

	return nil
}
// list all files and folders path in a repo on branch main
func (r *Repository) ListFiles(ctx context.Context, fullname, token string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/trees/main?recursive=1",fullname)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API request failed with status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Tree []struct {
			Path string `json:"path"`
		} `json:"tree"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Extract all file paths
	var paths []string
	for _, item := range result.Tree {
		paths = append(paths, item.Path)
	}

	return paths, nil
}
// create repo for the current user
func (r *Repository) CreateRepository(ctx context.Context, token string, repo models.CreateRepoRequest) (*models.CreateRepoResponse, error) {
	url := "https://api.github.com/user/repos"

	payload, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	var response models.CreateRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAccessToken exchanges the authorization code for an access token
func (r *Repository) GetAccessToken(ctx context.Context, code string) (string, error) {
	clientID := r.cfg.ClientID
	clientSecret := r.cfg.ClientSecret
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("GitHub OAuth credentials are missing")
	}

	reqBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("failed to retrieve access token: %v", result)
}