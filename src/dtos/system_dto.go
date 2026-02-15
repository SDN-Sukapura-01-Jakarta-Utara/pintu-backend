package dtos

import (
	"time"
)

// SystemCreateRequest represents the request payload for creating System
type SystemCreateRequest struct {
	Nama        string `json:"nama" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty"`
}

// SystemUpdateRequest represents the request payload for updating System
type SystemUpdateRequest struct {
	Nama        string `json:"nama" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty"`
}

// SystemResponse represents the response payload for System
type SystemResponse struct {
	ID          uint      `json:"id"`
	Nama        string    `json:"nama"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// SystemGetAllRequest represents the request payload for getting all systems with filters
type SystemGetAllRequest struct {
	Search struct {
		Nama   string `json:"nama"`
		Status string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}
