package models

import (
	"time"

	"gorm.io/gorm"
)

// Application represents the Application model
type Application struct {
	ID               uint           `gorm:"primaryKey"`
	Nama             string         `gorm:"type:varchar(255);not null"`
	Link             string         `gorm:"type:varchar(500);not null"`
	ShowInJumbotron  bool           `gorm:"default:false"`
	Status           string         `gorm:"type:varchar(20);default:'active'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	CreatedByID      *uint
	UpdatedByID      *uint
}

// TableName specifies the table name for Application
func (m *Application) TableName() string {
	return "applications"
}
