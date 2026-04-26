package dtos

import (
	"time"
)

// AnnouncementCreateRequest represents the request payload for creating Announcement
type AnnouncementCreateRequest struct {
	Judul           string `json:"judul" binding:"required"`
	Tanggal         string `json:"tanggal" binding:"required"`
	Deskripsi       string `json:"deskripsi" binding:"omitempty"`
	Penulis         string `json:"penulis" binding:"required"`
	StatusPublikasi string `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// AnnouncementUpdateRequest represents the request payload for updating Announcement
type AnnouncementUpdateRequest struct {
	ID              uint     `json:"id" binding:"required"`
	Judul           string   `json:"judul" binding:"omitempty"`
	Tanggal         string   `json:"tanggal" binding:"omitempty"`
	Deskripsi       string   `json:"deskripsi" binding:"omitempty"`
	Penulis         string   `json:"penulis" binding:"omitempty"`
	StatusPublikasi string   `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string   `json:"status" binding:"omitempty,oneof=active inactive"`
	FilesToDelete   []string `json:"files_to_delete" binding:"omitempty"`
}

// AnnouncementResponse represents the response payload for Announcement
type AnnouncementResponse struct {
	ID              uint          `json:"id"`
	Judul           string        `json:"judul"`
	Tanggal         time.Time     `json:"tanggal"`
	Deskripsi       string        `json:"deskripsi"`
	Gambar          string        `json:"gambar"`
	Files           []FileItemDTO `json:"files"`
	Penulis         string        `json:"penulis"`
	StatusPublikasi string        `json:"status_publikasi"`
	Status          string        `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	CreatedByID     *uint         `json:"created_by_id"`
	UpdatedByID     *uint         `json:"updated_by_id"`
}

// AnnouncementListResponse represents the response payload for listing Announcement
type AnnouncementListResponse struct {
	Data   []AnnouncementResponse `json:"data"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
	Total  int64                  `json:"total"`
}

// AnnouncementGetAllRequest represents the request payload for getting all announcements with filters
type AnnouncementGetAllRequest struct {
	Search struct {
		Judul            string `json:"judul"`
		StartDate        string `json:"start_date"` // Format: YYYY-MM-DD
		EndDate          string `json:"end_date"`   // Format: YYYY-MM-DD
		Penulis          string `json:"penulis"`
		StatusPublikasi  string `json:"status_publikasi"`
		Status           string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// AnnouncementListWithPaginationResponse represents the response with pagination
type AnnouncementListWithPaginationResponse struct {
	Data       []AnnouncementResponse `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
}

// AnnouncementPublicResponse represents the public response for announcement
type AnnouncementPublicResponse struct {
	ID        uint      `json:"id"`
	Judul     string    `json:"judul"`
	Tanggal   time.Time `json:"tanggal"`
	Deskripsi string    `json:"deskripsi"`
	Gambar    string    `json:"gambar"`
	Penulis   string    `json:"penulis"`
}

// AnnouncementPublicListResponse represents the public list response
type AnnouncementPublicListResponse struct {
	Data []AnnouncementPublicResponse `json:"data"`
}

// AnnouncementPublicListRequest represents the request payload for public announcement list with filters
type AnnouncementPublicListRequest struct {
	Filter struct {
		Sort string `json:"sort"` // "terbaru" or "terlama"
	} `json:"filter"`
	Offset int `json:"offset"` // Offset for pagination (default 0)
}

// AnnouncementPublicDaftarResponse represents the public list response with pagination
type AnnouncementPublicDaftarResponse struct {
	Data    []AnnouncementPublicResponse `json:"data"`
	Total   int64                        `json:"total"`
	Offset  int                          `json:"offset"`
	HasMore bool                         `json:"has_more"`
}

// AnnouncementPublicDetailResponse represents the public detail response for a single announcement
type AnnouncementPublicDetailResponse struct {
	ID        uint          `json:"id"`
	Judul     string        `json:"judul"`
	Tanggal   time.Time     `json:"tanggal"`
	Deskripsi string        `json:"deskripsi"`
	Gambar    string        `json:"gambar"`
	Penulis   string        `json:"penulis"`
	Files     []FileItemDTO `json:"files"`
}
