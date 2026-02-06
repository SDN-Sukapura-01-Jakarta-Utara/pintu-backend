package models

import (
	"time"
)

// RolePermission represents the RolePermission model (pivot table)
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	PermissionID uint      `gorm:"not null" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName specifies the table name for RolePermission
func (m *RolePermission) TableName() string {
	return "role_permissions"
}
