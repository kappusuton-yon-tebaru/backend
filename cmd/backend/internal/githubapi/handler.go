package githubapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/githubapi"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/projectrepository"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	cfg                *config.Config
	service            *githubapi.Service
	projectRepoService *projectrepository.Service
	validator          *validator.Validator
}

// NewHandler creates a new Handler instance
func NewHandler(cfg *config.Config, service *githubapi.Service, projectRepoService *projectrepository.Service, validator *validator.Validator) *Handler {
	return &Handler{
		cfg,
		service,
		projectRepoService,
		validator,
	}
}

// fetch GitHub user repositories
func (h *Handler) GetUserRepos(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token, found := strings.CutPrefix(token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	ctx := c.Request.Context() // Get context from the Gin request
	repos, err := h.service.GetUserRepos(ctx, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repos)
}

// Get repo files and folders, add ?path={foldername} to get folder contents  (can be use for updating file tree)
func (h *Handler) GetRepoContents(c *gin.Context) {
	req := GetRepoContentsRequest{
		Owner:  c.Param("owner"),
		Repo:   c.Param("repo"),
		Path:   c.DefaultQuery("path", ""),
		Branch: c.DefaultQuery("branch", ""),
		Token:  c.GetHeader("Authorization"),
	}
	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	contents, err := h.service.GetRepoContents(c.Request.Context(), req.Owner+"/"+req.Repo, req.Path, req.Branch, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contents)
}

// get all branch in that repo
func (h *Handler) GetRepoBranches(c *gin.Context) {
	req := GetRepoBranchesRequest{
		Owner: c.Param("owner"),
		Repo:  c.Param("repo"),
		Token: c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	fullname := req.Owner + "/" + req.Repo
	branches, err := h.service.GetRepoBranches(c.Request.Context(), fullname, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// Get "lastEditTime" and  "commitMessage" of a file in ... repo on ... branch
func (h *Handler) GetCommitMetadata(c *gin.Context) {
	req := GetCommitMetadataRequest{
		Owner:  c.Param("owner"),
		Repo:   c.Param("repo"),
		Path:   c.DefaultQuery("path", ""),
		Branch: c.DefaultQuery("branch", ""),
		Token:  c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	fullname := req.Owner + "/" + req.Repo
	metadata, err := h.service.GetCommitMetadata(c.Request.Context(), req.Path, req.Branch, fullname, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

// Get content in that file return "content", "sha"
func (h *Handler) FetchFileContent(c *gin.Context) {
	req := FetchFileContentRequest{
		Owner:    c.Param("owner"),
		Repo:     c.Param("repo"),
		FilePath: c.Query("path"),
		Branch:   c.Query("branch"),
		Token:    c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	fullname := req.Owner + "/" + req.Repo
	fileContent, err := h.service.FetchFileContent(c.Request.Context(), fullname, req.FilePath, req.Branch, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fileContent)
}

// CreateBranch handles requests to create a new branch.
func (h *Handler) CreateBranch(c *gin.Context) {
	req := CreateBranchRequest{
		Owner:          c.Param("owner"),
		Repo:           c.Param("repo"),
		BranchName:     c.DefaultQuery("branch_name", ""),
		SelectedBranch: c.DefaultQuery("selected_branch", ""),
		Token:          c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	fullname := req.Owner + "/" + req.Repo

	// Step 1: Get the base branch SHA (the selected branch)
	baseBranchSHA, err := h.service.GetBaseBranchSHA(c.Request.Context(), fullname, req.SelectedBranch, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Create the new branch
	branch, err := h.service.CreateBranch(c.Request.Context(), fullname, req.BranchName, baseBranchSHA, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branch)
}

// push commit
func (h *Handler) UpdateFileContent(c *gin.Context) {
	req := UpdateFileContentRequest{
		Owner: c.Param("owner"),
		Repo:  c.Param("repo"),
		Token: c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	// Bind JSON payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	fullname := req.Owner + "/" + req.Repo

	// Call service to update file content
	err := h.service.UpdateFileContent(c.Request.Context(), fullname, req.Path, req.Message, req.Base64Content, req.Sha, req.Branch, req.Token)
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
	// Initialize request struct
	req := GetServicesRequest{
		ProjectID: c.Param("id"),
		Token:     c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix
	token, found := strings.CutPrefix(req.Token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	req.Token = token

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	// Fetch project repository information by projectID
	projRepo, err := h.projectRepoService.GetProjectRepositoryByProjectId(c, req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract the full repo name from the Git repo URL
	fullname, err2 := ExtractFullName(projRepo.GitRepoUrl)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	// Retrieve the services from the repository
	services, err3 := h.service.FindServices(c, fullname, token)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	// Return the services as a JSON response
	c.JSON(http.StatusOK, services)
}

func (h *Handler) CreateRepository(c *gin.Context) {
	var req models.CreateRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"messages": h.validator.Translate(err),
		})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Remove "Bearer " prefix from the token
	token, found := strings.CutPrefix(token, "Bearer ")
	if !found || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	repo, err := h.service.CreateRepository(c.Request.Context(), token, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repo)
}

// Redirects user to GitHub OAuth authorization page
func (h *Handler) RedirectToGitHub(c *gin.Context) {

	if h.cfg.ClientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub Client ID missing"})
		return
	}

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=repo", h.cfg.ClientID)
	c.Redirect(http.StatusFound, redirectURL)
}

// Handles the callback from GitHub, exchanges code for token, and stores it
func (h *Handler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	err := h.service.AuthenticateUser(context.Background(), code, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}
