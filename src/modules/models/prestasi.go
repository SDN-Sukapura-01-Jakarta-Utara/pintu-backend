package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FotoItem represents a single photo in the foto array
type FotoItem struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	Thumbnail string `json:"thumbnail"` // "active" or "inactive"
}

// AnggotaTimPrestasi represents the anggota tim prestasi model
type AnggotaTimPrestasi struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	PrestasiID       uint           `gorm:"not null" json:"prestasi_id"`
	PesertaDidikID   uint           `gorm:"not null" json:"peserta_didik_id"`
	TahunPelajaranID uint           `gorm:"not null;column:tahun_pelajaran_id" json:"tahun_pelajaran_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedByID      *uint          `json:"created_by_id"`
	UpdatedByID      *uint          `json:"updated_by_id"`
	// Foreign key relationships
	Prestasi         *Prestasi      `gorm:"foreignKey:PrestasiID" json:"prestasi,omitempty"`
	PesertaDidik     *PesertaDidik  `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
	TahunPelajaran   *TahunPelajaran `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
}

// Prestasi represents the Prestasi model
type Prestasi struct {
	ID                 uint                 `gorm:"primaryKey" json:"id"`
	PesertaDidikID     *uint                `gorm:"column:peserta_didik_id" json:"peserta_didik_id"`
	Jenis              string               `gorm:"not null" json:"jenis"`
	NamaGrup           string               `json:"nama_grup"`
	NamaPrestasi       string               `gorm:"not null" json:"nama_prestasi"`
	TingkatPrestasi    string               `json:"tingkat_prestasi"`
	Penyelenggara      string               `json:"penyelenggara"`
	TanggalLomba       time.Time            `gorm:"not null" json:"tanggal_lomba"`
	Juara              string               `gorm:"not null" json:"juara"`
	Keterangan         string               `gorm:"type:text" json:"keterangan"`
	Foto               datatypes.JSON       `gorm:"type:jsonb;default:'[]'" json:"foto"`
	EkstrakurikulerID  *uint                `gorm:"column:ekstrakurikuler_id" json:"ekstrakurikuler_id"`
	TahunPelajaranID   uint                 `gorm:"not null;column:tahun_pelajaran_id" json:"tahun_pelajaran_id"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	DeletedAt          gorm.DeletedAt       `gorm:"index" json:"deleted_at,omitempty"`
	CreatedByID        *uint                `json:"created_by_id"`
	UpdatedByID        *uint                `json:"updated_by_id"`
	// Foreign key relationships
	PesertaDidik       *PesertaDidik        `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
	Ekstrakurikuler    *Ekstrakurikuler     `gorm:"foreignKey:EkstrakurikulerID" json:"ekstrakurikuler,omitempty"`
	TahunPelajaran     *TahunPelajaran      `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	// One-to-many relationship
	AnggotaTimPrestasi []AnggotaTimPrestasi `gorm:"foreignKey:PrestasiID" json:"anggota_tim_prestasi,omitempty"`
}

// TableName specifies the table name for Prestasi
func (m *Prestasi) TableName() string {
	return "prestasi"
}

// TableName specifies the table name for AnggotaTimPrestasi
func (m *AnggotaTimPrestasi) TableName() string {
	return "anggota_tim_prestasi"
}

// Scan implements the sql.Scanner interface for JSONB
func (f *FotoItem) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &f)
}

// Value implements the driver.Valuer interface for JSONB
func (f FotoItem) Value() (driver.Value, error) {
	return json.Marshal(f)
}