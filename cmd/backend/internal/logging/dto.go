package logging

import "github.com/kappusuton-yon-tebaru/backend/internal/models"

type GetLogResponse struct {
	Data []models.Log `json:"data"`
}
