package dtos

import "time"

// EkstrakurikulerCreateRequest represents the request payload for creating Ekstrakurikuler
type EkstrakurikulerCreateRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	KelasIDs []uint `json:"kelas_ids" binding:"required"`
	Kategori string `json:"kategori" binding:"required,min=2,max=50"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// EkstrakurikulerUpdateRequest represents the request payload for updating Ekstrakurikuler
type EkstrakurikulerUpdateRequest struct {
	ID       uint    `json:"id" binding:"required"`
	Name     *string `json:"name" binding:"omitempty,min=2,max=100"`
	KelasIDs []uint  `json:"kelas_ids" binding:"omitempty"`
	Kategori *string `json:"kategori" binding:"omitempty,min=2,max=50"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// EkstrakurikulerResponse represents the response payload for Ekstrakurikuler
type EkstrakurikulerResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	KelasIDs    []uint    `json:"kelas_ids"`
	Kategori    string    `json:"kategori"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// EkstrakurikulerListResponse represents the response payload for listing Ekstrakurikuler
type EkstrakurikulerListResponse struct {
	Data []EkstrakurikulerResponse `json:"data"`
}
