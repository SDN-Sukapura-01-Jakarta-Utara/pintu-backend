package models

import (
	"time"

	"gorm.io/gorm"
)

// KritikSaran represents the KritikSaran model
type KritikSaran struct {
	ID          uint           `gorm:"primaryKey"`
	Nama        string         `gorm:"type:varchar(255);not null"`
	KritikSaran string         `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for KritikSaran
func (m *KritikSaran) TableName() string {
	return "kritik_saran"
}
