package database

import (
	"os"
)

type MongoConfig struct {
	MongoHost string
	MongoPort string
	MongoDb   string
	MongoUser string
	Password  string
	UserCol   string
}

const (
	UserCol    = "users"
	PodcastCol = "podcasts"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
