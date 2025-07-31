package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

func miniGroceryRouter(r *gin.RouterGroup) {
	//initialize dependencies
	repo := repositories.NewMiniGroceryRepository(config.DB)
	service := services.NewMiniGroceryService(repo)
	miniGroceryHandler := handlers.NewMiniGroceryHandler(service)

	// user routes
	prefix := r.Group("/minigrocery")
	prefix.Use(middleware.AuthMiddleware())
	prefix.GET("/", miniGroceryHandler.GetMiniGrocery)
}
