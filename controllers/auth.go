package controllers

import (
	"net/http"
	"strings"

	"github.com/algohive/beeapi/models"
	"github.com/algohive/beeapi/services"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var loginRequest models.LoginRequest
	
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	token, err := a.authService.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, models.LoginResponse{
		Message:  "Login successful",
		Token:    token,
		Username: loginRequest.Username,
	})
}

// Register godoc
// @Summary User registration
// @Description Registers a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var registerRequest models.RegisterRequest
	
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	if len(registerRequest.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 5 characters"})
		return
	}
	
	err := a.authService.RegisterUser(registerRequest.Username, registerRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Logout godoc
// @Summary User logout
// @Description Logs out a user by invalidating their token
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
// @Security Bearer
func (a *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	
	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	err := a.authService.Logout(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Deletes the authenticated user's account
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/delete-account [delete]
// @Security Bearer
func (a *AuthController) DeleteAccount(c *gin.Context) {
	username := c.GetString("username")
	
	err := a.authService.DeleteUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

// GetUser godoc
// @Summary Get user information
// @Description Returns information about the authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Router /auth/user [get]
// @Security Bearer
func (a *AuthController) GetUser(c *gin.Context) {
	username := c.GetString("username")
	
	user, err := a.authService.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}
	
	c.JSON(http.StatusOK, models.UserResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

// CheckAuth godoc
// @Summary Check authentication
// @Description Checks if the provided token is valid
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]bool
// @Failure 401 {object} map[string]string
// @Router /auth/check [get]
// @Security Bearer
func (a *AuthController) CheckAuth(c *gin.Context) {
	// If middleware passed, the token is valid
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}
