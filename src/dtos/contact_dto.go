package dtos

import (
	"time"
)

// JamBukaItem represents opening hours for a specific day
type JamBukaItem struct {
	Hari      string `json:"hari"`
	JamBuka   string `json:"jam_buka"`
	JamTutup  string `json:"jam_tutup"`
}

// ContactCreateRequest represents the request payload for creating Contact
type ContactCreateRequest struct {
	Alamat    string         `json:"alamat" binding:"required"`
	Telepon   string         `json:"telepon" binding:"required"`
	Email     string         `json:"email" binding:"required,email"`
	JamBuka   []JamBukaItem  `json:"jam_buka" binding:"omitempty"`
	Gmaps     string         `json:"gmaps" binding:"omitempty"`
	Website   string         `json:"website" binding:"omitempty"`
	Youtube   string         `json:"youtube" binding:"omitempty"`
	Instagram string         `json:"instagram" binding:"omitempty"`
	Tiktok    string         `json:"tiktok" binding:"omitempty"`
	Facebook  string         `json:"facebook" binding:"omitempty"`
	Twitter   string         `json:"twitter" binding:"omitempty"`
}

// ContactUpdateRequest represents the request payload for updating Contact
type ContactUpdateRequest struct {
	ID        uint           `json:"id" binding:"required"`
	Alamat    string         `json:"alamat" binding:"omitempty"`
	Telepon   string         `json:"telepon" binding:"omitempty"`
	Email     string         `json:"email" binding:"omitempty,email"`
	JamBuka   []JamBukaItem  `json:"jam_buka" binding:"omitempty"`
	Gmaps     string         `json:"gmaps" binding:"omitempty"`
	Website   string         `json:"website" binding:"omitempty"`
	Youtube   string         `json:"youtube" binding:"omitempty"`
	Instagram string         `json:"instagram" binding:"omitempty"`
	Tiktok    string         `json:"tiktok" binding:"omitempty"`
	Facebook  string         `json:"facebook" binding:"omitempty"`
	Twitter   string         `json:"twitter" binding:"omitempty"`
}

// ContactResponse represents the response payload for Contact
type ContactResponse struct {
	ID        uint          `json:"id"`
	Alamat    string        `json:"alamat"`
	Telepon   string        `json:"telepon"`
	Email     string        `json:"email"`
	JamBuka   []JamBukaItem `json:"jam_buka"`
	Gmaps     string        `json:"gmaps"`
	Website   string        `json:"website"`
	Youtube   string        `json:"youtube"`
	Instagram string        `json:"instagram"`
	Tiktok    string        `json:"tiktok"`
	Facebook  string        `json:"facebook"`
	Twitter   string        `json:"twitter"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedByID *uint       `json:"created_by_id"`
	UpdatedByID *uint       `json:"updated_by_id"`
}
