package models

import (
	"time"

	"gorm.io/gorm"
)

// TahunPelajaran represents the TahunPelajaran model
type TahunPelajaran struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	TahunPelajaran  string          `gorm:"uniqueIndex;not null" json:"tahun_pelajaran"`
	Status          string          `gorm:"default:active" json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedByID     *uint           `json:"created_by_id"`
	UpdatedByID     *uint           `json:"updated_by_id"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for TahunPelajaran
func (m *TahunPelajaran) TableName() string {
	return "tahun_pelajarans"
}
