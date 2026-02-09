package models

import (
	"time"

	"gorm.io/gorm"
)

// TahunPelajaran represents the TahunPelajaran model
type TahunPelajaran struct {
	ID        uint            `gorm:"primaryKey"`
	// Add fields here
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for TahunPelajaran
func (m *TahunPelajaran) TableName() string {
	return "tahunpelajarans"
}
