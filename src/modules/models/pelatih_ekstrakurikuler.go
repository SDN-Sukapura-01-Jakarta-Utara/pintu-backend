package models

import (
	"time"

	"gorm.io/gorm"
)

// PelatihEkstrakurikuler represents mapping between coach and ekstrakurikuler
type PelatihEkstrakurikuler struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	PelatihID         uint            `gorm:"not null" json:"pelatih_id"`
	EkstrakurikulerID uint            `gorm:"not null" json:"ekstrakurikuler_id"`
	Status            string          `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	CreatedByID       *uint           `json:"created_by_id,omitempty"`
	UpdatedByID       *uint           `json:"updated_by_id,omitempty"`

	// Relations
	Pelatih         *Pelatih         `gorm:"foreignKey:PelatihID" json:"pelatih,omitempty"`
	Ekstrakurikuler *Ekstrakurikuler `gorm:"foreignKey:EkstrakurikulerID" json:"ekstrakurikuler,omitempty"`
}

// TableName specifies the table name
func (PelatihEkstrakurikuler) TableName() string {
	return "pelatih_ekstrakurikuler"
}
