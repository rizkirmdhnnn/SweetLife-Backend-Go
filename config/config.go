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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ENV = Env{
		APP_HOST: os.Getenv("APP_HOST"),
		APP_PORT: os.Getenv("APP_PORT"),
		APP_ENV:  os.Getenv("APP_ENV"),
		APP_KEY:  os.Getenv("APP_KEY"),

		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),

		JWTSIGNKEY: os.Getenv("JWTSIGNKEY"),

		MAILGUNKEY:    os.Getenv("MAILGUNKEY"),
		MAILGUNDOMAIN: os.Getenv("MAILGUNDOMAIN"),
		MAILFROM:      os.Getenv("MAILFROM"),

		STORAGE_BUCKET: os.Getenv("STORAGE_BUCKET"),
		STORAGE_FOLDER: os.Getenv("STORAGE_FOLDER"),

		PROJECT_ID:                os.Getenv("PROJECT_ID"),
		GOOGLE_CREDENTIALS_BASE64: os.Getenv("GOOGLE_CREDENTIALS_BASE64"),
	}

	if ENV.APP_ENV == "development" {
		log.Println("Running in development mode")
		ENV.APP_HOST = fmt.Sprintf("%s:%s", ENV.APP_HOST, ENV.APP_PORT)
	} else {
		log.Println("Running in production mode")
	}

	log.Println("Load server successfully")
}
