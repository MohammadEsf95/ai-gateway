package middleware

import (
	"auth/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtUtil utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		userId, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
