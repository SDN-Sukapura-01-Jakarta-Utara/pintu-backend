package dtos

import "time"

// VisiMisiCreateRequest represents the request payload for creating VisiMisi
type VisiMisiCreateRequest struct {
	Visi string `json:"visi" binding:"required,min=10"`
	Misi string `json:"misi" binding:"required,min=10"`
}

// VisiMisiUpdateRequest represents the request payload for updating VisiMisi
type VisiMisiUpdateRequest struct {
	ID   uint    `json:"id" binding:"required"`
	Visi *string `json:"visi" binding:"omitempty,min=10"`
	Misi *string `json:"misi" binding:"omitempty,min=10"`
}

// VisiMisiResponse represents the response payload for VisiMisi
type VisiMisiResponse struct {
	ID          uint      `json:"id"`
	Visi        string    `json:"visi"`
	Misi        string    `json:"misi"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// VisiMisiListResponse represents the response payload for listing VisiMisi
type VisiMisiListResponse struct {
	Data   []VisiMisiResponse `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
	Total  int64              `json:"total"`
}
