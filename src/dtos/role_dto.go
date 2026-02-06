package dtos

import (
	"time"
)

// RoleCreateRequest represents the request payload for creating Role
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	System      string `json:"system" binding:"required"`
}

// RoleUpdateRequest represents the request payload for updating Role
type RoleUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	System      string `json:"system" binding:"omitempty"`
}

// RoleResponse represents the response payload for Role
type RoleResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	System      string    `json:"system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleListResponse represents the response payload for listing Role
type RoleListResponse struct {
	Data  []RoleResponse `json:"data"`
	Total int64          `json:"total"`
}
