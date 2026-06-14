package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// SettingLayananSPMBRepository handles data operations for Setting Layanan SPMB
type SettingLayananSPMBRepository interface {
	GetByID(id uint) (*models.SettingLayananSPMB, error)
	Create(data *models.SettingLayananSPMB) error
	Update(data *models.SettingLayananSPMB) error
}

type SettingLayananSPMBRepositoryImpl struct {
	db *gorm.DB
}

// NewSettingLayananSPMBRepository creates a new Setting Layanan SPMB repository
func NewSettingLayananSPMBRepository(db *gorm.DB) SettingLayananSPMBRepository {
	return &SettingLayananSPMBRepositoryImpl{db: db}
}

// GetByID retrieves Setting Layanan SPMB by ID
func (r *SettingLayananSPMBRepositoryImpl) GetByID(id uint) (*models.SettingLayananSPMB, error) {
	var data models.SettingLayananSPMB
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Create creates a new Setting Layanan SPMB record
func (r *SettingLayananSPMBRepositoryImpl) Create(data *models.SettingLayananSPMB) error {
	return r.db.Create(data).Error
}

// Update updates Setting Layanan SPMB record
func (r *SettingLayananSPMBRepositoryImpl) Update(data *models.SettingLayananSPMB) error {
	return r.db.Save(data).Error
}
