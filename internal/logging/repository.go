package logging

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	log *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		log: db.Collection("log"),
	}
}
