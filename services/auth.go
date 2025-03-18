package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/algohive/beeapi/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidToken     = errors.New("invalid token")
	ErrSessionNotFound  = errors.New("session not found")
)

// JWT secret key
var jwtSecret = []byte("beeapi_secret_key") // In production, use a secure key from environment

// AuthService handles user authentication
type AuthService struct {
	users    map[string]models.User
	sessions map[string]models.Session
	mu       sync.RWMutex
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{
		users:    make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}

// RegisterUser registers a new user
func (a *AuthService) RegisterUser(username, password string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if _, exists := a.users[username]; exists {
		return ErrUserExists
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	a.users[username] = models.User{
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	
	return nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (a *AuthService) AuthenticateUser(username, password string) (string, error) {
	a.mu.RLock()
	user, exists := a.users[username]
	a.mu.RUnlock()
	
	if !exists {
		return "", ErrUserNotFound
	}
	
	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}
	
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24-hour token
	})
	
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	
	// Create session
	sessionID := generateSessionID()
	
	a.mu.Lock()
	a.sessions[sessionID] = models.Session{
		Username:  username,
		CreatedAt: time.Now(),
	}
	a.mu.Unlock()
	
	return tokenString, nil
}

// VerifyToken verifies a JWT token and returns the username
func (a *AuthService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})
	
	if err != nil {
		return "", err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", ErrInvalidToken
		}
		
		a.mu.RLock()
		_, exists := a.users[username]
		a.mu.RUnlock()
		
		if !exists {
			return "", ErrUserNotFound
		}
		
		return username, nil
	}
	
	return "", ErrInvalidToken
}

// DeleteUser deletes a user
func (a *AuthService) DeleteUser(username string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if _, exists := a.users[username]; !exists {
		return ErrUserNotFound
	}
	
	delete(a.users, username)
	
	// Delete any sessions for this user
	for id, session := range a.sessions {
		if session.Username == username {
			delete(a.sessions, id)
		}
	}
	
	return nil
}

// Logout removes a session
func (a *AuthService) Logout(tokenString string) error {
	username, err := a.VerifyToken(tokenString)
	if err != nil {
		return err
	}
	
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Delete any sessions for this user
	for id, session := range a.sessions {
		if session.Username == username {
			delete(a.sessions, id)
		}
	}
	
	return nil
}

// GetUser returns a user
func (a *AuthService) GetUser(username string) (*models.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	user, exists := a.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	return &user, nil
}

// LoadUsersFromEnv loads users from environment variables
func (a *AuthService) LoadUsersFromEnv() {
	// Look for USER_* environment variables
	foundUsers := false
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], "USER_") {
			username := strings.ToLower(strings.TrimPrefix(parts[0], "USER_"))
			password := parts[1]
			a.RegisterUser(username, password)
			foundUsers = true
		}
	}
	
	// If no users were found, create default admin user
	if !foundUsers {
		a.RegisterUser("admin", "admin")
	}
}

// generateSessionID generates a random session ID
func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
