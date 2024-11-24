package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

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

	PROJECT_ID string
}

var ENV Env

func LoadEnv() {
	viper.AddConfigPath("./")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal(err.Error())
	}

	if ENV.APP_ENV == "development" {
		log.Println("Running in development mode")
		ENV.APP_HOST = fmt.Sprintf("%s:%s", ENV.APP_HOST, ENV.APP_PORT)
	} else {
		log.Println("Running in production mode")
	}

	log.Println("Load server successfully")
}
