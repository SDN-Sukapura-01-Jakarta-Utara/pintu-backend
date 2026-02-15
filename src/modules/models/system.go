package models

import (
	"time"

	"gorm.io/gorm"
)

// System represents the System model
type System struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Nama        string          `gorm:"not null" json:"nama"`
	Description string          `json:"description"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for System
func (m *System) TableName() string {
	return "systems"
}
