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

func recomendationRouter(r *gin.RouterGroup) {
	//initialize dependencies
	httpClient := http.Client{}
	recomendationRepo := repositories.NewRecomendationRepo(&httpClient)
	healthRepo := repositories.NewHealthProfileRepository(config.DB)
	recomendationService := services.NewRecomendationService(recomendationRepo, healthRepo)
	recomendationHandler := handlers.NewRecomendationHandler(recomendationService)
	// user routes

	prefix := r.Group("food-recomendation")
	prefix.Use(middleware.AuthMiddleware())
	prefix.GET("", recomendationHandler.GetRecomendations)
}
