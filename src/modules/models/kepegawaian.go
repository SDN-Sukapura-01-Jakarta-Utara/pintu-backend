package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Kepegawaian represents the Kepegawaian model
type Kepegawaian struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Nama                  string         `gorm:"not null" json:"nama"`
	Username              string         `gorm:"not null" json:"username"`
	Password              string         `gorm:"not null" json:"password"`
	NIP                   string         `gorm:"column:nip;not null;unique" json:"nip"`
	NKKI                  string         `gorm:"column:nkki" json:"nkki"`
	Foto                  string         `gorm:"column:foto" json:"foto"`
	Kategori              string         `gorm:"column:kategori" json:"kategori"`
	Jabatan               string         `gorm:"column:jabatan" json:"jabatan"`
	BidangStudiID         *uint          `gorm:"column:bidang_studi_id" json:"bidang_studi_id"`
	RombelGuruKelasID     *uint          `gorm:"column:rombel_guru_kelas_id" json:"rombel_guru_kelas_id"`
	RombelBidangStudi     datatypes.JSON `gorm:"column:rombel_bidang_studi;type:jsonb;default:'[]'" json:"rombel_bidang_studi"`
	KK                    string         `gorm:"column:kk" json:"kk"`
	AktaLahir             string         `gorm:"column:akta_lahir" json:"akta_lahir"`
	KTP                   string         `gorm:"column:ktp" json:"ktp"`
	IjazahSD              string         `gorm:"column:ijazah_sd" json:"ijazah_sd"`
	IjazahSMP             string         `gorm:"column:ijazah_smp" json:"ijazah_smp"`
	IjazahSMA             string         `gorm:"column:ijazah_sma" json:"ijazah_sma"`
	IjazahS1              string         `gorm:"column:ijazah_s1" json:"ijazah_s1"`
	IjazahS2              string         `gorm:"column:ijazah_s2" json:"ijazah_s2"`
	IjazahS3              string         `gorm:"column:ijazah_s3" json:"ijazah_s3"`
	SertifikatPendidik    string         `gorm:"column:sertifikat_pendidik" json:"sertifikat_pendidik"`
	SertifikatLainnya     datatypes.JSON `gorm:"column:sertifikat_lainnya;type:jsonb;default:'[]'" json:"sertifikat_lainnya"`
	SK                    string         `gorm:"column:sk" json:"sk"`
	DokumenLainnya        datatypes.JSON `gorm:"column:dokumen_lainnya;type:jsonb;default:'[]'" json:"dokumen_lainnya"`
	Status                string         `gorm:"column:status;default:active" json:"status"`
	CreatedAt             time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at" json:"updated_at"`
	CreatedByID           *uint          `gorm:"column:created_by_id" json:"created_by_id"`
	UpdatedByID           *uint          `gorm:"column:updated_by_id" json:"updated_by_id"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Kepegawaian
func (m *Kepegawaian) TableName() string {
	return "kepegawaian"
}
