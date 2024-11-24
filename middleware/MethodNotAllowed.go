package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() == http.StatusNotFound {
			if c.Request.Method != "POST" {
				c.JSON(http.StatusMethodNotAllowed, gin.H{
					"error": "Method Not Allowed",
				})
				return
			}
		}
	}
}
