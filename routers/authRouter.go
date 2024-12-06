package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/email"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/handlers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/middleware"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

// AuthRouter is a router for authentication
func authRouter(r *gin.RouterGroup) {
	// initialize dependencies
	authRepo := repositories.NewAuthRepository(config.DB)
	emailClient := email.NewEmailClient(config.ENV.MAILGUNDOMAIN, config.ENV.MAILGUNKEY, config.ENV.MAILFROM)
	healthRepo := repositories.NewHealthProfileRepository(config.DB)
	authService := services.NewAuthService(authRepo, healthRepo, emailClient)
	authHandler := handlers.NewAuthHandler(authService)

	// auth routes
	prefix := r.Group("/auth")
	prefix.POST("/register", authHandler.Register)
	prefix.GET("/verify/:id", authHandler.VerifyAccount)
	prefix.POST("/login", authHandler.Login)
	prefix.POST("/forgot-password", authHandler.ForgotPassword)
	prefix.GET("/reset-password", authHandler.ShowResetPassword)
	prefix.POST("/reset-password", authHandler.ResetPassword)
	prefix.POST("/refresh-token", authHandler.RefreshToken)
	prefix.POST("/logout", authHandler.Logout)

	// protected routes
	prefix.Use(middleware.AuthMiddleware())
	prefix.POST("/change-password", authHandler.ChangePassword)
}
