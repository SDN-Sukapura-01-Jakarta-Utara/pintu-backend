package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FileItem represents a single file in the files array
type FileItem struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// Article represents the Article model
type Article struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Judul           string         `gorm:"not null" json:"judul"`
	Tanggal         time.Time      `gorm:"not null" json:"tanggal"`
	Kategori        string         `gorm:"not null" json:"kategori"`
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

// TableName specifies the table name for Article
func (m *Article) TableName() string {
	return "articles"
}

// Scan implements the sql.Scanner interface for JSONB
func (f *FileItem) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &f)
}

// Value implements the driver.Valuer interface for JSONB
func (f FileItem) Value() (driver.Value, error) {
	return json.Marshal(f)
}
