package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Contact represents the Contact model
type Contact struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Alamat    string         `gorm:"type:text;not null" json:"alamat"`
	Telepon   string         `gorm:"not null" json:"telepon"`
	Email     string         `gorm:"not null" json:"email"`
	JamBuka   datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"jam_buka"`
	Gmaps     string         `gorm:"type:text" json:"gmaps"`
	Website   string         `json:"website"`
	Youtube   string         `json:"youtube"`
	Instagram string         `json:"instagram"`
	Tiktok    string         `json:"tiktok"`
	Facebook  string         `json:"facebook"`
	Twitter   string         `json:"twitter"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedByID *uint        `json:"created_by_id"`
	UpdatedByID *uint        `json:"updated_by_id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Contact
func (m *Contact) TableName() string {
	return "contacts"
}
