package dtos

import (
	"time"
)

// FileItemDTO represents a single file in the files array
type FileItemDTO struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// ArticleCreateRequest represents the request payload for creating Article
type ArticleCreateRequest struct {
	Judul           string    `json:"judul" binding:"required"`
	Tanggal         string    `json:"tanggal" binding:"required"`
	Kategori        string    `json:"kategori" binding:"required"`
	Deskripsi       string    `json:"deskripsi" binding:"omitempty"`
	Penulis         string    `json:"penulis" binding:"required"`
	StatusPublikasi string    `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string    `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ArticleUpdateRequest represents the request payload for updating Article
type ArticleUpdateRequest struct {
	ID              uint      `json:"id" binding:"required"`
	Judul           string    `json:"judul" binding:"omitempty"`
	Tanggal         string    `json:"tanggal" binding:"omitempty"`
	Kategori        string    `json:"kategori" binding:"omitempty"`
	Deskripsi       string    `json:"deskripsi" binding:"omitempty"`
	Penulis         string    `json:"penulis" binding:"omitempty"`
	StatusPublikasi string    `json:"status_publikasi" binding:"omitempty,oneof=draft published archived"`
	Status          string    `json:"status" binding:"omitempty,oneof=active inactive"`
	FilesToDelete   []string  `json:"files_to_delete" binding:"omitempty"`
}

// ArticleResponse represents the response payload for Article
type ArticleResponse struct {
	ID              uint          `json:"id"`
	Judul           string        `json:"judul"`
	Tanggal         time.Time     `json:"tanggal"`
	Kategori        string        `json:"kategori"`
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

// ArticleListResponse represents the response payload for listing Article
type ArticleListResponse struct {
	Data   []ArticleResponse `json:"data"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Total  int64             `json:"total"`
}
