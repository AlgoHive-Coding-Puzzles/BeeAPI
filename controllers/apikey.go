package controllers

import (
	"net/http"

	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// CheckApiKey godoc
// @Summary Check API key
// @Description Returns the current API key
// @Tags API Key
// @Produce json
// @Success 200 {object} map[string]string
// CheckApiKey handles the API key check
// @Router /apikey [get]
// @Security Bearer
func CheckApiKey(c * gin.Context) {
	apiKeyManager, err := services.NewAPIKeyManager(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize API key manager"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": apiKeyManager.GetAPIKey()})
}