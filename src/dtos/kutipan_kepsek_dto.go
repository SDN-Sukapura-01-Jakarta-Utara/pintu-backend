package dtos

import "time"

// KutipanKepsekCreateRequest represents the request payload for creating KutipanKepsek
type KutipanKepsekCreateRequest struct {
	NamaKepsek    string `json:"nama_kepsek" binding:"required,min=2,max=255"`
	KutipanKepsek string `json:"kutipan_kepsek" binding:"required,min=5"`
}

// KutipanKepsekUpdateRequest represents the request payload for updating KutipanKepsek
type KutipanKepsekUpdateRequest struct {
	ID            uint    `json:"id" binding:"required"`
	NamaKepsek    *string `json:"nama_kepsek" binding:"omitempty,min=2,max=255"`
	KutipanKepsek *string `json:"kutipan_kepsek" binding:"omitempty,min=5"`
}

// KutipanKepsekResponse represents the response payload for KutipanKepsek
type KutipanKepsekResponse struct {
	ID            uint      `json:"id"`
	NamaKepsek    string    `json:"nama_kepsek"`
	FotoKepsek    string    `json:"foto_kepsek"`
	KutipanKepsek string    `json:"kutipan_kepsek"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedByID   *uint     `json:"created_by_id"`
	UpdatedByID   *uint     `json:"updated_by_id"`
}

// KutipanKepsekListResponse represents the response payload for listing KutipanKepsek
type KutipanKepsekListResponse struct {
	Data   []KutipanKepsekResponse `json:"data"`
	Limit  int                     `json:"limit"`
	Offset int                     `json:"offset"`
	Total  int64                   `json:"total"`
}
