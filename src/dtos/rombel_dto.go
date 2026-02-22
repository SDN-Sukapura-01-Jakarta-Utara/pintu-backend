package dtos

import "time"

// RombelCreateRequest represents the request payload for creating Rombel
type RombelCreateRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=20"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
	KelasID uint   `json:"kelas_id" binding:"required"`
}

// RombelUpdateRequest represents the request payload for updating Rombel
type RombelUpdateRequest struct {
	ID      uint    `json:"id" binding:"required"`
	Name    *string `json:"name" binding:"omitempty,min=1,max=20"`
	Status  *string `json:"status" binding:"omitempty,oneof=active inactive"`
	KelasID *uint   `json:"kelas_id" binding:"omitempty"`
}

// KelasDetail represents basic kelas information for RombelResponse
type KelasDetail struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// RombelResponse represents the response payload for Rombel
type RombelResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Status      string      `json:"status"`
	KelasID     uint        `json:"kelas_id"`
	Kelas       KelasDetail `json:"kelas"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CreatedByID *uint       `json:"created_by_id"`
	UpdatedByID *uint       `json:"updated_by_id"`
}

// RombelListResponse represents the response payload for listing Rombel
type RombelListResponse struct {
	Data []RombelResponse `json:"data"`
}

// RombelGetAllRequest represents the request payload for getting all rombel with filters
type RombelGetAllRequest struct {
	Search struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		KelasID uint   `json:"kelas_id"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// RombelListWithPaginationResponse represents the response with pagination
type RombelListWithPaginationResponse struct {
	Data       []RombelResponse `json:"data"`
	Pagination PaginationInfo   `json:"pagination"`
}
