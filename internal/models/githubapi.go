package models

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
    Name    string `json:"name"`
    Path    string `json:"path"`
    Sha     string `json:"sha"`
    Size    int    `json:"size"`
    DownloadURL string `json:"download_url"`
}

type Branch struct {
    Name string `json:"name"`
    Commit struct {
        Sha string `json:"sha"`
        Committer struct {
            Name  string `json:"name"`
            Date  string `json:"date"`
        } `json:"committer"`
    } `json:"commit"`
}