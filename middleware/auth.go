package middleware

import (
	"net/http"
	"strings"

	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware to protect routes that require authentication
func RequireAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		// Check if the header format is correct
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		
		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Verify token
		username, err := authService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
		// Store username in context for use in handlers
		c.Set("username", username)
		
		c.Next()
	}
}
