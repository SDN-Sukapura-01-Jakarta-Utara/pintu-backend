package models

import (
	"time"

	"gorm.io/gorm"
)

// Rombel represents the Rombel model
type Rombel struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"uniqueIndex;not null" json:"name"`
	Status      string          `gorm:"default:active" json:"status"`
	KelasID     uint            `gorm:"index" json:"kelas_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Rombel
func (m *Rombel) TableName() string {
	return "rombels"
}
