package dtos

import (
	"time"
)

// RoleCreateRequest represents the request payload for creating Role
type RoleCreateRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"omitempty"`
	SystemID       uint   `json:"system_id" binding:"required"`
	Status         string `json:"status" binding:"omitempty"`
	PermissionIDs  []uint `json:"permission_ids" binding:"omitempty"`
}

// RoleUpdateRequest represents the request payload for updating Role
type RoleUpdateRequest struct {
	ID            uint   `json:"id" binding:"required"`
	Name          string `json:"name" binding:"omitempty"`
	Description   string `json:"description" binding:"omitempty"`
	SystemID      *uint  `json:"system_id" binding:"omitempty"`
	Status        string `json:"status" binding:"omitempty"`
	PermissionIDs []uint `json:"permission_ids" binding:"omitempty"`
}

// RoleResponse represents the response payload for Role
type RoleResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	SystemID      *uint            `json:"system_id"`
	System        *SystemResponse  `json:"system"`
	Status        string           `json:"status"`
	Permissions   []PermissionData `json:"permissions"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	CreatedByID   *uint            `json:"created_by_id"`
	UpdatedByID   *uint            `json:"updated_by_id"`
}

// PermissionData represents permission detail in role response
type PermissionData struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupName   string `json:"group_name"`
	System      string `json:"system"`
}

// RoleListResponse represents the response payload for listing Role
type RoleListResponse struct {
	Data  []RoleResponse `json:"data"`
	Total int64          `json:"total"`
}

// RoleGetAllRequest represents the request payload for getting all roles with filters
type RoleGetAllRequest struct {
	Search struct {
		Name     string `json:"name"`
		SystemID uint   `json:"system_id"`
		Status   string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}
