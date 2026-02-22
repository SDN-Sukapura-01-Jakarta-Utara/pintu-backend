package dtos

import "time"

// KelasCreateRequest represents the request payload for creating Kelas
type KelasCreateRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// KelasUpdateRequest represents the request payload for updating Kelas
type KelasUpdateRequest struct {
	ID     uint    `json:"id" binding:"required"`
	Name   *string `json:"name" binding:"omitempty,min=1,max=100"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// KelasResponse represents the response payload for Kelas
type KelasResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// KelasListResponse represents the response payload for listing Kelas
type KelasListResponse struct {
	Data []KelasResponse `json:"data"`
}

// KelasGetAllRequest represents the request payload for getting all kelas with filters
type KelasGetAllRequest struct {
	Search struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// KelasListWithPaginationResponse represents the response with pagination
type KelasListWithPaginationResponse struct {
	Data       []KelasResponse `json:"data"`
	Pagination PaginationInfo  `json:"pagination"`
}
