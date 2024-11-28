package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

func userRouter(r *gin.RouterGroup) {
	//initialize dependencies
	userRepo := repositories.NewUserRepository(config.DB)
	authRepo := repositories.NewAuthRepository(config.DB)
	storageRepo := repositories.NewStorageBucketService(config.Client)
	userService := services.NewUserService(userRepo, authRepo, storageRepo)
	storageService := services.NewStorageBucketService(storageRepo)
	userHandler := handlers.NewUserHandler(userService, storageService)

	// user routes
	prefix := r.Group("/users/profile")
	prefix.Use(middleware.AuthMiddleware())
	prefix.GET("/", userHandler.GetProfile)
	prefix.PUT("/", userHandler.UpdateProfile)
}
