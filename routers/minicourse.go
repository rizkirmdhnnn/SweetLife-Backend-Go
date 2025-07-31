package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

func minicourseRouter(r *gin.RouterGroup) {
	//initialize dependencies
	repo := repositories.NewMiniCourseRepository(config.DB)
	service := services.NewMiniCourseService(repo)
	minicourseHandler := handlers.NewMiniCourseHandler(service)

	// user routes
	prefix := r.Group("/minicourse")
	prefix.Use(middleware.AuthMiddleware())
	prefix.GET("/", minicourseHandler.GetMiniCourse)
}
