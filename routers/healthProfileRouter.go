package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

func healthRouter(r *gin.RouterGroup) {
	//initialize dependencies
	healthRepo := repositories.NewHealthProfileRepository(config.DB)
	authRepo := repositories.NewAuthRepository(config.DB)
	healthService := services.NewHealthProfileService(healthRepo, authRepo)
	healthHandler := handlers.NewHealthProfileHandler(healthService)

	// user routes
	prefix := r.Group("/users/health")
	prefix.Use(middleware.AuthMiddleware())
	prefix.POST("/", healthHandler.CreateHealthProfile)

	// TODO: update health profile
}
