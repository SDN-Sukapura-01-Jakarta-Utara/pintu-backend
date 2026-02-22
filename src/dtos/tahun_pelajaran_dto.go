package dtos

import "time"

// TahunPelajaranCreateRequest represents the request payload for creating TahunPelajaran
type TahunPelajaranCreateRequest struct {
	TahunPelajaran string `json:"tahun_pelajaran" binding:"required,min=4,max=20"`
	Status         string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// TahunPelajaranUpdateRequest represents the request payload for updating TahunPelajaran
type TahunPelajaranUpdateRequest struct {
	ID             uint    `json:"id" binding:"required"`
	TahunPelajaran *string `json:"tahun_pelajaran" binding:"omitempty,min=4,max=20"`
	Status         *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// TahunPelajaranResponse represents the response payload for TahunPelajaran
type TahunPelajaranResponse struct {
	ID             uint      `json:"id"`
	TahunPelajaran string    `json:"tahun_pelajaran"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedByID    *uint     `json:"created_by_id"`
	UpdatedByID    *uint     `json:"updated_by_id"`
}

// TahunPelajaranListResponse represents the response payload for listing TahunPelajaran
type TahunPelajaranListResponse struct {
	Data []TahunPelajaranResponse `json:"data"`
}

// TahunPelajaranGetAllRequest represents the request payload for getting all tahun pelajaran with filters
type TahunPelajaranGetAllRequest struct {
	Search struct {
		TahunPelajaran string `json:"tahun_pelajaran"`
		Status         string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// TahunPelajaranListWithPaginationResponse represents the response with pagination
type TahunPelajaranListWithPaginationResponse struct {
	Data       []TahunPelajaranResponse `json:"data"`
	Pagination PaginationInfo           `json:"pagination"`
}
