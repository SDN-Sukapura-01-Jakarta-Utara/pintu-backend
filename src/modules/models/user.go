package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// User represents the User model
type User struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Nama             string         `gorm:"not null" json:"nama"`
	Username         string         `gorm:"uniqueIndex;not null" json:"username"`
	Password         string         `gorm:"not null" json:"-"`
	RoleID           *uint          `json:"role_id"`
	Role             *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	AccessibleSystem string         `gorm:"type:jsonb;default:'[]'" json:"accessible_system"`
	Status           string         `gorm:"default:active" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	CreatedByID      *uint          `json:"created_by_id"`
	UpdatedByID      *uint          `json:"updated_by_id"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for User
func (m *User) TableName() string {
	return "users"
}

// AccessibleSystems returns accessible systems as slice
func (u *User) AccessibleSystems() ([]string, error) {
	var systems []string
	if u.AccessibleSystem == "" || u.AccessibleSystem == "[]" {
		return systems, nil
	}

	err := json.Unmarshal([]byte(u.AccessibleSystem), &systems)
	return systems, err
}

// SetAccessibleSystems sets accessible systems from slice
func (u *User) SetAccessibleSystems(systems []string) error {
	data, err := json.Marshal(systems)
	if err != nil {
		return err
	}
	u.AccessibleSystem = string(data)
	return nil
}
