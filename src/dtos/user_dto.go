package dtos

import (
	"time"
)

// UserCreateRequest represents the request payload for creating User
type UserCreateRequest struct {
	Nama             string   `json:"nama" binding:"required"`
	Username         string   `json:"username" binding:"required,min=3"`
	Password         string   `json:"password" binding:"required,min=6"`
	RoleID           uint     `json:"role_id" binding:"required"`
	AccessibleSystem []string `json:"accessible_system" binding:"required"`
	Status           string   `json:"status" binding:"required,oneof=active inactive suspended"`
}

// UserUpdateRequest represents the request payload for updating User
type UserUpdateRequest struct {
	Nama             string   `json:"nama" binding:"omitempty"`
	Username         string   `json:"username" binding:"omitempty,min=3"`
	RoleID           uint     `json:"role_id" binding:"omitempty"`
	AccessibleSystem []string `json:"accessible_system" binding:"omitempty"`
	Status           string   `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserUpdatePasswordRequest represents the request payload for updating password
type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserResponse represents the response payload for User
type UserResponse struct {
	ID               uint      `json:"id"`
	Nama             string    `json:"nama"`
	Username         string    `json:"username"`
	RoleID           uint      `json:"role_id"`
	RoleName         string    `json:"role_name"`
	AccessibleSystem []string  `json:"accessible_system"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedByID      *uint     `json:"created_by_id"`
	UpdatedByID      *uint     `json:"updated_by_id"`
}

// UserListResponse represents the response payload for listing User
type UserListResponse struct {
	Data  []UserResponse `json:"data"`
	Total int64          `json:"total"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents the response payload for user login
type UserLoginResponse struct {
	ID               uint      `json:"id"`
	Nama             string    `json:"nama"`
	Username         string    `json:"username"`
	RoleID           uint      `json:"role_id"`
	RoleName         string    `json:"role_name"`
	AccessibleSystem []string  `json:"accessible_system"`
	Status           string    `json:"status"`
	Token            string    `json:"token"`
	CreatedAt        time.Time `json:"created_at"`
}
