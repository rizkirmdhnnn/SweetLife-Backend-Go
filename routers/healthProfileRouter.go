package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

func healthRouter(r *gin.RouterGroup) {
	//initialize dependencies
	httpClient := http.Client{}
	healthRepo := repositories.NewHealthProfileRepository(config.DB)
	authRepo := repositories.NewAuthRepository(config.DB)
	recomendRepo := repositories.NewRecomendationRepo(&httpClient)
	healthService := services.NewHealthProfileService(healthRepo, authRepo, recomendRepo)
	healthHandler := handlers.NewHealthProfileHandler(healthService)

	// user routes
	prefix := r.Group("/users/health")
	prefix.Use(middleware.AuthMiddleware())
	prefix.POST("/", healthHandler.CreateHealthProfile)
	prefix.GET("/", healthHandler.GetHealthProfile)

	// TODO: update health profile
}
