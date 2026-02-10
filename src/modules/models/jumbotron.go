package models

import (
	"time"

	"gorm.io/gorm"
)

// Jumbotron represents the Jumbotron model
type Jumbotron struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	File        string          `gorm:"not null" json:"file"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Jumbotron
func (m *Jumbotron) TableName() string {
	return "jumbotrons"
}
