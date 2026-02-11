package dtos

import (
	"time"
)

// ActivityGalleryCreateRequest represents the request payload for creating ActivityGallery
type ActivityGalleryCreateRequest struct {
	Judul           string `json:"judul" binding:"required"`
	Tanggal         string `json:"tanggal" binding:"required"`
	StatusPublikasi string `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ActivityGalleryUpdateRequest represents the request payload for updating ActivityGallery
type ActivityGalleryUpdateRequest struct {
	ID              uint     `json:"id" binding:"required"`
	Judul           string   `json:"judul" binding:"omitempty"`
	Tanggal         string   `json:"tanggal" binding:"omitempty"`
	StatusPublikasi string   `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string   `json:"status" binding:"omitempty,oneof=active inactive"`
	FotoToDelete    []string `json:"foto_to_delete" binding:"omitempty"`
}

// ActivityGalleryResponse represents the response payload for ActivityGallery
type ActivityGalleryResponse struct {
	ID              uint          `json:"id"`
	Judul           string        `json:"judul"`
	Tanggal         time.Time     `json:"tanggal"`
	Foto            []FileItemDTO `json:"foto"`
	StatusPublikasi string        `json:"status_publikasi"`
	Status          string        `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	CreatedByID     *uint         `json:"created_by_id"`
	UpdatedByID     *uint         `json:"updated_by_id"`
}

// ActivityGalleryListResponse represents the response payload for listing ActivityGallery
type ActivityGalleryListResponse struct {
	Data   []ActivityGalleryResponse `json:"data"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
	Total  int64                     `json:"total"`
}
