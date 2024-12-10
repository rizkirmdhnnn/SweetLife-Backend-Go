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
	healthRepo := repositories.NewHealthProfileRepository(config.DB)
	userService := services.NewUserService(userRepo, authRepo, storageRepo, healthRepo)
	storageService := services.NewStorageBucketService(storageRepo)
	userHandler := handlers.NewUserHandler(userService, storageService)

	// user routes
	prefix := r.Group("/users")
	prefix.Use(middleware.AuthMiddleware())
	prefix.GET("/profile", userHandler.GetProfile)
	prefix.PUT("/profile", userHandler.UpdateProfile)
	prefix.GET("/history", userHandler.GetHistory)
	prefix.GET("/dashboard", userHandler.GetDashboard)
}
