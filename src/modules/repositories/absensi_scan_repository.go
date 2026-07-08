package repositories

import (
	"pintu-backend/src/modules/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AbsensiScanRepository defines the interface for Absensi Scan repository
type AbsensiScanRepository interface {
	GetPesertaDidikByBarcode(barcode string) (*models.PesertaDidik, error)
	GetKonfigurasiAbsensi() (*models.KonfigurasiAbsensi, error)
	GetAbsensiByPesertaDidikAndDate(pesertaDidikID uint, tanggal time.Time) (*models.Absensi, error)
	UpsertAbsensi(absensi *models.Absensi) error
}

type AbsensiScanRepositoryImpl struct {
	db *gorm.DB
}

// NewAbsensiScanRepository creates a new Absensi Scan repository
func NewAbsensiScanRepository(db *gorm.DB) AbsensiScanRepository {
	return &AbsensiScanRepositoryImpl{db: db}
}

// GetPesertaDidikByBarcode retrieves peserta didik by barcode
func (r *AbsensiScanRepositoryImpl) GetPesertaDidikByBarcode(barcode string) (*models.PesertaDidik, error) {
	var pesertaDidik models.PesertaDidik
	if err := r.db.Where("barcode = ?", barcode).First(&pesertaDidik).Error; err != nil {
		return nil, err
	}
	return &pesertaDidik, nil
}

// GetKonfigurasiAbsensi retrieves konfigurasi absensi with ID = 1
func (r *AbsensiScanRepositoryImpl) GetKonfigurasiAbsensi() (*models.KonfigurasiAbsensi, error) {
	var config models.KonfigurasiAbsensi
	if err := r.db.First(&config, 1).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAbsensiByPesertaDidikAndDate retrieves absensi by peserta didik ID and date
func (r *AbsensiScanRepositoryImpl) GetAbsensiByPesertaDidikAndDate(pesertaDidikID uint, tanggal time.Time) (*models.Absensi, error) {
	var absensi models.Absensi
	if err := r.db.Where("peserta_didik_id = ? AND tanggal = ?", pesertaDidikID, tanggal.Format("2006-01-02")).First(&absensi).Error; err != nil {
		return nil, err
	}
	return &absensi, nil
}

// UpsertAbsensi creates or updates absensi record using UPSERT
func (r *AbsensiScanRepositoryImpl) UpsertAbsensi(absensi *models.Absensi) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "peserta_didik_id"}, {Name: "tanggal"}},
		DoUpdates: clause.AssignmentColumns([]string{"jam_datang", "jam_pulang", "status", "updated_at"}),
	}).Create(absensi).Error
}
