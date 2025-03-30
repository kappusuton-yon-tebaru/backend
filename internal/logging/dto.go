package logging

import "time"

type InsertLogDTO struct {
	Timestamp time.Time         `bson:"timestamp"`
	Log       string            `bson:"log"`
	Attribute map[string]string `bson:"attribute"`
}
