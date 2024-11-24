package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
)

// AuthMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// remove "Bearer " from token
		token = token[7:]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// validate token
		claims, err := helper.ValidateTokenAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// set userID to context
		c.Set("userID", claims.ID)
		c.Next()
	}
}
