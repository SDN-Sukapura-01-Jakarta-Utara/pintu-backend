package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Announcement represents the Announcement model
type Announcement struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Judul           string         `gorm:"not null" json:"judul"`
	Tanggal         time.Time      `gorm:"not null" json:"tanggal"`
	Deskripsi       string         `gorm:"type:text" json:"deskripsi"`
	Gambar          string         `json:"gambar"`
	Files           datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"files"`
	Penulis         string         `gorm:"not null" json:"penulis"`
	StatusPublikasi string         `gorm:"default:draft" json:"status_publikasi"`
	Status          string         `gorm:"default:active" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CreatedByID     *uint          `json:"created_by_id"`
	UpdatedByID     *uint          `json:"updated_by_id"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Announcement
func (m *Announcement) TableName() string {
	return "announcements"
}
