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