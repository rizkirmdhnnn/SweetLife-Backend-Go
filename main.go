package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/routers"
)

func main() {
	config.LoadEnv()
	config.LoadDatabase()
	config.LoadStorageBucket("config/service-account.json")

	router := gin.Default()
	router.LoadHTMLGlob("templates/views/*")

	routers.Routers(router)
	log.Println("Server started on port", config.ENV.APP_PORT)
	router.Run(":3000")
}
