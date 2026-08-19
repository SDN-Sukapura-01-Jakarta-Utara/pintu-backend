package models

import (
	"time"

	"gorm.io/gorm"
)

// Pelatih represents ekstrakurikuler coach/instructor
type Pelatih struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Nama        string          `gorm:"type:varchar(255);not null" json:"nama"`
	Username    *string         `gorm:"type:varchar(100);unique" json:"username"`
	Password    string          `gorm:"type:varchar(255)" json:"-"`
	Telepon     string          `gorm:"type:varchar(20)" json:"telepon"`
	Alamat      string          `gorm:"type:text" json:"alamat"`
	FotoProfil  *string         `gorm:"type:varchar(512)" json:"foto_profil"`
	Keahlian    string          `gorm:"type:text" json:"keahlian"`
	Sertifikat  *string         `gorm:"type:json" json:"sertifikat"` // JSON array
	Status      string          `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	CreatedByID *uint           `json:"created_by_id,omitempty"`
	UpdatedByID *uint           `json:"updated_by_id,omitempty"`

	// Relations
	Ekstrakurikuler []Ekstrakurikuler `gorm:"many2many:pelatih_ekstrakurikuler;" json:"ekstrakurikuler,omitempty"`
	Roles           []Role            `gorm:"many2many:pelatih_roles" json:"roles,omitempty"`
}

// TableName specifies the table name for Pelatih model
func (Pelatih) TableName() string {
	return "pelatih"
}
