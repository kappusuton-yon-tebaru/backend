package githubapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/githubapi"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
)

type Handler struct {
    service *githubapi.Service
	projectRepoService *projectrepository.Service
}

// NewHandler creates a new Handler instance
func NewHandler(service *githubapi.Service, projectRepoService *projectrepository.Service) *Handler {
    return &Handler{
		service,
		projectRepoService,
	}
}

// GetUserRepos handles the HTTP request to fetch GitHub user repositories
func (h *Handler) GetUserRepos(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]

    ctx := c.Request.Context() // Get context from the Gin request
    repos, err := h.service.GetUserRepos(ctx, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, repos)
}

func (h *Handler) GetRepoContents(c *gin.Context) {
    fullname := c.Param("owner")+"/"+c.Param("repo")
    path := c.DefaultQuery("path", "") // Default to an empty string if no path is provided

    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]
    contents, err := h.service.GetRepoContents(c.Request.Context(),fullname, path, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, contents)
}

// GetRepoBranches handles requests to fetch branches of a GitHub repository
func (h *Handler) GetRepoBranches(c *gin.Context) {
    fullname := c.Param("owner")+"/"+c.Param("repo")

    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]

    branches, err := h.service.GetRepoBranches(c.Request.Context(), fullname, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, branches)
}

// GetCommitMetadata handles requests to fetch commit metadata for a file in a GitHub repository
func (h *Handler) GetCommitMetadata(c *gin.Context) {
    fullname := c.Param("owner")+"/"+c.Param("repo")

    path := c.DefaultQuery("path", "") // Get the path query parameter (e.g., "README.md")
    branch := c.DefaultQuery("branch", "") // Get the branch query parameter (e.g., "main")

    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]

    metadata, err := h.service.GetCommitMetadata(c.Request.Context(), path, branch, fullname, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, metadata)
}

func (h *Handler) FetchFileContent(c *gin.Context) {
	fullname := c.Param("owner") + "/" + c.Param("repo")
	filePath := c.Query("path")
	branch := c.Query("branch")
	token := c.GetHeader("Authorization")

    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]

	content, sha, err := h.service.FetchFileContent(c.Request.Context(), fullname, filePath, branch, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sha":     sha,
		"content": content,
	})
}

// CreateBranch handles requests to create a new branch.
func (h *Handler) CreateBranch(c *gin.Context) {
	fullname := c.Param("owner") + "/" + c.Param("repo")
	branchName := c.DefaultQuery("branch_name", "")
	selectedBranchName := c.DefaultQuery("selected_branch", "")
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token = token[len("Bearer "):]

	// Step 1: Get the base branch SHA (the selected branch)
	baseBranchSHA, err := h.service.GetBaseBranchSHA(c.Request.Context(), fullname, selectedBranchName, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Create the new branch
	branch, err := h.service.CreateBranch(c.Request.Context(), fullname, branchName, baseBranchSHA, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branch)
}

// UpdateFileContent handles the file content update route using JSON payload
func (h *Handler) UpdateFileContent(c *gin.Context) {
	fullname := c.Param("owner") + "/" + c.Param("repo")
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token = token[len("Bearer "):]

	// Bind JSON payload to struct
	var req struct {
		Path      string `json:"path"`
		Message   string `json:"message"`
		Base64Content   string `json:"base64Content"`
		Sha       string `json:"sha"`
		Branch    string `json:"branch"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call service to update file content
	err := h.service.UpdateFileContent(c.Request.Context(), fullname, req.Path, req.Message, req.Base64Content, req.Sha, req.Branch, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File successfully updated"})
}

func ExtractFullName(gitURL string) (string, error) {
	parsedURL, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}

	// Ensure it's a GitHub URL
	if parsedURL.Host != "github.com" {
		return "", fmt.Errorf("invalid GitHub URL")
	}

	// Trim leading and trailing slashes
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	// Ensure we have exactly "owner/repo"
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid repository URL format")
	}

	return fmt.Sprintf("%s/%s", parts[0], parts[1]), nil
}

func (h *Handler) GetServices(c *gin.Context) {

	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}
	// Remove "Bearer " prefix from the token
	token = token[len("Bearer "):] 

	projRepos, err := h.projectRepoService.GetProjectRepositoriesByProjectID(c, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	var all_services []models.Service
	for _, projRepos := range projRepos{
		fullname, err := ExtractFullName(projRepos.GitRepoUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}

		services, err := h.service.FindServices(c,fullname,token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
			return
		}

		all_services = append(all_services, services...)
	}

	c.JSON(http.StatusOK, all_services)
}
