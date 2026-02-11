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
