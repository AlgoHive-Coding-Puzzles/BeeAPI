package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct {}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Ping godoc
// @Summary Health check endpoint
// @Description Returns a pong response to check if the API is alive
// @Tags App
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *HealthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// GetServerName godoc
// @Summary Get server name
// @Description Returns the name of the server
// @Tags App
// @Produce json
// @Success 200 {object} map[string]string
// @Router /name [get]
func (h *HealthController) GetServerName(c *gin.Context) {
	serverName := os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "Local"
	}
	
	serverDescription := os.Getenv("SERVER_DESCRIPTION")
	if serverDescription == "" {
		serverDescription = "Description not set"
	}
	
	c.JSON(http.StatusOK, gin.H{
		"name": serverName,
		"description": serverDescription,
	})
}
