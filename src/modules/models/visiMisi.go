package models

import (
	"time"
)

// VisiMisi represents the Visi Misi model
type VisiMisi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Visi        string    `gorm:"type:text;not null" json:"visi"`
	Misi        string    `gorm:"type:text;not null" json:"misi"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// TableName specifies the table name for VisiMisi
func (m *VisiMisi) TableName() string {
	return "visi_misi"
}
