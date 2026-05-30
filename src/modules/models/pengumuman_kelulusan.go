package models

import (
	"time"

	"gorm.io/gorm"
)

// PengumumanKelulusan represents the PengumumanKelulusan model
type PengumumanKelulusan struct {
	ID                         uint           `gorm:"primaryKey" json:"id"`
	SambutanKelulusan          string         `gorm:"type:text;not null" json:"sambutan_kelulusan"`
	TanggalPengumumanNilai     time.Time      `gorm:"type:timestamp;not null" json:"tanggal_pengumuman_nilai"`
	TanggalPengumumanKelulusan time.Time      `gorm:"type:timestamp;not null" json:"tanggal_pengumuman_kelulusan"`
	CreatedAt                  time.Time      `json:"created_at"`
	UpdatedAt                  time.Time      `json:"updated_at"`
	CreatedByID                *uint          `json:"created_by_id"`
	UpdatedByID                *uint          `json:"updated_by_id"`
	DeletedAt                  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign key relationships
	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	UpdatedBy *User `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`
}

// TableName specifies the table name for PengumumanKelulusan
func (m *PengumumanKelulusan) TableName() string {
	return "pengumuman_kelulusan"
}
