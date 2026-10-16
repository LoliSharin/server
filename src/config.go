package main

import (
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	UploadDir   string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://petya_db:arelun06_db@localhost:5666/Data_cloud"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploading"),
	}
}
