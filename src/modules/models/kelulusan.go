package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Kelulusan represents the Kelulusan model
type Kelulusan struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NomorPeserta  string         `gorm:"size:50;not null;unique" json:"nomor_peserta"`
	NISN          string         `gorm:"size:20;not null" json:"nisn"`
	Nama          string         `gorm:"size:255;not null" json:"nama"`
	TanggalLahir  time.Time      `gorm:"type:date;not null" json:"tanggal_lahir"`
	Nilai         datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"nilai"`
	Lulus         bool           `gorm:"not null;default:false" json:"lulus"`
	SKL           string         `gorm:"size:255" json:"skl"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedByID   *uint          `json:"created_by_id"`
	UpdatedByID   *uint          `json:"updated_by_id"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Foreign key relationships
	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	UpdatedBy *User `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`
}

// TableName specifies the table name for Kelulusan
func (m *Kelulusan) TableName() string {
	return "kelulusan"
}
