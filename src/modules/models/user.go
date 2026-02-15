package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the User model
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Nama        string         `gorm:"not null" json:"nama"`
	Username    string         `gorm:"uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"not null" json:"-"`
	Status      string         `gorm:"default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedByID *uint          `json:"created_by_id"`
	UpdatedByID *uint          `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// Many-to-many relationship
	Roles       []Role         `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName specifies the table name for User
func (m *User) TableName() string {
	return "users"
}
