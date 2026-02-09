package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents the Role model
type Role struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	System      string          `json:"system"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Role
func (m *Role) TableName() string {
	return "roles"
}
