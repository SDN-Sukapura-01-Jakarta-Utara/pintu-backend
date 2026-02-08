package dtos

import "time"

// LoginRequest represents the request payload for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response payload for successful login
type LoginResponse struct {
	Token     string    `json:"token"`
	User      UserLoginResponse `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}
