package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents the Permission model
type Permission struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	GroupName   string          `json:"group_name"`
	System      string          `json:"system"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Permission
func (m *Permission) TableName() string {
	return "permissions"
}
