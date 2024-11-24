package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/routers"
)

func main() {
	// Load environment variables, database, and storage
	config.LoadEnv()
	config.LoadDatabase()
	config.LoadStorageBucket()

	// Initialize Gin router
	router := gin.Default()
	router.LoadHTMLGlob("templates/views/*")

	// Add the middleware for handling 405 Method Not Allowed globally
	router.Use(middleware.MethodNotAllowedMiddleware())

	// Set up routes
	routers.Routers(router)

	// Log and start the server
	log.Println("Server started on port", config.ENV.APP_PORT)
	router.Run(":" + config.ENV.APP_PORT)
}
