package models

import (
	"time"

	"gorm.io/gorm"
)

// StrukturOrganisasi represents the Struktur Organisasi model
type StrukturOrganisasi struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	PegawaiID          *uint          `gorm:"column:pegawai_id" json:"pegawai_id"`
	NamaNonPegawai     string         `gorm:"column:nama_non_pegawai" json:"nama_non_pegawai"`
	JabatanNonPegawai  string         `gorm:"column:jabatan_non_pegawai" json:"jabatan_non_pegawai"`
	Urutan             int            `gorm:"column:urutan;not null" json:"urutan"`
	Relasi             string         `gorm:"column:relasi;not null" json:"relasi"`
	Status             string         `gorm:"column:status;default:active" json:"status"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	CreatedByID        *uint          `gorm:"column:created_by_id" json:"created_by_id"`
	UpdatedByID        *uint          `gorm:"column:updated_by_id" json:"updated_by_id"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	// Foreign key relationships
	Pegawai            *Kepegawaian   `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	CreatedBy          *User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	UpdatedBy          *User          `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`
}

// TableName specifies the table name for StrukturOrganisasi
func (m *StrukturOrganisasi) TableName() string {
	return "struktur_organisasi"
}
