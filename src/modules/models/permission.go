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
	SystemID    *uint           `json:"system_id"`
	System      *System         `gorm:"foreignKey:SystemID" json:"system,omitempty"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Permission
func (m *Permission) TableName() string {
	return "permissions"
}
