package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ENV Env

type Env struct {
	APP_HOST string
	APP_PORT string
	APP_ENV  string
	APP_KEY  string

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	JWTSIGNKEY string

	MAILGUNKEY    string
	MAILGUNDOMAIN string
	MAILFROM      string

	STORAGE_BUCKET string
	STORAGE_FOLDER string

	PROJECT_ID                string
	GOOGLE_CREDENTIALS_BASE64 string
}

func LoadEnv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load environment variables
	ENV = Env{
		APP_HOST: getEnv("APP_HOST", "127.0.0.1"),
		APP_PORT: getEnv("APP_PORT", "3000"),
		APP_ENV:  getEnv("APP_ENV", "development"),
		APP_KEY:  getEnv("APP_KEY", "anakepakyanto"),

		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_NAME:     getEnv("DB_NAME", "sweetlife"),

		JWTSIGNKEY: getEnv("JWTSIGNKEY", "anakepakyanto"),

		MAILGUNKEY:    getEnv("MAILGUNKEY", ""),
		MAILGUNDOMAIN: getEnv("MAILGUNDOMAIN", ""),
		MAILFROM:      getEnv("MAILFROM", ""),

		STORAGE_BUCKET: getEnv("STORAGE_BUCKET", ""),
		STORAGE_FOLDER: getEnv("STORAGE_FOLDER", ""),

		PROJECT_ID:                getEnv("PROJECT_ID", ""),
		GOOGLE_CREDENTIALS_BASE64: getEnv("GOOGLE_CREDENTIALS_BASE64", ""),
	}

	if ENV.APP_ENV == "development" {
		log.Println("Running in development mode")
		ENV.APP_HOST = fmt.Sprintf("%s:%s", ENV.APP_HOST, ENV.APP_PORT)
	} else {
		log.Println("Running in production mode")
	}

	log.Println("Environment loaded successfully")
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
