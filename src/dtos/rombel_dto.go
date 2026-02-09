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

// RombelResponse represents the response payload for Rombel
type RombelResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	KelasID     uint      `json:"kelas_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// RombelListResponse represents the response payload for listing Rombel
type RombelListResponse struct {
	Data []RombelResponse `json:"data"`
}
