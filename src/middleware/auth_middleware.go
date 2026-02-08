package middleware

import (
	"net/http"
	"strings"

	"pintu-backend/src/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("nama", claims.Nama)
		c.Set("roleID", claims.RoleID)
		c.Set("status", claims.Status)

		c.Next()
	}
}
