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
	Commit struct {
		Sha       string `json:"sha"`
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
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