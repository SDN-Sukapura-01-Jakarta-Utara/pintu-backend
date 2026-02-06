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
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Role
func (m *Role) TableName() string {
	return "roles"
}
