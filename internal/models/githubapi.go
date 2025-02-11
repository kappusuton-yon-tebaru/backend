package models

import "time"

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
}

type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
}

type Branch struct {
	Name   string `json:"name"`
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
	Commit struct {
		Sha       string `json:"sha"`
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit,omitempty"`  // optional field if needed for other GitHub API responses
}

// CommitMetadata represents metadata for a commit
type CommitMetadata struct {
	LastEditTime  *time.Time `json:"lastEditTime"`
	CommitMessage string     `json:"commitMessage"`
}

// Commit represents the structure of a commit
type Commit struct {
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}
type FileData struct {
	Sha      string `json:"sha"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// BaseBranchResponse represents the response structure when fetching a branch reference.
type BaseBranchResponse struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}

// CreateBranchRequest is the payload for creating a new branch.
type CreateBranchRequest struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type Service struct {
	Name           string `json:"name"`
	DockerfilePath string `json:"dockerfile_path"`
	OwnerRepo 	   string `json:"owner_repo"`
}
// CreateRepoRequest represents the request payload for creating a GitHub repository
type CreateRepoRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
}

// CreateRepoResponse represents the response after creating a repository
type CreateRepoResponse struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}