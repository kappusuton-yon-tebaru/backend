package githubapi

type GetRepoContentsRequest struct {
	Owner  string `validate:"required"`
	Repo   string `validate:"required"`
	Path   string `validate:"omitempty"`
	Branch string `validate:"omitempty"`
	Token  string `validate:"required"`
}

type GetRepoBranchesRequest struct {
	Owner string `validate:"required"`
	Repo  string `validate:"required"`
	Token string `validate:"required"`
}

type GetCommitMetadataRequest struct {
	Owner  string `validate:"required"`
	Repo   string `validate:"required"`
	Path   string `validate:"required"` // Path is required to fetch commit metadata
	Branch string `validate:"required"` // Branch is required to get correct commit history
	Token  string `validate:"required"`
}

type FetchFileContentRequest struct {
	Owner    string `validate:"required"`
	Repo     string `validate:"required"`
	FilePath string `validate:"required"` // File path is required to fetch content
	Branch   string `validate:"required"` // Branch is required for versioning
	Token    string `validate:"required"`
}

type CreateBranchRequest struct {
	Owner          string `validate:"required"`
	Repo           string `validate:"required"`
	BranchName     string `validate:"required"` // New branch name is required
	SelectedBranch string `validate:"required"` // Base branch to create from is required
	Token          string `validate:"required"`
}

type UpdateFileContentRequest struct {
	Owner         string `validate:"required"`                      // Repository owner
	Repo          string `validate:"required"`                      // Repository name
	Path          string `json:"path" validate:"required"`          // Path to the file
	Message       string `json:"message" validate:"required"`       // Commit message
	Base64Content string `json:"base64Content" validate:"required"` // Base64 encoded file content
	Sha           string `json:"sha" validate:"required"`           // SHA of the existing file
	Branch        string `json:"branch" validate:"required"`        // Target branch
	Token         string `validate:"required"`
}

type GetServicesRequest struct {
	ProjectID string `validate:"required"` // project ID is required
	Token     string `validate:"required"`
}
