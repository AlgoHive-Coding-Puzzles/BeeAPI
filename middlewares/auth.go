package middlewares

import (
	"net/http"
	"strings"

	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// RequireAPIKey crée un middleware qui valide la clé API
func RequireAPIKey(keyManager *services.APIKeyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère l'en-tête Authorization
		authHeader := c.GetHeader("Authorization")

		// Vérifie si l'en-tête Authorization est présent et a le format correct
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Clé API manquante"})
			return
		}

		// Le format doit être "Bearer <api-key>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format de clé API invalide"})
			return
		}

		key := parts[1]

		// Valide la clé API
		if !keyManager.ValidateKey(key) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Clé API invalide"})
			return
		}

		c.Next()
	}
}
