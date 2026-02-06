package dtos

import (
	"time"
)

// PermissionCreateRequest represents the request payload for creating Permission
type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	GroupName   string `json:"group_name" binding:"required"`
	System      string `json:"system" binding:"required"`
}

// PermissionUpdateRequest represents the request payload for updating Permission
type PermissionUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	GroupName   string `json:"group_name" binding:"omitempty"`
	System      string `json:"system" binding:"omitempty"`
}

// PermissionResponse represents the response payload for Permission
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupName   string    `json:"group_name"`
	System      string    `json:"system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionListResponse represents the response payload for listing Permission
type PermissionListResponse struct {
	Data   []PermissionResponse `json:"data"`
	Total  int64                `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}
