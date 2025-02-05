package githubapi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/githubapi"
)

// Handler is responsible for handling HTTP requests related to GitHub
type Handler struct {
    service *githubapi.Service
}

// NewHandler creates a new Handler instance
func NewHandler(service *githubapi.Service) *Handler {
	log.Println("Received request to /github/repos") // Log request

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
    name := c.Param("name")
    repo := c.Param("repo")
    path := c.DefaultQuery("path", "") // Default to an empty string if no path is provided

    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
        return
    }

    // Remove "Bearer " prefix from the token
    token = token[len("Bearer "):]
    full_name := name+"/"+repo
    contents, err := h.service.GetRepoContents(full_name, path, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, contents)
}