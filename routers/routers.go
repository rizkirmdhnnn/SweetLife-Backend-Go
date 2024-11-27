package routers

import (
	"github.com/gin-gonic/gin"
)

// Routers is a function to define all the routes
func Routers(r *gin.Engine) {
	prefix := r.Group("/api/v1/")
	authRouter(prefix)
	userRouter(prefix)
	healthRouter(prefix)
	recomendationRouter(prefix)

	// dummy route
	prefix.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Welcome to SweetLife Dashboard",
			"data": gin.H{
				"user": gin.H{
					"name":         "Jokowi",
					"diabetesType": "TYPE_2",
				},
				"dailyProgress": gin.H{
					"calorie": gin.H{
						"current":      250,
						"target":       1500,
						"percentage":   16.67,
						"satisfaction": "UNDER",
					},
					"glucose": gin.H{
						"current":      250,
						"target":       1500,
						"percentage":   16.67,
						"satisfaction": "UNDER",
					},
				},
				"status": gin.H{
					"satisfaction": "UNDER", // UNDER (<50%), GOOD (50-100%), EXCEED (>100%)
					"message":      "Ayo tingkatkan asupan kalori harianmu!",
				},
			},
		})
	})

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}
