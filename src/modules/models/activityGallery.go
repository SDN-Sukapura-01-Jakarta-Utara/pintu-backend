package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ActivityGallery represents the ActivityGallery model
type ActivityGallery struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Judul           string         `gorm:"not null" json:"judul"`
	Tanggal         time.Time      `gorm:"not null" json:"tanggal"`
	Foto            datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"foto"`
	StatusPublikasi string         `gorm:"default:draft" json:"status_publikasi"`
	Status          string         `gorm:"default:active" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CreatedByID     *uint          `json:"created_by_id"`
	UpdatedByID     *uint          `json:"updated_by_id"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for ActivityGallery
func (m *ActivityGallery) TableName() string {
	return "activity_galleries"
}
