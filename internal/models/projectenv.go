package models

import "github.com/kappusuton-yon-tebaru/backend/internal/enum"

type ProjectEnv struct {
	Id        string       `json:"id"`
	ProjectId string       `json:"project_id"`
	EnvType   enum.EnvType `json:"env_type"`
	Key       string       `json:"key"`
	Value     string       `json:"value"`
}
