package dtos

import (
	"time"
)

// PermissionCreateRequest represents the request payload for creating Permission
type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	GroupName   string `json:"group_name" binding:"required"`
	SystemID    uint   `json:"system_id" binding:"required"`
	Status      string `json:"status" binding:"omitempty"`
}

// PermissionUpdateRequest represents the request payload for updating Permission
type PermissionUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	GroupName   string `json:"group_name" binding:"omitempty"`
	SystemID    *uint  `json:"system_id" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty"`
}

// PermissionResponse represents the response payload for Permission
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupName   string    `json:"group_name"`
	SystemID    *uint     `json:"system_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// PermissionListResponse represents the response payload for listing Permission
type PermissionListResponse struct {
	Data   []PermissionResponse `json:"data"`
	Total  int64                `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}
