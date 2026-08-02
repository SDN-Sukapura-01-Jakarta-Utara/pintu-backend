package models

import (
	"time"

	"gorm.io/gorm"
)

// PesertaDidikEkstrakurikuler represents the junction table for student extracurricular registration
type PesertaDidikEkstrakurikuler struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	PesertaDidikRombelID  uint           `gorm:"not null;index" json:"peserta_didik_rombel_id"`
	EkstrakurikulerID     uint           `gorm:"not null;index" json:"ekstrakurikuler_id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedByID           *uint          `json:"created_by_id"`
	UpdatedByID           *uint          `json:"updated_by_id"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	PesertaDidikRombel *PesertaDidikRombel `gorm:"foreignKey:PesertaDidikRombelID" json:"peserta_didik_rombel,omitempty"`
	Ekstrakurikuler    *Ekstrakurikuler    `gorm:"foreignKey:EkstrakurikulerID" json:"ekstrakurikuler,omitempty"`
}

// TableName specifies the table name
func (m *PesertaDidikEkstrakurikuler) TableName() string {
	return "peserta_didik_ekstrakurikuler"
}
