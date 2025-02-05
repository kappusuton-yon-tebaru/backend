package githubapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/githubapi"
)

type Handler struct {
    service *githubapi.Service
}

// NewHandler creates a new Handler instance
func NewHandler(service *githubapi.Service) *Handler {
    return &Handler{
		service,
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