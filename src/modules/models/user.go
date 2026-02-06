package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the User model
type User struct {
	ID        uint            `gorm:"primaryKey"`
	// Add fields here
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for User
func (m *User) TableName() string {
	return "users"
}
