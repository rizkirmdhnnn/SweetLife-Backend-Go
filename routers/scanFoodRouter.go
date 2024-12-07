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

func scanFoodRouter(r *gin.RouterGroup) {
	//initialize dependencies
	repo := repositories.NewScanFoodRepository(&http.Client{}, config.DB, config.ENV.USDA_API_KEY)
	storageRepo := repositories.NewStorageBucketService(config.Client)
	service := services.NewScanFoodService(repo, storageRepo)
	scanFoodhandler := handlers.NewScanFoodHandler(service)

	// user routes
	prefix := r.Group("/food")
	prefix.Use(middleware.AuthMiddleware())
	prefix.POST("/scan", scanFoodhandler.ScanFood)
	prefix.POST("/find", scanFoodhandler.FindFood)
}
