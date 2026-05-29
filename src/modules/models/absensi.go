package models

import (
	"time"

	"gorm.io/gorm"
)

// Absensi represents the Absensi model
type Absensi struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	PesertaDidikID   uint           `gorm:"not null" json:"peserta_didik_id"`
	RombelID         *uint          `json:"rombel_id"`
	TahunPelajaranID uint           `gorm:"not null" json:"tahun_pelajaran_id"`
	Semester         int            `gorm:"not null" json:"semester"`
	Tanggal          time.Time      `gorm:"type:date;not null" json:"tanggal"`
	BidangStudiID    *uint          `json:"bidang_studi_id"` // NULL = guru kelas, NOT NULL = guru mapel
	PertemuanKe      *int           `json:"pertemuan_ke"`    // NULL = guru kelas, NOT NULL = guru mapel
	Status           string         `gorm:"not null" json:"status"` // hadir, sakit, izin, alpa
	WaktuAbsen       *time.Time     `json:"waktu_absen"`
	MetodeInput      string         `gorm:"not null" json:"metode_input"` // scan, manual
	Keterangan       string         `json:"keterangan"`
	FileSurat        string         `json:"file_surat"`
	DicatatOlehID    *uint          `json:"dicatat_oleh_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// Foreign key relationships
	PesertaDidik     *PesertaDidik     `gorm:"foreignKey:PesertaDidikID" json:"peserta_didik,omitempty"`
	Rombel           *Rombel           `gorm:"foreignKey:RombelID" json:"rombel,omitempty"`
	TahunPelajaran   *TahunPelajaran   `gorm:"foreignKey:TahunPelajaranID" json:"tahun_pelajaran,omitempty"`
	BidangStudi      *BidangStudi      `gorm:"foreignKey:BidangStudiID" json:"bidang_studi,omitempty"`
	DicatatOleh      *User             `gorm:"foreignKey:DicatatOlehID" json:"dicatat_oleh,omitempty"`
}

// TableName specifies the table name for Absensi
func (m *Absensi) TableName() string {
	return "absensi"
}
