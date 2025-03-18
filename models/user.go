package models

import "time"

// User represents a user in the system
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Don't expose password in JSON
	CreatedAt time.Time `json:"createdAt"`
}

// Session represents an active user session
type Session struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Message  string `json:"message"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

// UserResponse represents user data for API responses
type UserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
