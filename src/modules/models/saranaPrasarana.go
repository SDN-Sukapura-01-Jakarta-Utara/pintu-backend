package models

import (
	"time"

	"gorm.io/gorm"
)

// SaranaPrasarana represents the Sarana Prasarana model
type SaranaPrasarana struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Foto        string          `gorm:"not null" json:"foto"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for SaranaPrasarana
func (m *SaranaPrasarana) TableName() string {
	return "sarana_prasaranas"
}
