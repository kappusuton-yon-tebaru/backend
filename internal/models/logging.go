package models

import "time"

type Log struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Log       string    `json:"log"`
}
